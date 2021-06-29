package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/profiralex/go-bootstrap-redis/pkg/config"
	"sync"
)

var rdb *redis.Client
var rdbOnce sync.Once

func Init(cfg config.Config) {
	rdbOnce.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
		})
	})
}

func getRedisClient() *redis.Client {
	return rdb
}
