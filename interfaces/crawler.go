package interfaces

type Crawler interface {
	Init(secret string) error
	CrawlItemListings(itemName string, maxListings int) error
	CrawlItemTransactions(itemName string, maxTransactions int) error
}
