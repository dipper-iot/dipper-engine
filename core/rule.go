package core

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
)

type Rule interface {
	Id() string
	Initialize(ctx context.Context, option map[string]interface{}) error
	Run(ctx context.Context,
		subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error,
		pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error)
	Stop(ctx context.Context) error
}
