package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/while-loop/proxge/pkg"
	"log"
	"time"
)

type redisCache struct {
	client *redis.Client
	exp    time.Duration
}

func NewRedisClient(options *redis.Options) *redis.Client {
	rClient := redis.NewClient(options)

	err := rClient.Ping().Err()
	for err != nil {
		log.Println("unable to ping redis", err, options.Addr)
		time.Sleep(5 * time.Second)
		err = rClient.Ping().Err()
	}
	return rClient
}

func NewRedisCache(options *redis.Options, exp time.Duration) proxge.GECache {
	return &redisCache{
		client: NewRedisClient(options),
		exp:    exp,
	}
}

func (m *redisCache) Get(id int) (int, time.Duration, error) {
	p := m.client.Pipeline()
	priceRes := p.Get(s(id))
	ttlRes := p.TTL(s(id))
	_, err := p.Exec()
	if err != nil {
		return 0, 0, errors.Wrap(err, "redis pipeline failed to get")
	}

	price, _ := priceRes.Int()
	return price, ttlRes.Val(), nil
}

func (m *redisCache) Set(id int, price int) error {
	return m.client.Set(s(id), price, m.exp).Err()
}

func (m *redisCache) GetTTL() time.Duration {
	return m.exp
}

func s(i int) string {
	return fmt.Sprintf("%d", i)
}
