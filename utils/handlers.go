package utils

import (
	"fmt"
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-shared/database"
)

type HandlerFactoryInterface interface {
	GetListingsHandler() *types.Handler
	GetTransactionHandler() *types.Handler
}

type HandlerFactory struct {
	dbClient *database.DBClient
	config   *HandlerConfig
	repos    *database.Repositories
}

type HandlerConfig struct {
	staticOutputDir string
}

func OnError(err error) {
	log.Printf("Error: %v", err)
}

func OnComplete() {
	log.Println("Complete")
}

func NewHandlerFactory(dbClient *database.DBClient, config *HandlerConfig) *HandlerFactory {
	repos := database.NewRepositories(dbClient)
	return &HandlerFactory{
		dbClient: dbClient,
		config:   config,
		repos:    repos,
	}
}

func (f *HandlerFactory) GetListingsHandler() *types.Handler {
	return &types.Handler{
		OnResult: func(result interface{}) {
			itemRepo := f.repos.GetItemRepository()
			listingRepo := f.repos.GetListingRepository()
			data := result.(*types.ListingsData)
			// handle item
			item := data.Item
			if item != nil {
				itemRepo.UpdateItem(item)
				// save preview url
				previewPath := fmt.Sprintf("%s/%s.png", f.config.staticOutputDir, item.Name)
				DownloadImage(item.IconUrl, previewPath)
			}
			// handle listings
			listings := data.Listings
			listingRepo.InsertListings(listings)
		},
		OnError:    OnError,
		OnComplete: OnComplete,
	}
}

func (f *HandlerFactory) GetTransactionHandler() *types.Handler {
	return &types.Handler{
		OnResult: func(result interface{}) {
			transactionRepo := f.repos.GetTransactionRepository()
			data := result.(*types.TransactionData)
			transactions := data.Transactions
			transactionRepo.InsertTransactions(transactions)
		},
		OnError:    OnError,
		OnComplete: OnComplete,
	}
}

var DEFAULT_HANDLER_CONFIG = &HandlerConfig{
	staticOutputDir: "output/static",
}
