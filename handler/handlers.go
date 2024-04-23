package handler

import (
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

type IHandlerFactory interface {
	GetListingsHandler() types.IHandler
	GetTransactionHandler() types.IHandler
}

type HandlerFactory struct {
	config *HandlerConfig
	repos  repository.RepoFactory
}

type HandlerConfig struct {
	StaticOutputDir string
	SecretStore     *shared.JsonKvStore
}

func OnError(err error) {
	log.Printf("Error: %v", err)
}

func OnComplete(result interface{}) {
}

func NewHandlerFactory(repos repository.RepoFactory, config *HandlerConfig) *HandlerFactory {
	return &HandlerFactory{
		config: config,
		repos:  repos,
	}
}

func (f *HandlerFactory) GetListingsHandler() types.IHandler {
	return NewListingHandler(
		f.repos,
		f.config,
	)
}

func (f *HandlerFactory) GetTransactionHandler() types.IHandler {
	return NewBaseHandler(
		func(result interface{}) {
			transactionRepo := f.repos.GetTransactionRepository()
			data := result.(*types.TransactionData)
			transactions := data.Transactions
			err := transactionRepo.UpsertTransactionsByAssetID(transactions)
			if err != nil {
				log.Printf("Failed to upsert transactions: %v", err)
			}
		},
		OnError,
		OnComplete,
	)
}

var DEFAULT_SECRET_STORE, _ = shared.NewJsonKvStore("../secrets.json")
var DEFAULT_HANDLER_CONFIG = &HandlerConfig{
	StaticOutputDir: "output/static",
	SecretStore:     DEFAULT_SECRET_STORE,
}

func GetTestHandler() *BaseHandler {
	return NewBaseHandler(
		func(result interface{}) {
			log.Printf("Result: %v", result)
		},
		func(err error) {
			log.Printf("Error: %v", err)
		},
		func(result interface{}) {
			log.Printf("Complete: %v", result)
		},
	)
}
