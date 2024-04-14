package utils

import (
	"fmt"
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

type HandlerFactory struct {
	dbClient        *database.DBClient
	config          *HandlerConfig
	repos           *database.Repositories
	itemRepo        *repository.ItemRepository
	listingRepo     *repository.ListingRepository
	transactionRepo *repository.TransactionRepository
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
		dbClient:        dbClient,
		config:          config,
		repos:           repos,
		itemRepo:        repos.GetItemRepository(),
		listingRepo:     repos.GetListingRepository(),
		transactionRepo: repos.GetTransactionRepository(),
	}
}

func (f *HandlerFactory) NewListingsHandler() *types.Handler {
	return &types.Handler{
		OnResult: func(result interface{}) {
			data := result.(*types.ListingsData)
			// handle item
			item := data.Item
			if item != nil {
				f.itemRepo.UpdateItem(item)
				// save preview url
				previewPath := fmt.Sprintf("%s/%s.png", f.config.staticOutputDir, item.Name)
				DownloadImage(item.IconUrl, previewPath)
			}
			// handle listings
			listings := data.Listings
			f.listingRepo.InsertListings(listings)
		},
		OnError:    OnError,
		OnComplete: OnComplete,
	}
}

func (f *HandlerFactory) NewTransactionHandler() *types.Handler {
	return &types.Handler{
		OnResult: func(result interface{}) {
			data := result.(*types.TransactionData)
			transactions := data.Transactions
			f.transactionRepo.InsertTransactions(transactions)
		},
		OnError:    OnError,
		OnComplete: OnComplete,
	}
}

var DEFAULT_HANDLER_CONFIG = &HandlerConfig{
	staticOutputDir: "output/static",
}
