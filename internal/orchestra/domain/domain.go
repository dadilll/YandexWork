package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Task struct {
	ID         string    `json:"id"`
	Expression string    `json:"expression"`
	Status     string    `json:"status"`
	StartTime  time.Time `json:"start_time"`
	Result     float64   `json:"result"`
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
	taskKey := fmt.Sprintf(taskID)
	task := Task{ID: taskID, Expression: expression, Status: "pending", StartTime: time.Now()}

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

	err = o.Redis.Set(context.Background(), taskKey, data, 0).Err()
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
	for _, taskKey := range taskKeys {
		data, err := o.Redis.Get(context.Background(), taskKey).Result()
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
	taskKey := fmt.Sprintf("task:%s", id)
	data, err := o.Redis.Get(context.Background(), taskKey).Result()
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
