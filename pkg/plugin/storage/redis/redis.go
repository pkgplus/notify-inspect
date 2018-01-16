package redis

import (
	"github.com/go-redis/redis"

	"github.com/xuebing1110/notify-inspect/pkg/log"
	"github.com/xuebing1110/notify-inspect/pkg/plugin/storage"
	myredis "github.com/xuebing1110/notify-inspect/pkg/redis"
)

type RedisStorage struct {
	*redis.Client
	log.Logger
}

func init() {
	// RedisStorage
	storage.GlobalStorage = &RedisStorage{
		Logger: log.GlobalLogger,
		Client: myredis.GetClient(),
	}
}
