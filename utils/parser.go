package utils

import (
	"github.com/mikezzb/steam-trading-crawler/errors"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

// extract the lowest price from for each listing page, update lowest price at crawler
func ExtractLowestPrice(listing []model.Listing) string {
	if len(listing) == 0 {
		return errors.SafeInvalidPrice
	}

	lowestPrice := listing[0].Price
	for _, item := range listing {
		if item.Price < lowestPrice {
			lowestPrice = item.Price
		}
	}
	return lowestPrice
}
