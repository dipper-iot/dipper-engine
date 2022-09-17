package queue

import (
	"context"
	"time"
)

type defaultQueue[T any] struct {
	name  string
	queue chan T
}

func (d defaultQueue[T]) Name() string {
	return d.name
}

func (d defaultQueue[T]) Pushlish(ctx context.Context, input T) error {
	d.queue <- input
	return nil
}

func (d defaultQueue[T]) Subscribe(ctx context.Context, callback SubscribeFunction[T]) error {
	go d.loopData(ctx, callback)
	return nil
}

func (d defaultQueue[T]) loopData(ctx context.Context, callback SubscribeFunction[T]) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-d.queue:
			timeData := time.Now()
			callback(NewDeliver[T](
				ctx,
				data,
				&timeData,
				nil,
				nil,
			))
		}
	}
}

func NewDefaultQueue[T any](name string) *defaultQueue[T] {
	return &defaultQueue[T]{
		name:  name,
		queue: make(chan T),
	}
}
