package core

import "github.com/dipper-iot/dipper-engine/queue"

type FactoryQueue[T any] func(engine Rule) queue.QueueEngine[T]
type FactoryQueueName[T any] func(name string) queue.QueueEngine[T]

func FactoryQueueDefault[T any]() FactoryQueue[T] {
	return func(engine Rule) queue.QueueEngine[T] {
		return queue.NewDefaultQueue[T](engine.Id())
	}
}

func FactoryQueueNameDefault[T any]() FactoryQueueName[T] {
	return func(name string) queue.QueueEngine[T] {
		return queue.NewDefaultQueue[T](name)
	}
}
