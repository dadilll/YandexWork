package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Dadil/project/config"
	"github.com/Dadil/project/internal/orchestra/api"
	"github.com/Dadil/project/internal/orchestra/domain"
	_ "github.com/lib/pq"
)

func main() {
	// Установка соединения с базой данных PostgreSQL
	postgresDB, err := config.NewPostgreSQLDB()
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL database: %v", err)
	}

	// Проверка соединения с базой данных PostgreSQL
	err = postgresDB.Ping()
	if err != nil {
		log.Fatal("Ошибка при пинге базы данных PostgreSQL:", err)
	}

	// Создание экземпляра Orchestrator с использованием PostgreSQL
	orchestrator := domain.NewOrchestrator(postgresDB)
	api := api.NewOrchestratorAPI(orchestrator)

	// Передача Router в HTTP-сервер
	http.Handle("/", api.Router)

	// Запуск HTTP-сервера
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Запуск сервера на порту %s...", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}
