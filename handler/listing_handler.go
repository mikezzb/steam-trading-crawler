package handler

import (
	"fmt"
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database"
	"github.com/mikezzb/steam-trading-shared/database/repository"
	"github.com/mikezzb/steam-trading-shared/subscription"
)

type ListingHandler struct {
	itemRepo    *repository.ItemRepository
	listingRepo *repository.ListingRepository
	config      *HandlerConfig

	// chan
	listingCh chan *types.ListingsData

	// sub emitter
	subEmitter *subscription.NotificationEmitter
}

func NewListingHandler(repos *database.Repositories, config *HandlerConfig) *ListingHandler {

	handler := &ListingHandler{
		itemRepo:    repos.GetItemRepository(),
		listingRepo: repos.GetListingRepository(),
		config:      config,

		listingCh: make(chan *types.ListingsData, 100),
		subEmitter: subscription.NewNotificationEmitter(
			repos.GetSubscriptionRepository(),
			repos.GetItemRepository(),
			&subscription.NotifierConfig{
				TelegramToken: config.SecretStore.Get(shared.SECRET_TELEGRAM_TOKEN).(string),
			},
		),
	}

	go handler.onResult()

	return handler
}

func (h *ListingHandler) Close() {
	close(h.listingCh)
}

// Remove the unused method processUpdatedListing()
func (h *ListingHandler) onResult() {
	for data := range h.listingCh {
		// handle item
		item := data.Item
		if item != nil {
			h.itemRepo.UpdateItem(item)
			// save preview url
			previewPath := fmt.Sprintf("%s/%s.png", h.config.StaticOutputDir, item.Name)
			utils.DownloadImage(item.IconUrl, previewPath)
		}
		// handle listings
		listings := data.Listings
		updatedListings, err := h.listingRepo.UpsertListingsByAssetID(listings)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		// notify subscribers for updated / created listings
		h.subEmitter.EmitListings(updatedListings)
	}
}

func (h *ListingHandler) OnResult(result interface{}) {
	data := result.(*types.ListingsData)
	h.listingCh <- data
}

func (h *ListingHandler) OnError(err error) {
	log.Printf("Error: %v", err)
}

func (h *ListingHandler) OnComplete() {
}
