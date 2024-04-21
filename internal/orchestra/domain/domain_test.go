package domain_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dadil/project/internal/orchestra/domain" // Update with your project's import path
)

func TestAddTask(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new orchestrator with the mock database
	orchestrator := domain.NewOrchestrator(db)

	// Define the expected SQL query and mock behavior
	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function under test
	taskID, err := orchestrator.AddTask("2 + 2")
	if err != nil {
		t.Fatalf("Error adding task: %v", err)
	}

	// Check if task ID is returned
	if taskID == "" {
		t.Error("Expected non-empty task ID, got empty")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetTasks(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new orchestrator with the mock database
	orchestrator := domain.NewOrchestrator(db)

	// Define the expected SQL query and mock behavior
	rows := sqlmock.NewRows([]string{"id", "expression", "status", "result"}).
		AddRow("1", "2 + 2", "pending", 0.0).
		AddRow("2", "3 * 3", "completed", 9.0)
	mock.ExpectQuery("SELECT id, expression, status, result FROM tasks").WillReturnRows(rows)

	// Call the function under test
	tasks := orchestrator.GetTasks()

	// Check if tasks are returned
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddTaskForUser(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new orchestrator with the mock database
	orchestrator := domain.NewOrchestrator(db)

	// Define the expected SQL queries and mock behavior
	mock.ExpectQuery("SELECT id FROM users WHERE login = ?").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO user_tasks").
		WithArgs("1", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the function under test
	taskID, err := orchestrator.AddTaskForUser("2 + 2", "testuser")
	if err != nil {
		t.Fatalf("Error adding task for user: %v", err)
	}

	// Check if task ID is returned
	if taskID == "" {
		t.Error("Expected non-empty task ID, got empty")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
func TestGetUserByLogin(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new orchestrator with the mock database
	orchestrator := domain.NewOrchestrator(db)

	// Define the expected SQL query and mock behavior
	rows := sqlmock.NewRows([]string{"login", "password"}).
		AddRow("testuser", "hashed_password")
	mock.ExpectQuery("SELECT login, password FROM users WHERE login = ?").
		WithArgs("testuser").
		WillReturnRows(rows)

	// Call the function under test
	user, err := orchestrator.GetUserByLogin("testuser")
	if err != nil {
		t.Fatalf("Error getting user by login: %v", err)
	}

	// Check if user is returned
	if user == nil {
		t.Error("Expected user, got nil")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteAllTasksForUser(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new orchestrator with the mock database
	orchestrator := domain.NewOrchestrator(db)

	// Define the expected SQL queries and mock behavior
	mock.ExpectExec("DELETE FROM tasks WHERE id IN").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DELETE FROM user_tasks WHERE user_id =").
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the function under test
	err = orchestrator.DeleteAllTasksForUser("testuser")
	if err != nil {
		t.Fatalf("Error deleting tasks for user: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetTasksForUser(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Create a new orchestrator with the mock database
	orchestrator := domain.NewOrchestrator(db)

	// Define the expected SQL query and mock behavior
	rows := sqlmock.NewRows([]string{"id", "expression", "status", "result"}).
		AddRow("1", "2 + 2", "pending", 0.0).
		AddRow("2", "3 * 3", "completed", 9.0)
	mock.ExpectQuery("SELECT t.id, t.expression, t.status, t.result FROM tasks t").
		WillReturnRows(rows)

	// Call the function under test
	tasks := orchestrator.GetTasksForUser("testuser")

	// Check if tasks are returned
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
