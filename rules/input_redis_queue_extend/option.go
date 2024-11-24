package input_redis_queue_extend

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/go-redis/redis/v8"
)

type OptionSession struct {
	Queue         string `json:"queue"`
	RedisAddress  string `json:"redis_address"`
	RedisPassword string `json:"redis_password"`
	RedisDb       int    `json:"redis_db"`
	NextSuccess   string `json:"next_success"`
	NextError     string `json:"next_error"`
}

type RedisSessionInfo struct {
	client *redis.Client
	cancel context.CancelFunc
	option *OptionSession
}

type OptionLoop struct {
	pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error
	output          *data.OutputEngine
	option          *OptionSession
	client          *redis.Client
}
