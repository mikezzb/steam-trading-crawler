package types

type CrawlerConfig struct {
	MaxItems int               `json:"maxItems"`
	Filters  map[string]string `json:"filters"`
}

type Crawler interface {
	Init(secret string) error
	// RunTask(taskName, itemName string, handler *Handler, config *CrawlerConfig) error
	CrawlItemListings(itemName string, handler Handler, config *CrawlerConfig) error
	CrawlItemTransactions(itemName string, handler Handler, config *CrawlerConfig) error
	GetCookies() (string, error)
}
