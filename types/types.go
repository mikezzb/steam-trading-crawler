package types

import "github.com/mikezzb/steam-trading-shared/database/model"

type ListingsData struct {
	Item     *model.Item
	Listings []model.Listing
}

type ItemData struct {
	Item *model.Item
}

type TransactionData struct {
	Transactions []model.Transaction
}
