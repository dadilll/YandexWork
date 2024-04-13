package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Функция для создания нового подключения к базе данных PostgreSQL
func NewPostgreSQLDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=123456789 dbname=calc sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Установите максимальное количество открытых соединений
	db.SetMaxOpenConns(10)
	// Установите максимальное время жизни соединения
	db.SetConnMaxLifetime(time.Minute * 20)

	return db, nil
}

type AppConfig struct {
	NumAgents       int
	WorkersPerAgent int
	DurationMap     map[string]int
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		NumAgents:       3, // Настройка количества агентов
		WorkersPerAgent: 5, // Настройка количества воркеров
		DurationMap: map[string]int{
			"+": 10, // Пример времени задержки для сложения
			"-": 10, // Пример времени задержки для вычитания
			"*": 10, // Пример времени задержки для умножения
			"/": 10, // Пример времени задержки для деления
		},
	}
}
