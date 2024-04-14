package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/mikezzb/steam-trading-crawler/runner"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	"github.com/mikezzb/steam-trading-shared/database"
)

func main() {
	// db
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", time.Second*10)
	defer dbClient.Disconnect()

	// handlers
	handlerFactory := utils.NewHandlerFactory(dbClient, utils.DEFAULT_HANDLER_CONFIG)

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
		SecretPath:     "secrets.json",
		HandlerFactory: handlerFactory,
		MaxReruns:      4,
	})

	if err != nil {
		log.Printf("Failed to create runner: %v", err)
		return
	}

	runner.Run(tasks.Tasks)
}
