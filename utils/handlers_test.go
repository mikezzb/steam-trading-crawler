package utils_test

import (
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/utils"

	"github.com/mikezzb/steam-trading-shared/database"
)

func TestHandlerFactory_NewItemHandler(t *testing.T) {
	t.Run("Factory", func(t *testing.T) {
		dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", time.Second*10)
		factory := utils.NewHandlerFactory(dbClient, utils.DEFAULT_HANDLER_CONFIG)

		listingsHandler := factory.GetListingsHandler()
		if listingsHandler == nil {
			t.Error("ListingsHandler is nil")
		}

		defer dbClient.Disconnect()

	})
}
