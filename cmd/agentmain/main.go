package main

import (
	"log"

	"github.com/Dadil/project/config"
	"github.com/Dadil/project/internal/agent/agent"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	postgresDB, err := config.NewPostgreSQLDB()
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL database: %v", err)
	}

	err = postgresDB.Ping()
	if err != nil {
		log.Fatal("Ошибка при пинге базы данных PostgreSQL:", err)
	}

	log.Println("Connected to PostgreSQL database")

	// Convert postgresDB to *sqlx.DB
	postgresDBx := sqlx.NewDb(postgresDB, "postgres")

	appConfig := config.NewAppConfig()

	// Создание и запуск агентов
	for i := 1; i <= appConfig.NumAgents; i++ {
		agent := agent.NewAgent(i, postgresDBx, appConfig.WorkersPerAgent, appConfig.DurationMap)
		go agent.Start()
	}

	select {}
}
