package cache

import (
	"github.com/while-loop/proxge/pkg"
	"sync"
)

type memCache struct {
	data *sync.Map
}

func NewMemCache() proxge.GECache {
	return &memCache{
		data: &sync.Map{},
	}
}

func (m *memCache) Get(id int) (int, error) {
	price, ok := m.data.Load(id)
	if !ok {
		return 0, proxge.ErrDoesNotExist
	}

	return price.(int), nil
}

func (m *memCache) Set(id int, price int) error {
	m.data.Store(id, price)
	return nil
}
