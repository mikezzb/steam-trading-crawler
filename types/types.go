package types

import "steam-trading/shared"

type ListingsData struct {
	Item     *shared.Item
	Listings []shared.Listing
}

type TransactionData struct {
	Transactions []shared.Transaction
}
