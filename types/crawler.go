package types

type CrawlerConfig struct {
	MaxItems int
	Filters  map[string]string
}

type Crawler interface {
	Init(secret string) error
	CrawlItemListings(itemName string, handler Handler, config CrawlerConfig) error
	CrawlItemTransactions(itemName string, handler Handler, config CrawlerConfig) error
	GetCookies() (string, error)
}
