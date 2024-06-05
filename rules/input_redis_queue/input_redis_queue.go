package input_redis_queue

import (
	"context"
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/dipper-iot/dipper-engine/rules/common"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

type InputRedisQueueRule struct {
	client *redis.Client
	mapGet sync.Map
}

func (l *InputRedisQueueRule) Infinity() bool {
	return true
}

func NewInputRedisQueueRule() *InputRedisQueueRule {
	return &InputRedisQueueRule{
		mapGet: sync.Map{},
	}
}

func (l *InputRedisQueueRule) ListSession() []uint64 {
	list := make([]uint64, 0)

	l.mapGet.Range(func(key, value any) bool {
		list = append(list, key.(uint64))
		return true
	})

	return list
}

func (l *InputRedisQueueRule) StopSession(id uint64) {
	cancelFn, ok := l.mapGet.Load(id)
	if ok {
		cancel, ok := cancelFn.(context.CancelFunc)
		if ok {
			cancel()
		}
	}
}

func (l *InputRedisQueueRule) InfoSession(id uint64) map[string]interface{} {
	return nil
}

func (l *InputRedisQueueRule) Id() string {
	return "input-redis-queue"
}

func (l *InputRedisQueueRule) Initialize(ctx context.Context, optionRaw map[string]interface{}) error {

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

func (l *InputRedisQueueRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

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

		ctx2, cancel := context.WithCancel(ctx)

		_, ok := l.mapGet.Load(deliver.Data.SessionId)
		if !ok {
			go l.getData(ctx2, &OptionLoop{
				nextSuccess:     option.NextSuccess,
				nextError:       option.NextError,
				output:          output,
				pushQueueOutput: pushQueueOutput,
				queueName:       option.Queue,
			})

			l.mapGet.Store(deliver.Data.SessionId, cancel)
		}

		deliver.Ack()
	})
	if err != nil {
		log.Error(err)
		return
	}

}

func (l *InputRedisQueueRule) Stop(ctx context.Context) error {
	return l.client.Close()
}

func (l *InputRedisQueueRule) getData(ctx context.Context, option *OptionLoop) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				{
					send := option.output.Clone()
					dataRaw, err := l.client.RPop(ctx, option.queueName).Bytes()
					if err == io.EOF {
						return
					}
					if err == redis.Nil {
						continue
					}

					if err != nil {
						log.Error(err)
						l.sendError(ctx, err, "Redis POP message error", send, option)
						continue
					}

					var transferData map[string]interface{}
					err = json.Unmarshal(dataRaw, &transferData)
					if err != nil {
						log.Error(err)
						l.sendError(ctx, err, "Redis POP unmarshal error", send, option)
						continue
					}

					send.Next = []string{option.nextSuccess}
					send.Type = data.TypeOutputEngineSuccess
					send.Data = transferData
					err = option.pushQueueOutput(ctx, send)
					if err != nil {
						log.Error(err)
					}
					continue
				}
			}
		}
	}()
}

func (l *InputRedisQueueRule) sendError(ctx context.Context, e error, message string, send *data.OutputEngine, option *OptionLoop) {
	send.Error = &errors.ErrorEngine{
		ErrorDetail: e,
		Message:     message,
		Code:        errors.CodeProgress,
		SessionId:   send.SessionId,
		FromEngine:  l.Id(),
		Id:          send.IdNode,
	}
	send.Next = []string{option.nextError}
	send.Type = data.TypeOutputEngineError
	err := option.pushQueueOutput(ctx, send)
	if err != nil {
		log.Error(err)
	}
}
