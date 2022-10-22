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
	client *redis.Client
	name   string
}

func newRedisQueue[T any](client *redis.Client, id string) *redisQueue[T] {
	return &redisQueue[T]{
		client: client,
		name:   fmt.Sprintf("dipper-queue-%s", id),
	}
}

func FactoryQueueRedis[T any](client *redis.Client) core.FactoryQueue[T] {
	return func(engine core.Rule) queue.QueueEngine[T] {
		return newRedisQueue[T](client, engine.Id())
	}
}

func FactoryQueueNameRedis[T any](client *redis.Client) core.FactoryQueueName[T] {
	return func(name string) queue.QueueEngine[T] {
		return newRedisQueue[T](client, name)
	}
}

func (r redisQueue[T]) Name() string {
	return r.name
}

func (r redisQueue[T]) Publish(ctx context.Context, input T) error {
	data, err := util.ConvertToByte(input)
	if err != nil {
		return err
	}

	return r.client.RPush(ctx, r.name, data).Err()
}

func (r redisQueue[T]) Subscribe(ctx context.Context, callback queue.SubscribeFunction[T]) error {

	go func() {
		for {
			data, err := r.client.RPop(ctx, r.name).Bytes()
			if err == io.EOF {
				return
			}
			if err == redis.Nil {
				continue
			}
			if err != nil {
				log.Error(err)
				return
			}

			var transferData T
			err = json.Unmarshal(data, &transferData)
			if err != nil {
				log.Error(err)
				continue
			}
			timeNow := time.Now()
			callback(queue.NewDeliver[T](
				ctx,
				transferData,
				&timeNow,
				// ack
				func() {

				},
				// reject
				func() {
					go func() {
						// try 3 times
						for i := 0; i < 3; i++ {
							err := r.Publish(ctx, transferData)
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
	}()

	return nil
}
