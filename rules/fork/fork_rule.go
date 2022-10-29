package fork

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
)

type ForkRule struct {
}

func NewForkRule() *ForkRule {
	return &ForkRule{}
}

func (f ForkRule) Infinity() bool {
	return false
}

func (f ForkRule) Id() string {
	return "fork"
}

func (f ForkRule) Initialize(ctx context.Context, option map[string]interface{}) error {
	return nil
}

func (f ForkRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {
	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {
		output, err := f.handlerInput(deliver.Context, deliver.Data)
		if err != nil {
			log.Error(err)
			deliver.Reject()
			return
		}

		err = pushQueueOutput(ctx, output)
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

func (f ForkRule) handlerInput(ctx context.Context, input *data.InputEngine) (output *data.OutputEngine, errOutput error) {

	output = data.CreateOutput(input, f.Id())
	var option Option

	err := data.MapToStruct(input.Node.Option, &option)
	if err != nil {
		log.Error(err)

		output.Error = &errors.ErrorEngine{
			Message:     errors.MsgErrorOptionRuleNotMatch,
			ErrorDetail: err,
			FromEngine:  f.Id(),
			Code:        errors.CodeConvert,
			SessionId:   input.SessionId,
			Id:          input.ChanId,
		}
		option.Debug = true
		err = nil
		return
	}
	output.Debug = option.Debug
	output.Data = input.Data
	output.Next = option.NextSuccess
	output.Type = data.TypeOutputEngineSuccess

	return
}

func (f ForkRule) Stop(ctx context.Context) error {
	return nil
}
