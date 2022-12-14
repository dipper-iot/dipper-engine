package log

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/debug"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
)

type LogRule struct {
}

func NewLogRule() *LogRule {
	return &LogRule{}
}

func (l LogRule) Infinity() bool {
	return false
}

func (l LogRule) Id() string {
	return "log-core"
}

func (l LogRule) Initialize(ctx context.Context, option map[string]interface{}) error {

	return nil
}

func (l LogRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {
	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {
		debug.PrintJson(deliver.Data.Data, "Log Data: SessionId -> %d ; ChanId -> %s; From -> %s; Data => ", deliver.Data.SessionId, deliver.Data.ChanId, deliver.Data.ToEngine)
		if deliver.Data.Type == data.TypeOutputEngineError {
			debug.PrintJson(deliver.Data.Error, "Log Error: SessionId -> %d ; ChanId -> %s; From -> %s; Data => ", deliver.Data.SessionId, deliver.Data.ChanId, deliver.Data.ToEngine)
		}

		output := data.CreateOutput(deliver.Data, l.Id())
		output.Next = []string{}
		output.Type = data.TypeOutputEngineSuccess

		err := pushQueueOutput(ctx, output)
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

func (l LogRule) Stop(ctx context.Context) error {

	return nil
}
