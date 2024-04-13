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
	executingLock sync.Map
	durationMap   map[string]int
}

func NewAgent(id int, postgres *sqlx.DB, workers int, durationMap map[string]int) *Agent {
	log.Printf("Initializing agent with ID: %d", id)
	agent := &Agent{
		ID:          id,
		Postgres:    postgres,
		TaskQueues:  make([]chan Task, workers),
		Workers:     workers,
		durationMap: durationMap,
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

func (a *Agent) worker(workerID int) {
	for task := range a.TaskQueues[workerID] {
		// Пытаемся заблокировать задачу
		_, err := a.Postgres.Exec("INSERT INTO locks (id, status) VALUES ($1, 'locked') ON CONFLICT(id) DO NOTHING", task.ID)
		if err != nil {
			log.Printf("Error setting lock for task %s: %v", task.ID, err)
			continue
		}

		// Помечаем задачу как обрабатываемую этим воркером
		a.markTaskAsBeingProcessed(task.ID)

		log.Printf("Agent %d: Worker %d started processing task %s", a.ID, workerID, task.ID)
		a.processTask(task)
		log.Printf("Agent %d: Worker %d finished processing task %s", a.ID, workerID, task.ID)

		// Снимаем блокировку
		_, err = a.Postgres.Exec("DELETE FROM locks WHERE id = $1", task.ID)
		if err != nil {
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
		queueIndex := a.getQueueIndex(task.ID)
		a.TaskQueues[queueIndex] <- task
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over task rows: %v", err)
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
	result, err := expression.ParseExpression(task.Expression, a.durationMap)
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
