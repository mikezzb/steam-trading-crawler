package buff

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mikezzb/steam-trading-crawler/types"
	"github.com/mikezzb/steam-trading-crawler/utils"
	shared "github.com/mikezzb/steam-trading-shared"

	"github.com/mikezzb/steam-trading-shared/database/model"
)

type BuffParser struct {
}

type BuffListingResponseData struct {
	Code string `json:"code"`
	Data struct {
		GoodsInfos map[string]BuffGoodsInfo `json:"goods_infos"`
		Items      []BuffItem               `json:"items"`
	} `json:"data"`
}

func itemToListing(item BuffItem) model.Listing {
	return model.Listing{
		Price:            item.Price,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
		PreviewUrl:       item.PreviewUrl,
		GoodsId:          item.GoodsId,
		ClassId:          item.AssetInfo.ClassId,
		AssetId:          item.AssetInfo.AssetId,
		TradableCooldown: item.AssetInfo.TradableCooldown,
		PaintWear:        item.AssetInfo.PaintWear,
		PaintIndex:       item.AssetInfo.Info.PaintIndex,
		PaintSeed:        item.AssetInfo.Info.PaintSeed,
	}
}

func (p *BuffParser) formatItemListings(data BuffListingResponseData) ([]model.Listing, error) {
	items := data.Data.Items
	listing := make([]model.Listing, len(items))
	for i, item := range items {
		listing[i] = itemToListing(item)
	}
	return listing, nil
}

func getFirstValue(data map[string]BuffGoodsInfo) BuffGoodsInfo {
	for _, v := range data {
		return v
	}
	return BuffGoodsInfo{}
}

func ExtractLowestPrice(listing []model.Listing) string {
	lowestPrice := listing[0].Price
	for _, item := range listing {
		if item.Price < lowestPrice {
			lowestPrice = item.Price
		}
	}
	return lowestPrice
}

func (p *BuffParser) formatItem(name string, data BuffListingResponseData, listings []model.Listing) (*model.Item, error) {
	item := getFirstValue(data.Data.GoodsInfos)

	formattedItems := model.Item{
		Name:              name,
		IconUrl:           item.IconUrl,
		SteamPrice:        item.SteamPrice,
		LowestMarketPrice: ExtractLowestPrice(listings),
		LowestMarketName:  shared.MARKET_NAME_BUFF,
	}

	return &formattedItems, nil
}

func (p *BuffParser) ParseItemListings(name string, resp *http.Response, bodyBytes []byte) (*types.ListingsData, error) {
	decodedReader, err := utils.ReadBytes(bodyBytes)
	if err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		return nil, err
	}
	defer decodedReader.Close()

	// unmarshal response
	var data BuffListingResponseData
	if err := json.NewDecoder(decodedReader).Decode(&data); err != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", err)
		return nil, err
	}

	// format data
	if listings, err := p.formatItemListings(data); err != nil {
		return nil, err
	} else if item, err := p.formatItem(name, data, listings); err != nil {
		return nil, err
	} else {
		utils.PostFormatListings(name, listings)
		return &types.ListingsData{
			Item:     item,
			Listings: listings,
		}, nil
	}
}

type BuffTransactionResponseData struct {
	Code string `json:"code"`
	Data struct {
		Items []BuffItem `json:"items"`
	} `json:"data"`
}

func (p *BuffParser) formatItemTransactions(data *BuffTransactionResponseData) ([]model.Transaction, error) {
	items := data.Data.Items
	transactions := make([]model.Transaction, len(items))
	for i, item := range items {
		transactions[i] = model.Transaction(itemToListing(item))
	}
	return transactions, nil
}

func (p *BuffParser) ParseItemTransactions(name string, resp *http.Response, bodyBytes []byte) (*types.TransactionData, error) {
	decodedReader, err := utils.ReadBytes(bodyBytes)
	if err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		return nil, err
	}
	defer decodedReader.Close()

	// unmarshal response
	var data BuffTransactionResponseData
	if err := json.NewDecoder(decodedReader).Decode(&data); err != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", err)
		return nil, err
	}

	if transactions, err := p.formatItemTransactions(&data); err != nil {
		return nil, err
	} else {
		utils.PostFormatTransactions(name, transactions)
		return &types.TransactionData{
			Transactions: transactions,
		}, nil
	}
}
