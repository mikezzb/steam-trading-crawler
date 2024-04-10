package buff

// test
import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestBuffCrawler_DoReq(t *testing.T) {
	c := &BuffCrawler{}
	c.Init("test")
	c.DoReq("http://localhost:8080", nil, "GET")
}

func TestBuffParser_ParseItemListings(t *testing.T) {
	name := "★ Karambit | Marble Fade (Factory New)"
	mockResJsonPath := "mocks/listing_res.json"
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
	fmt.Printf("Item: %+v\n", item)
	fmt.Printf("Listings: %+v\n", listings)
}
