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
	"github.com/mikezzb/steam-trading-shared/database/repository"
	"github.com/mikezzb/steam-trading-shared/subscription"
)

func main() {
	// secret store
	secretStore, err := shared.NewJsonKvStore("secrets.json")

	if err != nil {
		log.Fatalf("Failed to load secrets: %v", err)
		return
	}

	// db
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", time.Second*10)
	defer dbClient.Disconnect()

	// event emitter
	emitter := subscription.NewNotificationEmitter(&subscription.NotifierConfig{
		TelegramToken: secretStore.Get(shared.SECRET_TELEGRAM_TOKEN).(string),
	})

	// change stream handlers
	changeStreamHandler := &repository.ChangeStreamHandlers{
		ListingChangeStreamCallback:      emitter.ListingChangeStreamHandler,
		SubscriptionChangeStreamCallback: emitter.SubChangeStreamHandler,
	}

	// repos
	repos := repository.NewRepoFactory(dbClient, changeStreamHandler)

	// init event emitter
	emitter.Init(repos)

	// handlers
	handlerFactory := handler.NewHandlerFactory(repos, &handler.HandlerConfig{
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
