package redis

import (
	"github.com/go-redis/redis"
	"os"

	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
)

func NewClientFromEnv() *redis.Client {
	// RedisClient
	addr := os.Getenv("REDIS_ADDR")
	passwd := os.Getenv("REDIS_PASSWD")
	if addr == "" {
		addr = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       0,
	})
}

type RedisStorage struct {
	*redis.Client
}

func init() {
	// RedisStorage
	storage.GlobalStorage = &RedisStorage{
		Client: NewClientFromEnv(),
	}
}
