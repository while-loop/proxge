package cache

import (
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Config struct {
	RedisAddr string        `split_words:"true"`
	RedisPass string        `split_words:"true"`
	RedisDB   int           `split_words:"true" default:"0"`
	Exp       time.Duration `default:"30m"`
}

func TestRedis(t *testing.T) {
	var config Config
	assert.NoError(t, envconfig.Process("proxge", &config))

	o := &redis.Options{
		Addr:     config.RedisAddr,
		DB:       config.RedisDB,
		Password: config.RedisPass,
	}

	rClient := redis.NewClient(o)

	assert.NoError(t, rClient.Ping().Err())
	r := NewRedisCache(rClient, config.Exp)

	assert.NoError(t, r.Set(-3, 300))
	p, err := r.Get(-3)
	assert.NoError(t, err)
	assert.Equal(t, 300, p)
}
