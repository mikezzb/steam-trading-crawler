package buff_test

// test
import (
	"net/url"
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/buff"
	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-crawler/handler"
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

func InitBuffCrawler(t *testing.T, cookie string) *buff.BuffCrawler {
	c, err := buff.NewCrawler(cookie)
	if err != nil {
		t.Errorf("Failed to init buff crawler: %v", err)
	}
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

	// repos
	repos := repository.NewRepoFactory(dbClient, nil)

	// handler
	factory := handler.NewHandlerFactory(repos, handler.DEFAULT_HANDLER_CONFIG)
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
		buffCrawler := InitBuffCrawler(t, "")
		// no throttle
		buffCrawler.DoReq("localhost:8000", url.Values{}, "GET")
		// throttled
		buffCrawler.DoReq("localhost:8000", url.Values{}, "GET")
	})
}

func TestStringPrice(t *testing.T) {
	t.Run("StringPrice", func(t *testing.T) {
		testVals := []string{
			"0.01",
			"0.1",
			"1",
			"10",
			"900",
			"23",
			"952340",
			"100000",
		}

		for _, val := range testVals {
			if errors.SafeInvalidPrice < val {
				t.Errorf("invalid price smaller than 0.01")
			}
		}
	})
}
