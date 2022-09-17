package bus

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/dipper-iot/dipper-engine/internal/util"
)

type defaultBus struct {
	bus EventBus.Bus
}

func (d defaultBus) Subscribe(ctx context.Context, name string, callback func(data []byte)) error {
	d.bus.Subscribe(name, callback)
	return nil
}

func NewDefaultBus() *defaultBus {
	return &defaultBus{
		bus: EventBus.New(),
	}
}

func (d defaultBus) Pushlish(ctx context.Context, name string, input interface{}) error {
	data, err := util.ConvertToByte(input)
	if err != nil {
		return err
	}
	d.bus.Publish(name, data)
	return nil
}
