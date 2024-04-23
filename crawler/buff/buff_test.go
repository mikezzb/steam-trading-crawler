package buff_test

// test
import (
	"testing"
	"time"

	"github.com/mikezzb/steam-trading-crawler/crawler/buff"
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
	var secretStore, _ = shared.NewJsonKvStore(
		"../secrets.json",
	)
	buffSecretName := utils.GetSecretName(shared.MARKET_NAME_BUFF)
	buffCrawler := InitBuffCrawler(t, secretStore.Get(buffSecretName).(string))
	defer utils.UpdateSecrets(buffCrawler, secretStore, buffSecretName)

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
	err := buffCrawler.CrawlItemListings(name, listingsHandler, &types.CrawlTaskConfig{
		MaxItems: 20,
	})
	if err != nil {
		t.Errorf("Failed to crawl item listings: %v", err)
	}
}
