package ge

import (
	"github.com/stretchr/testify/assert"
	"github.com/while-loop/proxge/pkg/cache"
	"testing"
)

func TestCannonBall(t *testing.T) {
	c := cache.NewMemCache()
	rsbuddy := NewRsBuddyGe(c).(*rsBuddyGe)

	price, err := rsbuddy.PriceById(2)
	assert.NoError(t, err)
	assert.True(t, price > 0)
}
