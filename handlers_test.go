package crawler_test

import (
	"testing"
	"time"

	crawler "github.com/mikezzb/steam-trading-crawler"

	"github.com/mikezzb/steam-trading-shared/database"
)

func TestHandlerFactory_NewItemHandler(t *testing.T) {
	t.Run("Factory", func(t *testing.T) {
		dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		factory := crawler.NewHandlerFactory(dbClient)
		itemHandler := factory.NewItemHandler()
		if itemHandler == nil {
			t.Error("ItemHandler is nil")
		}

		listingsHandler := factory.NewListingsHandler()
		if listingsHandler == nil {
			t.Error("ListingsHandler is nil")
		}

		defer dbClient.Disconnect()

	})
}
