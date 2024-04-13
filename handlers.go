package crawler

import (
	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

type HandlerFactory struct {
	dbClient        *database.DBClient
	repos           *database.Repositories
	itemRepo        *repository.ItemRepository
	listingRepo     *repository.ListingRepository
	transactionRepo *repository.TransactionRepository
}

func NewHandlerFactory(dbClient *database.DBClient) *HandlerFactory {
	repos := database.NewRepositories(dbClient)
	return &HandlerFactory{
		dbClient:        dbClient,
		repos:           repos,
		itemRepo:        repos.GetItemRepository(),
		listingRepo:     repos.GetListingRepository(),
		transactionRepo: repos.GetTransactionRepository(),
	}
}

func (f *HandlerFactory) NewItemHandler() *types.Handler {
	return &types.Handler{
		OnResult: func(result interface{}) {
			// save to db
		},
		OnError: func(err error) {
		},
		OnComplete: func() {
		}}
}

func (f *HandlerFactory) NewListingsHandler() *types.Handler {
	return &types.Handler{
		OnResult: func(result interface{}) {
			// save to db

			// save preview urls
		},
		OnError: func(err error) {
		},
		OnComplete: func() {
		}}
}
