package domain

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Task struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
}

type User struct {
	Login    string
	Password string
	ID       int
}

type Orchestrator struct {
	DB             *sql.DB
	Agents         []*Agent
	processedTasks map[string]bool
}

type Agent struct {
	ID          int
	TaskChannel chan Task
}

func NewOrchestrator(db *sql.DB) *Orchestrator {
	return &Orchestrator{
		DB:             db,
		processedTasks: make(map[string]bool),
	}
}

func (o *Orchestrator) AddTask(expression string) (string, error) {
	taskID := generateTaskID()
	task := Task{ID: taskID, Expression: expression, Status: "pending"}

	// Проверяем, была ли уже обработана задача с таким ID
	if _, exists := o.processedTasks[taskID]; exists {
		// Если да, возвращаем успешный ответ с кодом 200
		return taskID, nil
	}

	_, err := o.DB.Exec("INSERT INTO tasks (id, expression, status, result) VALUES ($1, $2, $3, $4)",
		taskID, task.Expression, task.Status, task.Result)
	if err != nil {
		log.Println("Error saving task to PostgreSQL:", err)
		return "", err
	}

	// Добавляем ID задачи в список обработанных
	o.processedTasks[taskID] = true

	return taskID, nil
}

func (o *Orchestrator) GetTasks() []Task {
	rows, err := o.DB.Query("SELECT id, expression, status, result FROM tasks")
	if err != nil {
		log.Println("Error getting tasks from PostgreSQL:", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Expression, &task.Status, &task.Result)
		if err != nil {
			log.Println("Error scanning task:", err)
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func (o *Orchestrator) GetTaskByID(id string) *Task {
	var task Task
	err := o.DB.QueryRow("SELECT id, expression, status, result FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Expression, &task.Status, &task.Result)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Task not found:", id)
			return nil
		}
		log.Println("Error getting task from PostgreSQL:", err)
		return nil
	}
	return &task
}

func (o *Orchestrator) GetTasksByUserID(userID string) []Task {
	rows, err := o.DB.Query("SELECT t.id, t.expression, t.status, t.result FROM tasks t JOIN user_tasks ut ON t.id = ut.task_id WHERE ut.user_id = $1", userID)
	if err != nil {
		log.Println("Error getting tasks from PostgreSQL:", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Expression, &task.Status, &task.Result)
		if err != nil {
			log.Println("Error scanning task:", err)
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func generateTaskID() string {
	taskID := uuid.New()
	return taskID.String()
}

func (o *Orchestrator) CreateUser(login, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	_, err = o.DB.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, string(hashedPassword))
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}

func (o *Orchestrator) GetUserByLogin(login string) (*User, error) {
	var user User
	err := o.DB.QueryRow("SELECT login, password FROM users WHERE login = $1", login).
		Scan(&user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found:", login)
			return nil, err
		}
		log.Println("Error getting user from PostgreSQL:", err)
		return nil, err
	}

	return &user, nil
}
