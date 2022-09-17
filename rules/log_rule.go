package rules

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
)

type LogRule struct {
}

func (l LogRule) Id() string {
	return "logger"
}

func (l LogRule) Initialize(ctx context.Context, option map[string]interface{}) error {

	return nil
}

func (l LogRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

}

func (l LogRule) Stop(ctx context.Context) error {

	return nil
}
