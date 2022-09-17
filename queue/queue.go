package queue

import (
	"context"
	"time"
)

type SubscribeFunction[T any] func(d *Deliver[T])

type Deliver[T any] struct {
	ackCallBack    func()
	rejectCallBack func()
	Context        context.Context
	Data           T
	Time           *time.Time
}

func NewDeliver[T any](
	context context.Context,
	data T,
	time *time.Time,
	ackCallBack func(),
	rejectCallBack func(),
) *Deliver[T] {
	return &Deliver[T]{
		ackCallBack:    ackCallBack,
		rejectCallBack: rejectCallBack,
		Context:        context,
		Data:           data,
		Time:           time,
	}
}

func (d Deliver[T]) Ack() {
	if d.ackCallBack == nil {
		return
	}
	d.ackCallBack()
}

func (d Deliver[T]) Reject() {
	if d.rejectCallBack == nil {
		return
	}
	d.rejectCallBack()
}

type QueueEngine[T any] interface {
	Name() string
	Pushlish(ctx context.Context, input T) error
	Subscribe(ctx context.Context, callback SubscribeFunction[T]) error
}
