package utils_test

import (
	"testing"

	shared "github.com/mikezzb/steam-trading-shared"
	"github.com/mikezzb/steam-trading-shared/database/model"
)

func TestGetTier(t *testing.T) {
	t.Run("GetTier", func(t *testing.T) {
		testPairs := []struct {
			listing  model.Listing
			expected string
		}{
			{
				model.Listing{Name: "★ Flip Knife | Marble Fade (Factory New)", PaintSeed: 872},
				"Tricolor",
			},
			{
				model.Listing{Name: "★ Karambit | Doppler (Factory New)", PaintSeed: 741},
				"Good Phase 2",
			},
			{
				model.Listing{Name: "★ Bayonet | Marble Fade (Factory New)", PaintSeed: 727},
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
