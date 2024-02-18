package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/Dadil/project/internal/agent/expression"
	"github.com/go-redis/redis/v8"
)

type Task struct {
	ID         string    `json:"id"`
	Expression string    `json:"expression"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	Result     float64   `json:"result"`
}

type Agent struct {
	ID         int
	Redis      *redis.Client
	TaskQueues []chan Task // Вот это новое поле
	Workers    int
	taskLocks  map[string]*sync.Mutex
	lockMap    sync.Map
}

func NewAgent(id int, redis *redis.Client, workers int) *Agent {
	log.Printf("Initializing agent with ID: %d", id)
	agent := &Agent{
		ID:         id,
		Redis:      redis,
		TaskQueues: make([]chan Task, workers), // Создание слайса каналов
		Workers:    workers,
		taskLocks:  make(map[string]*sync.Mutex),
	}

	// Инициализируем мьютексы для блокировки задач
	for i := 0; i < workers; i++ {
		agent.TaskQueues[i] = make(chan Task)
		agent.taskLocks[fmt.Sprintf("worker-%d", i)] = &sync.Mutex{}
	}

	return agent
}

func (a *Agent) Start() {
	log.Printf("Agent %d is starting %d workers", a.ID, a.Workers)
	// Запуск воркеров
	for i := 0; i < a.Workers; i++ {
		go a.worker(i) // Передаем индекс воркера в качестве аргумента
	}

	// Начало проверки задач из Redis
	go func() {
		// Бесконечный цикл для периодической проверки
		for {
			// Проверяем задачи в базе данных Redis
			a.checkTasks()
			// Ждем 5 секунд перед следующей проверкой
			time.Sleep(5 * time.Second)
		}
	}()
}

func (a *Agent) worker(workerID int) {
	for task := range a.TaskQueues[workerID] {
		if a.isTaskBeingProcessed(task.ID) {
			// Задача уже обрабатывается другим воркером, пропускаем её
			continue
		}

		// Помечаем задачу как обрабатываемую этим воркером
		a.markTaskAsBeingProcessed(task.ID)

		log.Printf("Agent %d: Worker %d started processing task %s", a.ID, workerID, task.ID)
		a.processTask(task)
		log.Printf("Agent %d: Worker %d finished processing task %s", a.ID, workerID, task.ID)
	}
}

func (a *Agent) isTaskBeingProcessed(taskID string) bool {
	_, ok := a.lockMap.Load(taskID)
	return ok
}

func (a *Agent) markTaskAsBeingProcessed(taskID string) {
	a.lockMap.Store(taskID, true)
}

func (a *Agent) checkTasks() {
	tasks, err := a.Redis.Keys(context.Background(), "*").Result()
	if err != nil {
		log.Printf("Error getting tasks from Redis: %v", err)
		return
	}

	for _, taskKey := range tasks {
		taskJSON, err := a.Redis.Get(context.Background(), taskKey).Bytes()
		if err != nil {
			log.Printf("Error getting task %s from Redis: %v", taskKey, err)
			continue
		}

		var task Task
		err = json.Unmarshal(taskJSON, &task)
		if err != nil {
			log.Printf("Error decoding task %s: %v", taskKey, err)
			continue
		}

		// Проверяем статус задачи, пропускаем выполненные и задачи с ошибкой
		if task.Status == "completed" || task.Status == "error" {
			continue
		}

		// Определяем индекс очереди задач
		queueIndex := a.getQueueIndex(task.ID)
		a.TaskQueues[queueIndex] <- task
	}
}

func (a *Agent) getQueueIndex(taskID string) int {
	return int(hash(taskID)) % a.Workers
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func (a *Agent) processTask(task Task) {
	op1, op2, operator, duration, err := expression.ParseExpression(task.Expression)
	if err != nil {
		log.Printf("Error parsing expression for task %s: %s", task.ID, err)
		task.Status = "error"
	} else {
		task.Result, err = expression.EvaluateExpression(op1, op2, operator, duration)
		if err != nil {
			log.Printf("Error evaluating expression for task %s: %s", task.ID, err)
			task.Status = "error"
		} else {
			task.Status = "completed"
		}
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		log.Printf("Error marshaling task %s: %s", task.ID, err)
		return
	}

	err = a.Redis.Set(context.Background(), task.ID, taskJSON, 0).Err()
	if err != nil {
		log.Printf("Error updating task %s: %s", task.ID, err)
	}
}
