package handler

import (
	"fmt"
	"log"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	"github.com/mikezzb/steam-trading-shared/database/repository"
)

type ListingHandler struct {
	itemRepo    *repository.ItemRepository
	listingRepo *repository.ListingRepository
	config      *HandlerConfig

	// chan
	listingCh chan *types.ListingsData
}

func NewListingHandler(repos repository.RepoFactory, config *HandlerConfig) *ListingHandler {

	handler := &ListingHandler{
		itemRepo:    repos.GetItemRepository(),
		listingRepo: repos.GetListingRepository(),
		config:      config,

		listingCh: make(chan *types.ListingsData, 100),
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
		_, err := h.listingRepo.UpsertListingsByAssetID(listings)
		if err != nil {
			log.Printf("Error: %v", err)
		}
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
