package main

import (
	"context"
	"log"

	"github.com/Dadil/project/config"
	"github.com/Dadil/project/internal/agent/agent"
)

func main() {

	redisClient := config.NewRedisClient()

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

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
