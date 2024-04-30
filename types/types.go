package types

type ListingsData struct {
	Item     *Item
	Listings []Listing
}

type ItemData struct {
	Item *Item
}

type TransactionData struct {
	Transactions []Transaction
}
