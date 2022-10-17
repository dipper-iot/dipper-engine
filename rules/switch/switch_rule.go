package _switch

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
	"time"
)

type SwitchRule struct {
}

func (a SwitchRule) Id() string {
	return "switch"
}

func (a SwitchRule) Initialize(ctx context.Context, options map[string]interface{}) error {
	return nil
}

func (a SwitchRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

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

func (a SwitchRule) Stop(ctx context.Context) error {
	return nil
}

func (a SwitchRule) handlerInput(ctx context.Context, input *data.InputEngine) (output *data.OutputEngine, errOutput error) {

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

	mathRunner := NewSwitch(input.BranchMain, input.Data)

	var result string
	result, err = mathRunner.Run(option.MapSwitch, option.Key)
	if err != nil {
		log.Errorf("Run Switch error -> %s", err.Error())

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

	output.Data = input.Data
	output.Next = []string{result}
	output.Type = data.TypeOutputEngineSuccess

	return
}

func (a SwitchRule) createOutput(input *data.InputEngine) (output *data.OutputEngine) {

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
