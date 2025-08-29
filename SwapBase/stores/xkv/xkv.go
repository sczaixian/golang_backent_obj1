package xkv

import (
	"log"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Store struct {
	kv.Store
	Redis *redis.Redis
}

func NewStore(c kv.KvConf) *Store {
	if len(c) == 0 || cache.TotalWeights(c) <= 0{
		log.Fatal("no cache nodes")
	}

	cn := redis.MustNewRedis(c[0].RedisConf)
	return &Store{Store: kv.NewStore(c), Redis: cn}
}