package arithmetic

import (
	"context"
	"fmt"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
)

type Arithmetic struct {
}

func NewArithmetic() *Arithmetic {
	return &Arithmetic{}
}

func (a Arithmetic) Infinity() bool {
	return false
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

	output = data.CreateOutput(input, a.Id())
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

	mathRunner := NewMath(input.BranchMain, input.Data)

	for keyResult, exp := range option.Operators {
		err = mathRunner.Run(exp, keyResult)
		if err != nil {
			err = fmt.Errorf("%s=%s have error: %s ", keyResult, exp, err.Error())
			break
		}
	}

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
