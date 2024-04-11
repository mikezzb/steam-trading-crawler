package buff

// test
import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mikezzb/steam-trading-crawler/utils"
)

func InitBuffCrawler(t *testing.T) *BuffCrawler {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Errorf("Error loading .env file: %v", err)
	}

	c := &BuffCrawler{}
	c.Init(os.Getenv("BUFF_SECRET"))
	return c
}

func TestBuffCrawler_CrawlListings(t *testing.T) {
	name := "★ M9 Bayonet | Marble Fade (Factory New)"
	buffCrawler := InitBuffCrawler(t)
	err := buffCrawler.CrawlItemListings(name, 10, nil)
	if err != nil {
		t.Errorf("Failed to crawl item listings: %v", err)
	}
}

func TestBuffParser_ParseItemListings(t *testing.T) {
	name := "★ Karambit | Marble Fade (Factory New)"
	mockResJsonPath := "mocks/gzip_encode.json"
	mockRes, err := os.ReadFile(mockResJsonPath)
	if err != nil {
		t.Errorf("Failed to read mock response file: %s", mockResJsonPath)
	}

	// convert mock response to http.Response
	mockResReader := bytes.NewReader(mockRes)
	mockResHttp := &http.Response{
		Body: io.NopCloser(mockResReader),
	}

	p := &BuffParser{}
	item, listings, err := p.ParseItemListings(name, mockResHttp)
	if err != nil {
		t.Errorf("Failed to parse item listings: %v", err)
	}

	// save item and listings to JSON files
	itemJsonPath := "mocks/item.json"
	if err := utils.WriteJSONToFile(item, itemJsonPath); err != nil {
		t.Errorf("Failed to write item JSON to file: %v", err)
	}
	listingsJsonPath := "mocks/listings.json"
	if err := utils.WriteJSONToFile(listings, listingsJsonPath); err != nil {
		t.Errorf("Failed to write listings JSON to file: %v", err)
	}

}
