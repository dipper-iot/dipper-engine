package log

import (
	"context"
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
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
	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {
		defer deliver.Ack()
		dataByte, err := json.Marshal(deliver.Data.Data)
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("Log Data: SessionId -> %d ; ChanId -> %s; Data -> %s", deliver.Data.SessionId, deliver.Data.ChanId, string(dataByte))
		if deliver.Data.Type == data.TypeOutputEngineError {
			dataByte, err := json.Marshal(deliver.Data.Error)
			if err != nil {
				log.Error(err)
				return
			}
			log.Infof("Log Error: SessionId -> %d ; ChanId -> %s; Data -> %s", deliver.Data.SessionId, deliver.Data.ChanId, string(dataByte))
		}
	})
	if err != nil {
		log.Error(err)
		return
	}

}

func (l LogRule) Stop(ctx context.Context) error {

	return nil
}
