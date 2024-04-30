package buff

import (
	"fmt"
	"net/http"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"
)

type BuffParser struct {
}

// Shared by listing and transaction
func toListing(item *BuffItem) *types.Listing {
	return &types.Listing{
		Price:            item.Price,
		CreatedAt:        shared.UnixToTime(item.CreatedAt),
		UpdatedAt:        shared.UnixToTime(item.UpdatedAt),
		PreviewUrl:       item.PreviewUrl,
		GoodsId:          item.GoodsId,
		ClassId:          item.AssetInfo.ClassId,
		AssetId:          item.AssetInfo.AssetId,
		TradableCooldown: item.AssetInfo.TradableCooldown,
		PaintWear:        item.AssetInfo.PaintWear,
		PaintIndex:       item.AssetInfo.Info.PaintIndex,
		PaintSeed:        item.AssetInfo.Info.PaintSeed,
		InstanceId:       item.AssetInfo.InstanceId,
		Market:           shared.MARKET_NAME_BUFF,
	}
}

func (p *BuffParser) formatItemListings(data *BuffListingResponseData) ([]types.Listing, error) {
	items := data.Data.Items
	listing := make([]types.Listing, len(items))
	for i, item := range items {
		listing[i] = *toListing(&item)
	}
	return listing, nil
}

func getFirstValue(data map[string]BuffGoodsInfo) BuffGoodsInfo {
	for _, v := range data {
		return v
	}
	return BuffGoodsInfo{}
}

func (p *BuffParser) formatItem(name string, data *BuffListingResponseData, listings []types.Listing) (*types.Item, error) {
	item := getFirstValue(data.Data.GoodsInfos)

	now := shared.GetNow()

	steamPrice := &types.MarketPrice{
		Price:     item.SteamPrice,
		UpdatedAt: now,
	}

	buffPrice := &types.MarketPrice{
		Price:     utils.ExtractLowestPrice(listings),
		UpdatedAt: now,
	}

	formattedItems := types.Item{
		Name:       name,
		IconUrl:    item.IconUrl,
		SteamPrice: steamPrice,
		BuffPrice:  buffPrice,
	}

	return &formattedItems, nil
}

func (p *BuffParser) ParseItemListings(name string, resp *http.Response, resData *BuffListingResponseData) (*types.ListingsData, error) {
	if resp.StatusCode != http.StatusOK || resData.Code != "OK" {
		return nil, fmt.Errorf("invalid response: %d %s", resp.StatusCode, resData.Code)
	}

	// format data
	if listings, err := p.formatItemListings(resData); err != nil {
		return nil, err
	} else if item, err := p.formatItem(name, resData, listings); err != nil {
		return nil, err
	} else {
		utils.PostFormatListings(name, listings)
		return &types.ListingsData{
			Item:     item,
			Listings: listings,
		}, nil
	}
}

func (p *BuffParser) toTransaction(item *BuffItem) *types.Transaction {
	return &types.Transaction{
		AssetId:          item.AssetInfo.AssetId,
		Market:           shared.MARKET_NAME_BUFF,
		Price:            item.Price,
		CreatedAt:        shared.UnixToTime(item.CreatedAt),
		PreviewUrl:       item.PreviewUrl,
		GoodsId:          item.GoodsId,
		ClassId:          item.AssetInfo.ClassId,
		TradableCooldown: item.AssetInfo.TradableCooldown,
		PaintWear:        item.AssetInfo.PaintWear,
		PaintIndex:       item.AssetInfo.Info.PaintIndex,
		PaintSeed:        item.AssetInfo.Info.PaintSeed,
		InstanceId:       item.AssetInfo.InstanceId,
	}
}

func (p *BuffParser) formatItemTransactions(data *BuffTransactionResponseData) ([]types.Transaction, error) {
	items := data.Data.Items
	transactions := make([]types.Transaction, len(items))
	for i, item := range items {
		transactions[i] = *p.toTransaction(&item)
	}
	return transactions, nil
}

func (p *BuffParser) ParseItemTransactions(name string, resp *http.Response, resData *BuffTransactionResponseData) (*types.TransactionData, error) {
	if resp.StatusCode != http.StatusOK || resData.Code != "OK" {
		return nil, fmt.Errorf("invalid response: %d %s", resp.StatusCode, resData.Code)
	}

	if transactions, err := p.formatItemTransactions(resData); err != nil {
		return nil, err
	} else {
		utils.PostFormatTransactions(name, transactions)

		return &types.TransactionData{
			Transactions: transactions,
		}, nil
	}
}

func (p *BuffParser) ParseListingControl(resData *BuffListingResponseData) *types.CrawlerControl {
	return &types.CrawlerControl{
		TotalPages: resData.Data.TotalPages,
	}
}

func (p *BuffParser) ParseTransactionControl(resData *BuffTransactionResponseData) *types.CrawlerControl {
	return &types.CrawlerControl{
		TotalPages: resData.Data.TotalPages,
	}
}
