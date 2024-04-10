package interfaces

import (
	"net/http"
	"steam-trading/shared"
)

type Parser interface {
	ParseItemListings(resp *http.Response) ([]shared.Item, error)
	ParseItemTransactions(resp *http.Response) ([]shared.Transaction, error)
}
