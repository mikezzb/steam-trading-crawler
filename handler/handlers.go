package handler

import (
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-shared/database"
)

type IHandlerFactory interface {
	GetListingsHandler() types.Handler
	GetTransactionHandler() types.Handler
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

func (f *HandlerFactory) GetListingsHandler() types.Handler {
	return NewListingHandler(
		f.repos.GetItemRepository(),
		f.repos.GetListingRepository(),
		f.config,
	)
}

func (f *HandlerFactory) GetTransactionHandler() types.Handler {
	return NewBaseHandler(
		func(result interface{}) {
			transactionRepo := f.repos.GetTransactionRepository()
			data := result.(*types.TransactionData)
			transactions := data.Transactions
			transactionRepo.UpsertTransactionsByAssetID(transactions)
		},
		OnError,
		OnComplete,
	)
}

var DEFAULT_HANDLER_CONFIG = &HandlerConfig{
	staticOutputDir: "output/static",
}
