package types

type CrawlTaskConfig struct {
	MaxItems int               `json:"maxItems"`
	Filters  map[string]string `json:"filters"`
}

type Crawler interface {
	// RunTask(taskName, itemName string, handler *Handler, config *CrawlerConfig) error
	CrawlItemListings(itemName string, handler Handler, config *CrawlTaskConfig) error
	CrawlItemTransactions(itemName string, handler Handler, config *CrawlTaskConfig) error
	GetCookies() (string, error)
	Stop()
}

type CrawlerControl struct {
	TotalPages int `json:"total_page"`
}
