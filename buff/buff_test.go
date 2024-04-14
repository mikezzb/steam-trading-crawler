package buff

// test
import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database"
)

func InitBuffCrawler(t *testing.T, cookie string) *BuffCrawler {
	c := &BuffCrawler{}
	c.Init(cookie)
	return c
}

func TestBuffCrawler_CrawlListings(t *testing.T) {
	// Init
	var secretStore, _ = shared.NewPersisitedStore(
		"../secrets.json",
	)
	buffCrawler := InitBuffCrawler(t, secretStore.Get("buff_secret").(string))
	defer utils.UpdateSecrets(buffCrawler, *secretStore, "buff_secret")

	// db
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", 10*time.Second)
	defer dbClient.Disconnect()

	// handler
	factory := utils.NewHandlerFactory(dbClient, utils.DEFAULT_HANDLER_CONFIG)
	listingsHandler := factory.NewListingsHandler()

	// Run
	name := "★ Bayonet | Marble Fade (Factory New)"
	err := buffCrawler.CrawlItemListings(name, listingsHandler, &types.CrawlerConfig{
		MaxItems: 20,
	})
	if err != nil {
		t.Errorf("Failed to crawl item listings: %v", err)
	}
}

func TestBuffParser_ParseItemListings(t *testing.T) {
	name := "★ Karambit | Marble Fade (Factory New)"
	testCases := []struct {
		mockResJsonPath string
	}{
		{"mocks/gzip_encode.json"},
		// {"output/buff_l_★ M9 Bayonet | Fade (Factory New)_1_1712820647114.json"},
	}

	for _, tc := range testCases {
		mockResJsonPath := tc.mockResJsonPath
		mockRes, err := os.ReadFile(mockResJsonPath)
		if err != nil {
			t.Errorf("Failed to read mock response file: %s", mockResJsonPath)
		}

		p := &BuffParser{}
		data, err := p.ParseItemListings(name, nil, mockRes)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Item: %v", data.Item)
		t.Logf("Listings: %v", data.Listings)

		// save item and listings to JSON files
		itemJsonPath := "mocks/item.json"
		if err := utils.WriteJSONToFile(data.Item, itemJsonPath); err != nil {
			t.Errorf("Failed to write item JSON to file: %v", err)
		}
		listingsJsonPath := "mocks/listings.json"
		if err := utils.WriteJSONToFile(data.Listings, listingsJsonPath); err != nil {
			t.Errorf("Failed to write listings JSON to file: %v", err)
		}
	}
}

func TestBuffCrawler(t *testing.T) {
	// Init
	var secretStore, _ = shared.NewPersisitedStore(
		"../secrets.json",
	)
	buffCrawler := InitBuffCrawler(t, secretStore.Get("buff_secret").(string))
	defer utils.UpdateSecrets(buffCrawler, *secretStore, "buff_secret")

	// db
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", 10*time.Second)
	defer dbClient.Disconnect()

	// handler
	factory := utils.NewHandlerFactory(dbClient, utils.DEFAULT_HANDLER_CONFIG)
	handler := factory.NewTransactionHandler()

	t.Run("CrawlItemTransactions", func(t *testing.T) {
		name := "★ Bayonet | Marble Fade (Factory New)"
		err := buffCrawler.CrawlItemTransactions(name, handler, &types.CrawlerConfig{})
		if err != nil {
			t.Errorf("Failed to crawl item transactions: %v", err)
		}
	})

}

func TestBuffParser_ParseTransactions(t *testing.T) {
	name := "★ Bayonet | Marble Fade (Factory New)"
	testCases := []struct {
		mockResJsonPath string
	}{
		{"mocks/transactions.json"},
	}

	id := shared.GetBuffIds()["★ Bayonet | Marble Fade (Factory New)"]
	fmt.Printf("ID: %v\n", id)

	// DB & REPO
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading-unit-test", 10*time.Second)
	defer dbClient.Disconnect()

	factory := utils.NewHandlerFactory(dbClient, utils.DEFAULT_HANDLER_CONFIG)
	handler := factory.NewTransactionHandler()

	for _, tc := range testCases {
		mockResJsonPath := tc.mockResJsonPath
		mockRes, err := os.ReadFile(mockResJsonPath)
		if err != nil {
			t.Errorf("Failed to read mock response file: %s", mockResJsonPath)
		}

		// p := &BuffParser{}
		mp := &MockBuffCrawler{}
		mp.Init()

		data, err := mp.MockParserTransactions(name, mockRes)
		if err != nil {
			t.Error(err)
		}

		fmt.Printf("Transactions: %v", data)

		// TEST DB
		handler.OnResult(data)
	}
}
