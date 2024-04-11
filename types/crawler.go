package types

type CrawlerConfig struct {
	MaxItems   int
	Filters    map[string]string
	OnResult   func(interface{})
	OnError    func(error)
	OnComplete func()
}

type Crawler interface {
	Init(secret string) error
	CrawlItemListings(itemName string, config CrawlerConfig) error
	CrawlItemTransactions(itemName string, config CrawlerConfig) error
	GetCookies() (string, error)
}
