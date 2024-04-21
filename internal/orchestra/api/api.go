package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Dadil/project/internal/orchestra/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type expressionRequest struct {
	Expression string `json:"expression"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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
	api.Router.HandleFunc("/register", api.RegisterUser).Methods("POST")
	api.Router.HandleFunc("/login", api.LoginUser).Methods("POST")
	api.Router.HandleFunc("/add", api.AddExpression).Methods("POST")
	api.Router.HandleFunc("/expressions", api.GetExpressions).Methods("GET")
	api.Router.HandleFunc("/delete-tasks", api.DeleteAllTasksForUser).Methods("DELETE")
}

func (api *OrchestratorAPI) DeleteAllTasksForUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete all tasks for user")

	authHeader := r.Header.Get("Authorization")
	if _, err := ValidateJWTTokenFromHeader(authHeader); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	login, err := ValidateJWTTokenFromHeader(authHeader)
	if err != nil {
		log.Println("Error validating JWT token:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = api.Orchestrator.DeleteAllTasksForUser(login)
	if err != nil {
		log.Println("Error deleting tasks for user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("All tasks deleted successfully for user:", login)
	response := map[string]string{"message": "All tasks deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *OrchestratorAPI) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Received registration request")

	var registerRequest User
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = api.Orchestrator.CreateUser(registerRequest.Login, registerRequest.Password)
	if err != nil {
		log.Println("Failed to register user:", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	log.Println("User registered successfully:", registerRequest.Login)
	response := map[string]string{"message": "User registered successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *OrchestratorAPI) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Received login request")

	var loginRequest User
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hashedPassword, err := api.GetUserHashByLogin(loginRequest.Login)
	if err != nil {
		log.Println("Error retrieving user hash:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ComparePassword(loginRequest.Password, hashedPassword)
	if err != nil {
		log.Println("Incorrect password:", err)
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	log.Println("User logged in successfully:", loginRequest.Login)

	// Генерируем JWT токен для залогинившегося пользователя
	tokenString, err := GenerateJWTToken(loginRequest.Login)
	if err != nil {
		log.Println("Error generating JWT token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Включаем токен в ответ
	response := map[string]string{
		"message": "User logged in successfully",
		"token":   tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Compare(hash string, s string) error {
	existing := []byte(hash)
	incoming := []byte(s)
	log.Printf("Comparing password hash: %s with: %s\n", existing, incoming) // Добавленный лог
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func GenerateJWTToken(login string) (string, error) {
	// Устанавливаем время истечения срока действия токена на 24 часа от текущего момента
	expirationTime := time.Now().Add(24 * time.Hour)

	// Создаем JWT токен с указанием времени истечения срока действия и логина пользователя в качестве утверждения
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   expirationTime.Unix(), // Установка времени истечения срока действия
	})

	// Подписываем токен
	tokenString, err := token.SignedString([]byte("DbJ8nHvYX5QhragK8T5G4vsGxOgkorBtqBPMbPhPAw4="))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ComparePassword(password string, hashedPassword string) error {
	existing := []byte(hashedPassword)
	incoming := []byte(password)
	log.Printf("Comparing password hash: %s with: %s\n", existing, incoming) // Добавленный лог
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func (api *OrchestratorAPI) GetUserHashByLogin(login string) (string, error) {
	user, err := api.Orchestrator.GetUserByLogin(login)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	return user.Password, nil
}

func ValidateJWTTokenFromHeader(header string) (string, error) {
	tokenString := extractTokenFromHeader(header)
	if tokenString == "" {
		return "", fmt.Errorf("no token found in Authorization header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte("DbJ8nHvYX5QhragK8T5G4vsGxOgkorBtqBPMbPhPAw4="), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	login, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return login, nil
}

func (api *OrchestratorAPI) GetExpressions(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add expression")

	authHeader := r.Header.Get("Authorization")
	if _, err := ValidateJWTTokenFromHeader(authHeader); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	login, err := ValidateJWTTokenFromHeader(authHeader)
	if err != nil {
		log.Println("Error validating JWT token:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Продолжаем выполнение запроса
	expressions := api.Orchestrator.GetTasksForUser(login)
	var tasks []*domain.Task
	for _, expr := range expressions {
		task := expr
		tasks = append(tasks, &task)
	}
	jsonResponse(w, tasks)
}

func (api *OrchestratorAPI) AddExpression(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to add expression")

	authHeader := r.Header.Get("Authorization")
	if _, err := ValidateJWTTokenFromHeader(authHeader); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var expressionRequest expressionRequest
	err := json.NewDecoder(r.Body).Decode(&expressionRequest)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !IsValidExpression(expressionRequest.Expression) {
		log.Println("Invalid expression:", expressionRequest.Expression)
		http.Error(w, "Expression is invalid", http.StatusBadRequest)
		return
	}

	login, err := ValidateJWTTokenFromHeader(authHeader)
	if err != nil {
		log.Println("Error validating JWT token:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	existingTasks := api.Orchestrator.GetTasksForUser(login)
	for _, task := range existingTasks {
		if task.Expression == expressionRequest.Expression {
			log.Println("Task with the same expression already exists")
			http.Error(w, "Task with the same expression already exists", http.StatusBadRequest)
			return
		}
	}

	id, err := api.Orchestrator.AddTaskForUser(expressionRequest.Expression, login)
	if err != nil {
		log.Println("Error adding task:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func IsValidExpression(expr string) bool {
	validExpression := regexp.MustCompile(`^[\d+\-*/\s]+$`)
	return validExpression.MatchString(expr)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func extractTokenFromHeader(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
