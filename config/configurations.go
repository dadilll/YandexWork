package config

import (
	"github.com/go-redis/redis/v8"
)

// настройка Redis
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

type AppConfig struct {
	NumAgents       int
	WorkersPerAgent int
	DurationMap     map[string]int
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		NumAgents:       3, //настройка кол. агентов
		WorkersPerAgent: 5, //настройка кол. воркеров
		DurationMap: map[string]int{
			"+": 3, // Пример времени задержки для сложения
			"-": 3, // Пример времени задержки для вычитания
			"*": 1, // Пример времени задержки для умножения
			"/": 6, // Пример времени задержки для деления
		},
	}
}
