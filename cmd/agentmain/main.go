package main

import (
	"context"
	"log"

	"github.com/Dadil/project/internal/agent/agent"
	"github.com/go-redis/redis/v8"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	// Если мы дошли до этой точки, подключение к Redis прошло успешно
	log.Println("Connected to Redis")
	// Создание и запуск агентов
	numAgents := 3
	workersPerAgent := 5 // количество воркеров на каждого агента
	for i := 1; i <= numAgents; i++ {
		agent := agent.NewAgent(i, redisClient, workersPerAgent)
		go agent.Start()
	}

	select {}
}
