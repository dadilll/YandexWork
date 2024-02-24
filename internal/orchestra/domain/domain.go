package domain

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Task struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

type Orchestrator struct {
	Redis          *redis.Client
	Agents         []*Agent
	processedTasks map[string]bool
}

type Agent struct {
	ID          int
	TaskChannel chan Task
}

func NewOrchestrator(redis *redis.Client) *Orchestrator {
	return &Orchestrator{
		Redis:          redis,
		processedTasks: make(map[string]bool),
	}
}

func (o *Orchestrator) AddTask(expression string) (string, error) {
	taskID := generateTaskID()
	task := Task{ID: taskID, Expression: expression, Status: "pending"} // Используем config.Task здесь

	// Проверяем, была ли уже обработана задача с таким ID
	if _, exists := o.processedTasks[taskID]; exists {
		// Если да, возвращаем успешный ответ с кодом 200
		return taskID, nil
	}

	data, err := json.Marshal(task)
	if err != nil {
		log.Println("Error marshaling task:", err)
		return "", err
	}

	err = o.Redis.Set(context.Background(), taskID, data, 0).Err()
	if err != nil {
		log.Println("Error saving task to Redis:", err)
		return "", err
	}

	// Добавляем ID задачи в список обработанных
	o.processedTasks[taskID] = true

	return taskID, nil
}

func (o *Orchestrator) GetTasks() []Task {
	taskKeys, err := o.Redis.Keys(context.Background(), "*").Result()
	if err != nil {
		log.Println("Error getting task keys from Redis:", err)
		return nil
	}

	var tasks []Task
	for _, taskID := range taskKeys { // Используйте просто taskID вместо "task:" + taskID
		data, err := o.Redis.Get(context.Background(), taskID).Result()
		if err != nil {
			log.Println("Error getting task from Redis:", err)
			continue
		}

		var task Task
		err = json.Unmarshal([]byte(data), &task)
		if err != nil {
			log.Println("Error unmarshaling task:", err)
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func (o *Orchestrator) GetTaskByID(id string) *Task {
	data, err := o.Redis.Get(context.Background(), id).Result()
	if err != nil {
		log.Println("Error getting task from Redis:", err)
		return nil
	}

	var task Task
	err = json.Unmarshal([]byte(data), &task)
	if err != nil {
		log.Println("Error unmarshaling task:", err)
		return nil
	}

	return &task
}
func generateTaskID() string {
	taskID := uuid.New()
	return taskID.String()
}
