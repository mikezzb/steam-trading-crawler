package utils

import (
	"steam-trading/shared"
)

func PostFormatListing(name string, listings []shared.Listing) []shared.Listing {
	for i := range listings {
		// add name to listings
		listings[i].Name = name
		listings[i].Rarity = shared.GetListingTier(listings[i])
	}
	return listings
}
