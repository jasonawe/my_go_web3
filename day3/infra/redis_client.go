package infra

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedisClient() {
	log.Info("asdasd", TheAppConfig.RedisAddr)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     TheAppConfig.RedisAddr,
		Password: "",
		DB:       0,
	})
}
