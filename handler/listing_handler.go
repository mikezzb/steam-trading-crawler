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
	itemCh    chan *types.ItemData
}

func NewListingHandler(repos repository.RepoFactory, config *HandlerConfig) *ListingHandler {

	handler := &ListingHandler{
		itemRepo:    repos.GetItemRepository(),
		listingRepo: repos.GetListingRepository(),
		config:      config,

		listingCh: make(chan *types.ListingsData, 100),
		itemCh:    make(chan *types.ItemData, 100),
	}

	go handler.onListingResult()
	go handler.OnItemResult()

	return handler
}

func (h *ListingHandler) Close() {
	close(h.listingCh)
}

// Remove the unused method processUpdatedListing()
func (h *ListingHandler) onListingResult() {
	for data := range h.listingCh {
		listings := data.Listings
		_, err := h.listingRepo.UpsertListingsByAssetID(listings)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func (h *ListingHandler) OnItemResult() {
	for data := range h.itemCh {
		item := data.Item
		if item != nil {
			h.itemRepo.UpdateItem(item)
			// save preview url
			previewPath := fmt.Sprintf("%s/%s.png", h.config.StaticOutputDir, item.Name)
			utils.DownloadImage(item.IconUrl, previewPath)
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

func (h *ListingHandler) OnComplete(result interface{}) {
	// for the listing handler, the on complete will returns the updated item
	data := result.(*types.ItemData)
	h.itemCh <- data
}
