package bus

import (
	"context"
)

type Bus interface {
	Pushlish(ctx context.Context, name string, input interface{}) error
	Subscribe(ctx context.Context, name string, callback func(data []byte)) error
}
