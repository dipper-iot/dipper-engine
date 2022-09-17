package redis

import (
	"context"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"io"
)

type redisBus struct {
	client *redis.Client
}

func NewRedisBus(client *redis.Client) *redisBus {
	return &redisBus{client: client}
}

func (r redisBus) Pushlish(ctx context.Context, name string, input interface{}) error {
	data, err := util.ConvertToByte(input)
	if err != nil {
		return err
	}

	return r.client.Publish(ctx, name, data).Err()
}

func (r redisBus) Subscribe(ctx context.Context, name string, callback func(data []byte)) error {
	pubsub := r.client.Subscribe(ctx, name)

	// Close the subscription when we are done.
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Error(err)
			continue
		}

		callback([]byte(msg.Payload))
	}

	return nil
}
