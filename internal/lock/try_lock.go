package lock

import "context"

type TryLock interface {
	Lock(ctx context.Context) bool
	Unlock(ctx context.Context)
}
