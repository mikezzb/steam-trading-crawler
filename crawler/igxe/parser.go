package igxe

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

// Implements type.Parser interface
type IgxeParser struct {
}

func toListing(item *IgxeListing) *types.Listing {
	updatedAt, _ := shared.ParseDateHhmmss(item.LastUpdated)
	return &types.Listing{
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

func (p *IgxeParser) formatListings(data *IgxeListingResponseData) ([]types.Listing, error) {
	items := data.Listings
	listing := make([]types.Listing, len(items))
	for i, item := range items {
		listing[i] = *toListing(&item)
	}
	return listing, nil
}

func (p *IgxeParser) getPriceItem(name string, listings []types.Listing) (*types.Item, error) {
	igxePrice := &types.MarketPrice{
		Price:     utils.ExtractLowestPrice(listings),
		UpdatedAt: shared.GetNow(),
	}

	item := &types.Item{
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

func toTransactions(item *IgxeTransaction) *types.Transaction {
	timestamp, _ := shared.ParseChineseDate(item.LastUpdated)
	return &types.Transaction{
		AssetId: strconv.Itoa(item.ID),
		Market:  shared.MARKET_NAME_IGXE,

		Price:     item.Price,
		CreatedAt: timestamp,

		PaintWear:  item.ExteriorWear,
		PaintIndex: item.PaintIndex,
		PaintSeed:  item.PaintSeed,
		InstanceId: strconv.Itoa(item.ID),
	}
}

func (p *IgxeParser) formatTransactions(data *IgxeTransactionResponseData) ([]types.Transaction, error) {
	items := data.Data
	transactions := make([]types.Transaction, len(items))
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
