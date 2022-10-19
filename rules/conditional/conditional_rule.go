package conditional

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
	"time"
)

type ConditionalRule struct {
}

func (a ConditionalRule) Id() string {
	return "conditional"
}

func (a ConditionalRule) Initialize(ctx context.Context, options map[string]interface{}) error {
	return nil
}

func (a ConditionalRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {
		output, err := a.handlerInput(deliver.Context, deliver.Data)
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

func (a ConditionalRule) Stop(ctx context.Context) error {
	return nil
}

func (a ConditionalRule) handlerInput(ctx context.Context, input *data.InputEngine) (output *data.OutputEngine, errOutput error) {

	output = a.createOutput(input)
	var option Option

	err := data.MapToStruct(input.Node.Option, &option)
	if err != nil {
		log.Error(err)

		output.Error = &errors.ErrorEngine{
			Message:     errors.MsgErrorOptionRuleNotMatch,
			ErrorDetail: err,
			FromEngine:  a.Id(),
			Code:        errors.CodeConvert,
			SessionId:   input.SessionId,
			Id:          input.ChanId,
		}
		option.Debug = true
		err = nil
		return
	}
	output.Next = []string{option.NextError}
	output.Debug = option.Debug

	mathRunner := NewConditional(input.BranchMain, input.Data)

	var result bool
	result, err = mathRunner.Run(option.Operator, option.SetParamResultTo)
	if err != nil {
		log.Errorf("Run Math error -> %s", err.Error())

		output.Error = &errors.ErrorEngine{
			Message:     errors.MsgErrorServerErrorProgress,
			ErrorDetail: err,
			FromEngine:  a.Id(),
			Code:        errors.CodeProgress,
			SessionId:   input.SessionId,
			Id:          input.ChanId,
		}
		err = nil
		return
	}

	output.Data = mathRunner.Data()
	if result {
		output.Next = []string{option.NextTrue}
	} else {
		output.Next = []string{option.NextFalse}
	}
	output.Type = data.TypeOutputEngineSuccess

	return
}

func (a ConditionalRule) createOutput(input *data.InputEngine) (output *data.OutputEngine) {

	timeData := time.Now()

	output = new(data.OutputEngine)
	output.BranchMain = input.BranchMain
	output.ChanId = input.ChanId
	output.FromEngine = a.Id()
	output.SessionId = input.SessionId
	output.Time = &timeData
	output.Data = input.Data
	output.Type = data.TypeOutputEngineError

	return
}
