package main

import (
	"log"

	"github.com/Dadil/project/config"
	"github.com/Dadil/project/internal/agent/agent"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Установка соединения с базой данных PostgreSQL
	postgresDB, err := sqlx.Open("postgres", "user=postgres password=123456789 dbname=calc sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open PostgreSQL database: %v", err)
	}

	// Проверка соединения
	if err = postgresDB.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL database: %v", err)
	}

	log.Println("Connected to PostgreSQL database")

	appConfig := config.NewAppConfig()

	// Создание и запуск агентов
	for i := 1; i <= appConfig.NumAgents; i++ {
		agent := agent.NewAgent(i, postgresDB, appConfig.WorkersPerAgent, appConfig.DurationMap)
		go agent.Start()
	}

	select {}
}
