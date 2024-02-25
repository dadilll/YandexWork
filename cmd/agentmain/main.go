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

	appConfig := config.NewAppConfig()

	// Создание и запуск агентов
	for i := 1; i <= appConfig.NumAgents; i++ {
		agent := agent.NewAgent(i, redisClient, appConfig.WorkersPerAgent, appConfig.DurationMap)
		go agent.Start()
	}

	select {}
}
