package utils

import (
	"push_service/src/core"
	"github.com/go-redis/redis"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var BaseRedis baseRedisModel

type baseRedisModel struct{}

func (baseRedisModel) Connect(Key string) *redis.Client {
	config := core.Config.Redis[Key]

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return redisClient
}

func (baseRedisModel) Close(redisClient *redis.Client) {
	redisClient.Close()
}
