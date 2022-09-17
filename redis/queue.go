package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

type redisQueue[T any] struct {
	client      *redis.Client
	dataDefault T
	name        string
}

func newRedisQueue[T any](client *redis.Client, id string, dataDefault T) *redisQueue[T] {
	return &redisQueue[T]{
		client:      client,
		dataDefault: dataDefault,
		name:        fmt.Sprintf("dipper-queue-%s", id),
	}
}

func FactoryQueueRedis[T any](client *redis.Client, defaultQueue T) core.FactoryQueue[T] {
	return func(engine core.Rule) queue.QueueEngine[T] {
		return newRedisQueue[T](client, engine.Id(), defaultQueue)
	}
}

func FactoryQueueNameRedis[T any](client *redis.Client, defaultQueue T) core.FactoryQueueName[T] {
	return func(name string) queue.QueueEngine[T] {
		return newRedisQueue[T](client, name, defaultQueue)
	}
}

func (r redisQueue[T]) Name() string {
	return r.name
}

func (r redisQueue[T]) Pushlish(ctx context.Context, input T) error {
	data, err := util.ConvertToByte(input)
	if err != nil {
		return err
	}

	return r.client.RPush(ctx, r.name, data).Err()
}

func (r redisQueue[T]) Subscribe(ctx context.Context, callback queue.SubscribeFunction[T]) error {

	go func(dataQueue T) {
		for {
			data, err := r.client.RPop(ctx, r.name).Bytes()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Error(err)
				return
			}

			err = json.Unmarshal(data, dataQueue)
			if err != nil {
				log.Error(err)
				continue
			}
			timeNow := time.Now()
			callback(queue.NewDeliver[T](
				ctx,
				dataQueue,
				&timeNow,
				// ack
				func() {

				},
				// reject
				func() {
					go func() {
						// try 3 times
						for i := 0; i < 3; i++ {
							err := r.Pushlish(ctx, dataQueue)
							if err != nil {
								log.Error(err)
								continue
							}
							break
						}
					}()
				},
			))
		}
	}(r.dataDefault)

	return nil
}
