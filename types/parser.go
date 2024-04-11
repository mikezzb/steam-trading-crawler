package types

import (
	"net/http"
)

type Parser interface {
	ParseItemListings(name string, resp *http.Response) (*ListingsData, error)
	ParseItemTransactions(name string, resp *http.Response) (*TransactionData, error)
}
