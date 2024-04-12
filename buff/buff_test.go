package buff

// test
import (
	"os"
	"steam-trading/shared"
	"testing"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
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

	// Run
	name := "★ Karambit | Marble Fade (Factory New)"
	err := buffCrawler.CrawlItemListings(name, types.Handler{
		OnResult: func(result interface{}) {
			t.Logf("Result: %v", result)
		},
		OnError: func(err error) {
			t.Errorf("Error: %v", err)
		},
	}, types.CrawlerConfig{
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
		data, err := p.ParseItemListings(name, mockRes)
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
