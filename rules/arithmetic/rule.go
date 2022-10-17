package arithmetic

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
	"time"
)

type Arithmetic struct {
}

func (a Arithmetic) Id() string {
	return "arithmetic"
}

func (a Arithmetic) Initialize(ctx context.Context, options map[string]interface{}) error {
	return nil
}

func (a Arithmetic) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

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

func (a Arithmetic) Stop(ctx context.Context) error {
	return nil
}

func (a Arithmetic) handlerInput(ctx context.Context, input *data.InputEngine) (output *data.OutputEngine, errOutput error) {

	output = a.createOutput(input)
	var option Option

	err := util.MapToStruct(input.Node.Option, &option)
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

	mathRunner := NewMath(input.BranchMain, input.Data)

	err = mathRunner.Run(option.List)
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
	output.Next = []string{option.NextSuccess}
	output.Type = data.TypeOutputEngineSuccess

	return
}

func (a Arithmetic) createOutput(input *data.InputEngine) (output *data.OutputEngine) {

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
