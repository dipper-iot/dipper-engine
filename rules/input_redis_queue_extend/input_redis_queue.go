package input_redis_queue_extend

import (
	"context"
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/dipper-iot/dipper-engine/rules/common"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

type InputRedisQueueExtendRule struct {
	mapGet sync.Map
}

func (l *InputRedisQueueExtendRule) Infinity() bool {
	return true
}

func NewInputRedisQueueExtendRule() *InputRedisQueueExtendRule {
	return &InputRedisQueueExtendRule{
		mapGet: sync.Map{},
	}
}

func (l *InputRedisQueueExtendRule) ListSession() []uint64 {
	list := make([]uint64, 0)

	l.mapGet.Range(func(key, value any) bool {
		list = append(list, key.(uint64))
		return true
	})

	return list
}

func (l *InputRedisQueueExtendRule) StopSession(id uint64) {
	cancelFn, ok := l.mapGet.Load(id)
	if ok {
		redisInfo, ok := cancelFn.(*RedisSessionInfo)
		if ok {
			redisInfo.cancel()
			redisInfo.client.Close()
		}
	}
}

func (l *InputRedisQueueExtendRule) InfoSession(id uint64) map[string]interface{} {
	return nil
}

func (l *InputRedisQueueExtendRule) Id() string {
	return "input-redis-queue-extend"
}

func (l *InputRedisQueueExtendRule) Initialize(ctx context.Context, optionRaw map[string]interface{}) error {

	return nil
}

func (l *InputRedisQueueExtendRule) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {

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

		client, err := common.ConnectRedis(ctx, &common.OptionRedis{
			Address:  option.RedisAddress,
			Password: option.RedisPassword,
			Db:       option.RedisDb,
		})
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

		_, ok := l.mapGet.Load(deliver.Data.SessionId)
		if !ok {
			go l.getData(ctx2, &OptionLoop{
				option:          &option,
				output:          output,
				pushQueueOutput: pushQueueOutput,
				client:          client,
			})

			l.mapGet.Store(deliver.Data.SessionId, &RedisSessionInfo{
				option: &option,
				client: client,
				cancel: cancel,
			})
		}

		deliver.Ack()
	})
	if err != nil {
		log.Error(err)
		return
	}

}

func (l *InputRedisQueueExtendRule) Stop(ctx context.Context) error {

	l.mapGet.Range(func(key, value any) bool {
		l.StopSession(key.(uint64))
		return true
	})

	return nil

}

func (l *InputRedisQueueExtendRule) getData(ctx context.Context, option *OptionLoop) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				{
					send := option.output.Clone()
					dataRaw, err := option.client.RPop(ctx, option.option.Queue).Bytes()
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

					send.Next = []string{option.option.NextSuccess}
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

func (l *InputRedisQueueExtendRule) sendError(ctx context.Context, e error, message string, send *data.OutputEngine, option *OptionLoop) {
	send.Error = &errors.ErrorEngine{
		ErrorDetail: e,
		Message:     message,
		Code:        errors.CodeProgress,
		SessionId:   send.SessionId,
		FromEngine:  l.Id(),
		Id:          send.IdNode,
	}
	send.Next = []string{option.option.NextError}
	send.Type = data.TypeOutputEngineError
	err := option.pushQueueOutput(ctx, send)
	if err != nil {
		log.Error(err)
	}
}
