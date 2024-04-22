package types

import (
	"net/http"
)

type CrawlTaskConfig struct {
	MaxItems int               `json:"maxItems"`
	Filters  map[string]string `json:"filters"`
}

type ICrawler interface {
	CrawlItemListings(itemName string, handler IHandler, config *CrawlTaskConfig) error
	CrawlItemTransactions(itemName string, handler IHandler, config *CrawlTaskConfig) error
	GetCookies() (string, error)
	Stop()
}

type CrawlerControl struct {
	TotalPages int `json:"total_page"`
}

type IParser interface {
	ParseItemListings(name string, resp *http.Response) (*ListingsData, error)
	ParseItemTransactions(name string, resp *http.Response) (*TransactionData, error)
}
