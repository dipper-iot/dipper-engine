package input_redis_queue

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
)

type Option struct {
	RedisAddress  string `json:"redis_address"`
	RedisPassword string `json:"redis_password"`
	RedisDb       int    `json:"redis_db"`
}

type OptionSession struct {
	Queue       string `json:"queue"`
	NextSuccess string `json:"next_success"`
	NextError   string `json:"next_error"`
}

type OptionLoop struct {
	pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error
	queueName       string
	output          *data.OutputEngine
	nextSuccess     string
	nextError       string
}
