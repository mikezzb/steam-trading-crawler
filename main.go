package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mikezzb/steam-trading-crawler/handler"
	"github.com/mikezzb/steam-trading-crawler/runner"
	"github.com/mikezzb/steam-trading-crawler/types"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database"
)

func main() {
	// secret store
	secretStore, err := shared.NewPersisitedStore("secrets.json")

	if err != nil {
		log.Fatalf("Failed to load secrets: %v", err)
		return
	}

	// db
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", time.Second*10)
	defer dbClient.Disconnect()

	// handlers
	handlerFactory := handler.NewHandlerFactory(dbClient, &handler.HandlerConfig{
		SecretStore:     secretStore,
		StaticOutputDir: "output/static",
	})

	// tasks

	fileBytes, err := os.ReadFile("crawler_tasks.json")

	if err != nil {
		log.Fatalf("Failed to load tasks: %v", err)
		return
	}

	var tasks types.CrawlerTasks

	json.Unmarshal(fileBytes, &tasks)

	runner, err := runner.NewRunner(&runner.RunnerConfig{
		LogFolder:      "logs",
		SecretStore:    secretStore,
		HandlerFactory: handlerFactory,
		MaxReruns:      4,
	})

	if err != nil {
		log.Printf("Failed to create runner: %v", err)
		return
	}

	runner.Run(tasks.Tasks)
}
