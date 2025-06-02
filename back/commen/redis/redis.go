package redis

import (
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	once sync.Once
	rdb  *redis.Client
)

func GetRedis() *redis.Client {
	if rdb == nil {
		once.Do(func() {
			rdb = redis.NewClient(&redis.Options{
				Addr:     "124.70.56.84:6379",
				Password: "320930",
				DB:       0,
			})
		})
	}
	return rdb
}
