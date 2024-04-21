package agent_test

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Dadil/project/internal/agent/agent"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAgent(t *testing.T) {
	// Создаем мок базы данных
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Создаем экземпляр *sql.DB
	sqlDB := sqlx.NewDb(mockDB, "sqlmock")

	// Создаем экземпляр агента с моком базы данных
	testAgent := &agent.Agent{Postgres: sqlDB}

	// Устанавливаем ожидания для запросов к базе данных
	rows := sqlmock.NewRows([]string{"id", "expression", "status", "result"}).
		AddRow(1, "test1", "completed", 1.0).
		AddRow(2, "test2", "completed", 1.0).
		AddRow(3, "test3", "completed", 1.0)

	mock.ExpectQuery("SELECT id, expression, status FROM tasks").
		WillReturnRows(rows)

	// Запускаем агента
	testAgent.Start()

	// Ждем, чтобы агент обработал задачи
	time.Sleep(1 * time.Second)

	// Проверяем, что все ожидания были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestMarkTaskAsBeingProcessed(t *testing.T) {
	// Создаем тестовый агент
	testAgent := &agent.Agent{}

	// Вызываем функцию MarkTaskAsBeingProcessed с тестовым taskID
	testTaskID := "test_task_id"
	testAgent.MarkTaskAsBeingProcessed(testTaskID)

	// Проверяем, что задача помечена как обрабатываемая
	_, ok := testAgent.ExecutingLock.Load(testTaskID)
	if !ok {
		t.Errorf("Task with ID %s should be marked as being processed", testTaskID)
	}
}

func TestMarkTaskAsFinished(t *testing.T) {
	// Создаем тестовый агент
	testAgent := &agent.Agent{}

	// Помечаем тестовую задачу как обрабатываемую
	testTaskID := "test_task_id"
	testAgent.MarkTaskAsBeingProcessed(testTaskID)

	// Вызываем функцию MarkTaskAsFinished для тестового taskID
	testAgent.MarkTaskAsFinished(testTaskID)

	// Проверяем, что задача помечена как завершенная
	_, ok := testAgent.ExecutingLock.Load(testTaskID)
	if ok {
		t.Errorf("Task with ID %s should be marked as finished", testTaskID)
	}
}

func TestNewAgent(t *testing.T) {
	// Подготовка
	id := 1
	workers := 3
	postgres := &sqlx.DB{}
	durationMap := map[string]int{"+": 1, "-": 2, "*": 3, "/": 4}

	// Выполнение
	testAgent := agent.NewAgent(id, postgres, workers, durationMap)

	// Проверка
	assert.NotNil(t, testAgent, "The agent should not be nil")
	assert.Equal(t, id, testAgent.ID, "The agent ID should match the provided ID")
	assert.Equal(t, postgres, testAgent.Postgres, "The agent Postgres should match the provided postgres instance")
	assert.Len(t, testAgent.TaskQueues, workers, "The length of TaskQueues should be equal to workers")
	assert.Equal(t, workers, testAgent.Workers, "The number of workers should match the provided workers")
	assert.Equal(t, durationMap, testAgent.DurationMap, "The durationMap should match the provided durationMap")

	// Проверка, что каждый канал в TaskQueues не равен nil
	for i := 0; i < workers; i++ {
		assert.NotNil(t, testAgent.TaskQueues[i], "The TaskQueues[%d] channel should not be nil", i)
	}
}

func TestGetQueueIndex(t *testing.T) {
	// Создаем тестовый агент
	testAgent := &agent.Agent{Workers: 3}

	// Тестовые данные
	testTaskID := "test_task_id"

	// Получаем индекс очереди для задачи
	queueIndex := testAgent.GetQueueIndex(testTaskID)

	// Проверяем, что индекс находится в диапазоне от 0 до (Workers - 1)
	if queueIndex < 0 || queueIndex >= testAgent.Workers {
		t.Errorf("Unexpected queue index: %d", queueIndex)
	}
}

func TestProcessTask(t *testing.T) {
	// Создаем мок базы данных
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	// Создаем экземпляр *sql.DB
	sqlDB := sqlx.NewDb(mockDB, "sqlmock")

	// Создаем экземпляр агента с моком базы данных
	testAgent := &agent.Agent{Postgres: sqlDB}

	// Устанавливаем ожидания для запроса к базе данных
	mock.ExpectExec("UPDATE tasks SET result = (.+), status = (.+) WHERE id = (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Создаем тестовую задачу
	testTask := agent.Task{
		ID:         "test_task_id",
		Expression: "test_expression",
	}

	// Обрабатываем тестовую задачу
	testAgent.ProcessTask(testTask)

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
