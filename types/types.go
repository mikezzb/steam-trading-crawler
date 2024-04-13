package types

import "github.com/mikezzb/steam-trading-shared/database/model"

type ListingsData struct {
	Item     *model.Item
	Listings []model.Listing
}

type TransactionData struct {
	Transactions []model.Transaction
}
