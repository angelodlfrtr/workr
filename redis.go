package workr

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// RedisConnect connect to redis and ping it
func (wrkr *Workr) RedisConnect() error {
	if wrkr.RedisClient != nil {
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", wrkr.Config.RedisHost, wrkr.Config.RedisPort),
		Password: wrkr.Config.RedisPassword,
		DB:       wrkr.Config.RedisDB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	wrkr.RedisClient = rdb

	return err
}

// NewZ returns a redis Z type
func NewZ(score float64, member interface{}) *redis.Z {
	return &redis.Z{
		Score:  score,
		Member: member,
	}
}
