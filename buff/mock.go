package buff

import "github.com/mikezzb/steam-trading-crawler/types"

type MockBuffCrawler struct {
	parser *BuffParser
}

func (c *MockBuffCrawler) Init() error {
	c.parser = &BuffParser{}
	return nil
}

func (c *MockBuffCrawler) MockParserTransactions(name string, bodyBytes []byte) (*types.TransactionData, error) {
	return c.parser.ParseItemTransactions(name, nil, bodyBytes)
}
