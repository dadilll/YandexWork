package agent

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/Dadil/project/internal/agent/expression"
	"github.com/go-redis/redis/v8"
)

type Task struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

type Agent struct {
	ID            int
	Redis         *redis.Client
	TaskQueues    []chan Task // Используем config.Task здесь
	Workers       int
	executingLock sync.Map // Map для отслеживания выполняющихся задач
}

func NewAgent(id int, redis *redis.Client, workers int) *Agent {
	log.Printf("Initializing agent with ID: %d", id)
	agent := &Agent{
		ID:         id,
		Redis:      redis,
		TaskQueues: make([]chan Task, workers), // Используем config.Task здесь
		Workers:    workers,
	}

	// Инициализируем каналы задач для воркеров
	for i := 0; i < workers; i++ {
		agent.TaskQueues[i] = make(chan Task)
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
		// Пытаемся заблокировать задачу
		lockKey := "lock:" + task.ID
		locked, err := a.Redis.SetNX(context.Background(), lockKey, "locked", time.Minute).Result()
		if err != nil {
			log.Printf("Error setting lock for task %s: %v", task.ID, err)
			continue
		}
		if !locked {
			// Задача уже обрабатывается другим агентом, пропускаем её
			continue
		}

		// Помечаем задачу как обрабатываемую этим воркером
		a.markTaskAsBeingProcessed(task.ID)

		log.Printf("Agent %d: Worker %d started processing task %s", a.ID, workerID, task.ID)
		a.processTask(task)
		log.Printf("Agent %d: Worker %d finished processing task %s", a.ID, workerID, task.ID)

		// Снимаем блокировку
		if err := a.Redis.Del(context.Background(), lockKey).Err(); err != nil {
			log.Printf("Error removing lock for task %s: %v", task.ID, err)
		}

		// По завершении обработки задачи освобождаем её
		a.markTaskAsFinished(task.ID)
	}
}

func (a *Agent) markTaskAsBeingProcessed(taskID string) {
	a.executingLock.Store(taskID, true)
}

func (a *Agent) markTaskAsFinished(taskID string) {
	a.executingLock.Delete(taskID)
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
	result, err := expression.ParseExpression(task.Expression)
	if err != nil {
		log.Printf("Error parsing expression for task %s: %s", task.ID, err)
		task.Status = "error"
	} else {
		task.Result = result
		task.Status = "completed"
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
