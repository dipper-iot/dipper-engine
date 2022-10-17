package fork

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
	"time"
)

type ForkRule struct {
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

	timeData := time.Now()
	output = new(data.OutputEngine)
	output.BranchMain = input.BranchMain
	output.ChanId = input.ChanId
	output.FromEngine = f.Id()
	output.SessionId = input.SessionId
	output.Time = &timeData
	output.Data = input.Data
	output.Type = data.TypeOutputEngineError
	var option Option

	err := util.MapToStruct(input.Node.Option, &option)
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
