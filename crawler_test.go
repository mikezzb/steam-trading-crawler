package main_test

import (
	"log"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/buff"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

func TestBuff_CrawlTransactions(t *testing.T) {
	t.Run("CrawlTransactions", func(t *testing.T) {

		// Init
		var secretStore, _ = shared.NewPersisitedStore(
			"secrets.json",
		)
		buffCrawler := &buff.BuffCrawler{}
		buffCrawler.Init(secretStore.Get("buff_secret").(string))
		defer utils.UpdateSecrets(buffCrawler, *secretStore, "buff_secret")

		// db
		dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", 10*time.Second)
		defer dbClient.Disconnect()

		// handler
		factory := utils.NewHandlerFactory(dbClient, utils.DEFAULT_HANDLER_CONFIG)
		handler := factory.GetTransactionHandler()

		t.Run("CrawlItemTransactions", func(t *testing.T) {
			name := "★ M9 Bayonet | Marble Fade (Factory New)"
			err := buffCrawler.CrawlItemTransactions(name, handler, &types.CrawlerConfig{})
			if err != nil {
				t.Errorf("Failed to crawl item transactions: %v", err)
			}
		})
	})
}

func TestPostProcessors(t *testing.T) {
	t.Run("PostFormatTransactions", func(t *testing.T) {
		name := "★ Bayonet | Marble Fade (Factory New)"
		transactions := []model.Transaction{
			{
				PaintSeed: 727,
			},
		}

		utils.PostFormatTransactions(name, transactions)

		if transactions[0].Name != "★ Bayonet | Marble Fade (Factory New)" {
			t.Error("Failed to format transaction name")
		}

		if transactions[0].Rarity != "FFI" {
			t.Error("Failed to format rarity")
		}

		log.Printf("%+v\n", transactions[0])
	})
}
