package tonicpow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsInList will test the method isInList()
func TestIsInList(t *testing.T) {
	t.Parallel()

	t.Run("test app list", func(t *testing.T) {
		assert.Equal(t, true, isInList(SortByFieldName, appSortFields))
		assert.Equal(t, true, isInList(SortByFieldCreatedAt, appSortFields))
		assert.Equal(t, false, isInList("bad_field", appSortFields))
	})

	t.Run("test campaign list", func(t *testing.T) {
		assert.Equal(t, true, isInList(SortByFieldBalance, campaignSortFields))
		assert.Equal(t, true, isInList(SortByFieldCreatedAt, campaignSortFields))
		assert.Equal(t, true, isInList(SortByFieldLinksCreated, campaignSortFields))
		assert.Equal(t, true, isInList(SortByFieldPaidClicks, campaignSortFields))
		assert.Equal(t, true, isInList(SortByFieldPayPerClick, campaignSortFields))
		assert.Equal(t, false, isInList("bad_field", campaignSortFields))
	})
}
