package interfaces

import (
	"net/http"
	"steam-trading/shared"
)

type Parser interface {
	ParseItemListings(name string, resp *http.Response) (*shared.Item, []shared.Listing, error)
	ParseItemTransactions(name string, resp *http.Response) ([]shared.Transaction, error)
}
