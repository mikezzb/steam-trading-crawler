package igxe

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

// Implements type.Parser interface
type IgxeParser struct {
}

func toListing(item *IgxeListing) *model.Listing {
	updatedAt, _ := shared.ConvertToUnixTimestamp(item.LastUpdated)
	return &model.Listing{
		Price:      item.Price,
		CreatedAt:  updatedAt,
		UpdatedAt:  updatedAt,
		PreviewUrl: item.InspectImgSmall,

		AssetId:          item.SteamPid,
		TradableCooldown: item.LockSpanStr,
		PaintWear:        item.ExteriorWear,
		PaintIndex:       item.PaintIndex,
		PaintSeed:        item.PaintSeed,
		InstanceId:       strconv.Itoa(item.ID),

		Market: shared.MARKET_NAME_IGXE,
	}
}

func (p *IgxeParser) formatListings(data *IgxeListingResponseData) ([]model.Listing, error) {
	items := data.Listings
	listing := make([]model.Listing, len(items))
	for i, item := range items {
		listing[i] = *toListing(&item)
	}
	return listing, nil
}

func (p *IgxeParser) getPriceItem(name string, listings []model.Listing) (*model.Item, error) {
	now := shared.GetUnixNow()

	igxePrice := model.MarketPrice{
		Price:     utils.ExtractLowestPrice(listings),
		UpdatedAt: now,
	}

	item := &model.Item{
		Name:      name,
		IgxePrice: igxePrice,
	}

	return item, nil
}

func (p *IgxeParser) ParseItemListings(name string, resp *http.Response, resData *IgxeListingResponseData) (*types.ListingsData, error) {
	if resp.StatusCode != http.StatusOK || !resData.Succ {
		return nil, fmt.Errorf("invalid listing response: %d %v", resp.StatusCode, resData)
	}

	// format data
	if listings, err := p.formatListings(resData); err != nil {
		return nil, err
	} else if item, err := p.getPriceItem(name, listings); err != nil {
		return nil, err
	} else {
		utils.PostFormatListings(name, listings)
		return &types.ListingsData{
			Item:     item,
			Listings: listings,
		}, nil
	}
}

func (p *IgxeParser) ParseListingControl(resData *IgxeListingResponseData) *types.CrawlerControl {
	return &types.CrawlerControl{
		TotalPages: resData.Page.PageCount,
	}
}

func toTransactions(item *IgxeTransaction) *model.Transaction {
	timestamp, _ := shared.ConvertChineseDateToUnix(item.LastUpdated)
	return &model.Transaction{
		Price:     item.Price,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,

		PaintWear:  item.ExteriorWear,
		PaintIndex: item.PaintIndex,
		PaintSeed:  item.PaintSeed,
		InstanceId: strconv.Itoa(item.ID),
		// MUST provide an unique asset id to upsert
		AssetId: strconv.Itoa(item.ID),

		Market: shared.MARKET_NAME_IGXE,
	}
}

func (p *IgxeParser) formatTransactions(data *IgxeTransactionResponseData) ([]model.Transaction, error) {
	items := data.Data
	transactions := make([]model.Transaction, len(items))
	for i, item := range items {
		transactions[i] = *toTransactions(&item)
	}
	return transactions, nil
}

func (p *IgxeParser) ParseItemTransactions(name string, resp *http.Response, resData *IgxeTransactionResponseData) (*types.TransactionData, error) {
	if resp.StatusCode != http.StatusOK || !resData.Succ {
		return nil, fmt.Errorf("invalid transaction response: %d %v", resp.StatusCode, resData)
	}

	// format data
	if transactions, err := p.formatTransactions(resData); err != nil {
		return nil, err
	} else {
		utils.PostFormatTransactions(name, transactions)
		return &types.TransactionData{
			Transactions: transactions,
		}, nil
	}
}

func (p *IgxeParser) ParseTransactionControl(resData *IgxeTransactionResponseData) *types.CrawlerControl {
	return nil
}
