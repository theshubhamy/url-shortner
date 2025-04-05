package database

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/theshubhamy/urlshortner/config"
)

var Ctx = context.Background()

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.DbUrl,
		Password: config.DbPass,
		DB:       dbNo,
	})
	return rdb
}
