package handler_test

import (
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/handler"

	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

func TestHandlerFactory_NewItemHandler(t *testing.T) {
	t.Run("Factory", func(t *testing.T) {
		dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		repos := repository.NewRepoFactory(dbClient, nil)
		factory := handler.NewHandlerFactory(repos, handler.DEFAULT_HANDLER_CONFIG)

		listingsHandler := factory.GetListingsHandler()
		if listingsHandler == nil {
			t.Error("ListingsHandler is nil")
		}

		defer dbClient.Disconnect()

	})
}
