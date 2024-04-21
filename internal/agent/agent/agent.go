package agent

import (
	"hash/fnv"
	"log"
	"sync"
	"time"

	"github.com/Dadil/project/internal/agent/expression"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

type Agent struct {
	ID            int
	Postgres      *sqlx.DB
	TaskQueues    []chan Task
	Workers       int
	ExecutingLock sync.Map
	DurationMap   map[string]int
}

func NewAgent(id int, postgres *sqlx.DB, workers int, durationMap map[string]int) *Agent {
	log.Printf("Initializing agent with ID: %d", id)
	agent := &Agent{
		ID:          id,
		Postgres:    postgres,
		TaskQueues:  make([]chan Task, workers),
		Workers:     workers,
		DurationMap: durationMap,
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
		go a.Worker(i) // Передаем индекс воркера в качестве аргумента
	}

	// Начало проверки задач из PostgreSQL
	go func() {
		// Бесконечный цикл для периодической проверки
		for {
			// Проверяем задачи в базе данных PostgreSQL
			a.checkTasks()
			// Ждем 5 секунд перед следующей проверкой
			time.Sleep(5 * time.Second)
		}
	}()
}

func (a *Agent) Worker(workerID int) {
	for task := range a.TaskQueues[workerID] {
		// Пытаемся заблокировать задачу
		_, err := a.Postgres.Exec("INSERT INTO locks (id, status) VALUES ($1, 'locked') ON CONFLICT(id) DO NOTHING", task.ID)
		if err != nil {
			log.Printf("Error setting lock for task %s: %v", task.ID, err)
			continue
		}

		// Помечаем задачу как обрабатываемую этим воркером
		a.MarkTaskAsBeingProcessed(task.ID)

		log.Printf("Agent %d: Worker %d started processing task %s", a.ID, workerID, task.ID)
		a.ProcessTask(task)
		log.Printf("Agent %d: Worker %d finished processing task %s", a.ID, workerID, task.ID)

		// Снимаем блокировку
		_, err = a.Postgres.Exec("DELETE FROM locks WHERE id = $1", task.ID)
		if err != nil {
			log.Printf("Error removing lock for task %s: %v", task.ID, err)
		}

		// По завершении обработки задачи освобождаем её
		a.MarkTaskAsFinished(task.ID)
	}
}

func (a *Agent) MarkTaskAsBeingProcessed(taskID string) {
	a.ExecutingLock.Store(taskID, true)
}

func (a *Agent) MarkTaskAsFinished(taskID string) {
	a.ExecutingLock.Delete(taskID)
}

func (a *Agent) checkTasks() {
	rows, err := a.Postgres.Query("SELECT id, expression, status FROM tasks WHERE status != 'completed' AND status != 'error'")
	if err != nil {
		log.Printf("Error getting tasks from PostgreSQL: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Expression, &task.Status); err != nil {
			log.Printf("Error scanning task from PostgreSQL: %v", err)
			continue
		}
		// Определяем индекс очереди задач
		queueIndex := a.GetQueueIndex(task.ID)
		a.TaskQueues[queueIndex] <- task
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over task rows: %v", err)
	}
}

func (a *Agent) GetQueueIndex(taskID string) int {
	return int(hash(taskID)) % a.Workers
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func (a *Agent) ProcessTask(task Task) {
	result, err := expression.ParseExpression(task.Expression, a.DurationMap)
	if err != nil {
		log.Printf("Error parsing expression for task %s: %s", task.ID, err)
		task.Status = "error"
	} else {
		task.Result = result
		task.Status = "completed"
	}

	// Обновляем задачу в базе данных PostgreSQL
	_, err = a.Postgres.Exec("UPDATE tasks SET result = $1, status = $2 WHERE id = $3", task.Result, task.Status, task.ID)
	if err != nil {
		log.Printf("Error updating task %s in PostgreSQL: %v", task.ID, err)
		return
	}
}
