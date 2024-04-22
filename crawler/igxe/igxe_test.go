package igxe_test

import (
	"testing"

	"github.com/mikezzb/steam-trading-crawler/crawler/igxe"
	"github.com/mikezzb/steam-trading-crawler/handler"
	"github.com/mikezzb/steam-trading-crawler/types"
)

func TestIgxeListing(t *testing.T) {
	c, err := igxe.NewCrawler("")
	if err != nil {
		t.Errorf("Failed to init igxe crawler: %v", err)
	}

	t.Run("IgxeListing", func(t *testing.T) {
		handler := handler.GetTestHandler()
		// test
		c.CrawlItemListings(
			"★ Bayonet | Marble Fade (Factory New)",
			handler,
			&types.CrawlTaskConfig{
				MaxItems: 10,
			},
		)
	})
}

func TestIgxeTransaction(t *testing.T) {
	c, err := igxe.NewCrawler("")
	if err != nil {
		t.Errorf("Failed to init igxe crawler: %v", err)
	}

	t.Run("IgxeTransactions", func(t *testing.T) {
		handler := handler.GetTestHandler()
		// test
		c.CrawlItemTransactions(
			"★ Bayonet | Marble Fade (Factory New)",
			handler,
			&types.CrawlTaskConfig{
				MaxItems: 10,
			},
		)
	})
}
