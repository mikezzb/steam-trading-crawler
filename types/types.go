package types

import "steam-trading/shared/database/model"

type ListingsData struct {
	Item     *model.Item
	Listings []model.Listing
}

type TransactionData struct {
	Transactions []model.Transaction
}
