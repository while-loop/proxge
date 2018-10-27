package cache

import (
	"fmt"
	"github.com/go-redis/redis"
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

func NewRedisCache(client *redis.Client, exp time.Duration) proxge.GECache {
	return &redisCache{
		client: client,
		exp:    exp,
	}
}

func (m *redisCache) Get(id int) (int, error) {
	return m.client.Get(s(id)).Int()
}

func (m *redisCache) Set(id int, price int) error {
	return m.client.Set(s(id), price, m.exp).Err()
}

func s(i int) string {
	return fmt.Sprintf("%d", i)
}
