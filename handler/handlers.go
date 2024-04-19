package handler

import (
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

type IHandlerFactory interface {
	GetListingsHandler() types.Handler
	GetTransactionHandler() types.Handler
}

type HandlerFactory struct {
	config *HandlerConfig
	repos  repository.RepoFactory
}

type HandlerConfig struct {
	StaticOutputDir string
	SecretStore     *shared.PersisitedStore
}

func OnError(err error) {
	log.Printf("Error: %v", err)
}

func OnComplete(result interface{}) {
	log.Println("Complete")
}

func NewHandlerFactory(repos repository.RepoFactory, config *HandlerConfig) *HandlerFactory {
	return &HandlerFactory{
		config: config,
		repos:  repos,
	}
}

func (f *HandlerFactory) GetListingsHandler() types.Handler {
	return NewListingHandler(
		f.repos,
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

var DEFAULT_SECRET_STORE, _ = shared.NewPersisitedStore("../secrets.json")
var DEFAULT_HANDLER_CONFIG = &HandlerConfig{
	StaticOutputDir: "output/static",
	SecretStore:     DEFAULT_SECRET_STORE,
}
