package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Dadil/project/internal/orchestra/api"
	"github.com/Dadil/project/internal/orchestra/domain"
	"github.com/go-redis/redis/v8"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Проверка соединения с Redis
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Создание экземпляра оркестратора
	orchestrator := domain.NewOrchestrator(redisClient)

	// Создание экземпляра HTTP API оркестратора
	orchestratorAPI := api.NewOrchestratorAPI(orchestrator)

	// Запуск HTTP сервера
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port %s...", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, orchestratorAPI.Router))
}
