package common

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis(ctx context.Context, option *OptionRedis) (client *redis.Client, err error) {

	client = redis.NewClient(&redis.Options{
		Addr:     option.Address,
		Password: option.Password,
		DB:       option.Db,
	})

	err = client.Ping(ctx).Err()
	if err == nil {
		return
	}

	return
}
