package rules

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
)

type ArithmeticRule struct {
}

func (a ArithmeticRule) Id() string {
	return "arithmetic"
}

func (a ArithmeticRule) Initialize(ctx context.Context, option map[string]interface{}) error {
	return nil
}

func (a ArithmeticRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {
		err := a.handlerInput(deliver.Context, deliver.Data)
		if err != nil {
			log.Error(err)
			deliver.Reject()
			return
		}

		deliver.Ack()
	})
	if err != nil {
		log.Error(err)
		return
	}

}

func (a ArithmeticRule) Stop(ctx context.Context) error {
	return nil
}

func (a ArithmeticRule) handlerInput(ctx context.Context, input *data.InputEngine) error {

	return nil
}
