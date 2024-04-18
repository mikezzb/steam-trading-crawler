package buff

// test
import (
	"net/url"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/handler"
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
	buffSecretName := utils.GetSecretName(shared.MARKET_NAME_BUFF)
	buffCrawler := InitBuffCrawler(t, secretStore.Get(buffSecretName).(string))
	defer utils.UpdateSecrets(buffCrawler, *secretStore, buffSecretName)

	// db
	dbClient, _ := database.NewDBClient("mongodb://localhost:27017", "steam-trading", 10*time.Second)
	defer dbClient.Disconnect()

	// handler
	factory := handler.NewHandlerFactory(dbClient, handler.DEFAULT_HANDLER_CONFIG)
	listingsHandler := factory.GetListingsHandler()

	// Run
	name := "★ Bayonet | Marble Fade (Factory New)"
	err := buffCrawler.CrawlItemListings(name, listingsHandler, &types.CrawlerConfig{
		MaxItems: 20,
	})
	if err != nil {
		t.Errorf("Failed to crawl item listings: %v", err)
	}
}

func TestBuffSleep(t *testing.T) {
	t.Run("Sleep", func(t *testing.T) {
		buffCrawler := &BuffCrawler{}

		buffCrawler.DoReq("localhost:8000", url.Values{}, "GET")
		buffCrawler.DoReq("localhost:8000", url.Values{}, "GET")
	})
}
