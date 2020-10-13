package cache

import (
	"github.com/while-loop/proxge/pkg"
	"sync"
	"time"
)

type memCache struct {
	data *sync.Map
}

func NewMemCache() proxge.GECache {
	return &memCache{
		data: &sync.Map{},
	}
}

func (m *memCache) Get(id int) (int, time.Duration, error) {
	price, ok := m.data.Load(id)
	if !ok {
		return 0, 0, proxge.ErrDoesNotExist
	}

	return price.(int), time.Duration(0), nil
}

func (m *memCache) Set(id int, price int) error {
	m.data.Store(id, price)
	return nil
}

func (m *memCache) GetTTL() time.Duration {
	return time.Duration(0)
}
