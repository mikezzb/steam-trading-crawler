package utils_test

import (
	"testing"

	"github.com/mikezzb/steam-trading-crawler/types"
	shared "github.com/mikezzb/steam-trading-shared"
)

func TestGetTier(t *testing.T) {
	t.Run("GetTier", func(t *testing.T) {
		testPairs := []struct {
			listing  types.Listing
			expected string
		}{
			{
				types.Listing{Name: "★ Flip Knife | Marble Fade (Factory New)", PaintSeed: 872},
				"Tricolor",
			},
			{
				types.Listing{Name: "★ Karambit | Doppler (Factory New)", PaintSeed: 741},
				"Good Phase 2",
			},
			{
				types.Listing{Name: "★ Bayonet | Marble Fade (Factory New)", PaintSeed: 727},
				"FFI",
			},
		}

		for _, pair := range testPairs {
			listing := pair.listing
			actual := shared.GetTier(listing.Name, listing.PaintSeed)
			if actual != pair.expected {
				t.Errorf("Expected %s, got %s", pair.expected, actual)
			}
		}
	})
}
