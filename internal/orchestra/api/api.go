package api

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Dadil/project/internal/orchestra/domain"
	"github.com/gorilla/mux"
)

type expressionRequest struct {
	Expression string `json:"expression"`
}

type OrchestratorAPI struct {
	Router       *mux.Router
	Orchestrator *domain.Orchestrator
}

func NewOrchestratorAPI(orchestrator *domain.Orchestrator) *OrchestratorAPI {
	api := &OrchestratorAPI{
		Router:       mux.NewRouter(),
		Orchestrator: orchestrator,
	}

	api.setupRoutes()
	return api
}

func (api *OrchestratorAPI) setupRoutes() {
	api.Router.HandleFunc("/add", api.AddExpression).Methods("POST")
	api.Router.HandleFunc("/expressions", api.GetExpressions).Methods("GET")
	api.Router.HandleFunc("/delete-all", api.DeleteAllTasks).Methods("DELETE")
}

func (api *OrchestratorAPI) DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	// Удаляем все задачи из хранилища
	err := api.Orchestrator.Redis.FlushAll(context.Background()).Err()
	if err != nil {
		http.Error(w, "Failed to delete tasks", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All tasks deleted successfully"))
}

func (api *OrchestratorAPI) AddExpression(w http.ResponseWriter, r *http.Request) {
	var expressionRequest expressionRequest
	err := json.NewDecoder(r.Body).Decode(&expressionRequest)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверка на валидность выражения
	if !isValidExpression(expressionRequest.Expression) {
		http.Error(w, "Expression is invalid", http.StatusBadRequest)
		return
	}

	id, err := api.Orchestrator.AddTask(expressionRequest.Expression)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция для проверки валидности выражения
func isValidExpression(expr string) bool {
	// Регулярное выражение для проверки допустимости выражения
	validExpression := regexp.MustCompile(`^[\d+\-*/\s]+$`)
	return validExpression.MatchString(expr)
}

func (api *OrchestratorAPI) GetExpressions(w http.ResponseWriter, r *http.Request) {
	expressions := api.Orchestrator.GetTasks()
	var tasks []*domain.Task
	for _, expr := range expressions {
		// Создаем новый экземпляр Task для каждой задачи
		task := expr
		tasks = append(tasks, &task)
	}
	jsonResponse(w, tasks)
}

func jsonResponse(w http.ResponseWriter, data []*domain.Task) {
	w.Header().Set("Content-Type", "application/json")
	for _, task := range data {
		err := json.NewEncoder(w).Encode(task)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
