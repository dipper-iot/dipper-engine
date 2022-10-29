package output_redis_queue

import (
	"context"
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/dipper-iot/dipper-engine/rules/common"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
)

type OutputRedisQueueRule struct {
	client *redis.Client
}

func (l *OutputRedisQueueRule) Infinity() bool {
	return false
}

func NewOutputRedisQueueRule() *OutputRedisQueueRule {
	return &OutputRedisQueueRule{}
}
func (l *OutputRedisQueueRule) Id() string {
	return "output-redis-queue"
}

func (l *OutputRedisQueueRule) Initialize(ctx context.Context, optionRaw map[string]interface{}) error {

	var option Option
	err := data.MapToStruct(optionRaw, &option)
	if err != nil {
		log.Error(err)
		return err
	}

	l.client, err = common.ConnectRedis(ctx, &common.OptionRedis{
		Address:  option.RedisAddress,
		Password: option.RedisPassword,
		Db:       option.RedisDb,
	})
	if err == nil {
		return err
	}

	return nil
}

func (l *OutputRedisQueueRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {

		output := data.CreateOutput(deliver.Data, l.Id())

		var option OptionSession
		err := data.MapToStruct(deliver.Data.Node.Option, &option)
		if err != nil {
			log.Error(err)
			output.Error = &errors.ErrorEngine{
				Message:     errors.MsgErrorOptionRuleNotMatch,
				ErrorDetail: err,
				FromEngine:  l.Id(),
				Code:        errors.CodeConvert,
				SessionId:   deliver.Data.SessionId,
				Id:          deliver.Data.ChanId,
			}
			output.Debug = deliver.Data.Node.Debug

			pushQueueOutput(ctx, output)
			err = nil
			return
		}

		dataByte, err := json.Marshal(deliver.Data.Data)
		if err != nil {
			log.Error(err)
			l.sendError(ctx, err, "Redis POP unmarshal error", output, &option, pushQueueOutput)
			return
		}

		err = l.client.RPush(ctx, option.Queue, dataByte).Err()
		if err != nil {
			log.Error(err)
			l.sendError(ctx, err, "Redis RPush error", output, &option, pushQueueOutput)
			return
		}

		output.Next = []string{option.NextSuccess}
		output.Type = data.TypeOutputEngineSuccess
		output.Data = deliver.Data.Data
		err = pushQueueOutput(ctx, output)
		if err != nil {
			log.Error(err)
		}

		deliver.Ack()
	})
	if err != nil {
		log.Error(err)
		return
	}

}

func (l *OutputRedisQueueRule) Stop(ctx context.Context) error {
	return l.client.Close()
}

func (l *OutputRedisQueueRule) sendError(ctx context.Context, e error, message string, send *data.OutputEngine, option *OptionSession, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {
	send.Error = &errors.ErrorEngine{
		ErrorDetail: e,
		Message:     message,
		Code:        errors.CodeProgress,
		SessionId:   send.SessionId,
		FromEngine:  l.Id(),
		Id:          send.IdNode,
	}
	send.Next = []string{option.NextSuccess}
	send.Type = data.TypeOutputEngineError
	err := pushQueueOutput(ctx, send)
	if err != nil {
		log.Error(err)
	}
}
