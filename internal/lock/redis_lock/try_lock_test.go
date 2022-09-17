package redis_lock

import (
	"context"
	"gitlab.com/dipper-iot/shared/load/rs"
	"gitlab.com/dipper-iot/shared/logger"
	"sync"
	"testing"
)

func init() {

	//logger.Init(logger.WithLevel(logger.TraceLevel))
}

func TestTryLock1(t *testing.T) {
	client, err := rs.NewRedisToEnv()
	if err != nil {
		logger.Error(err)
	}
	locker := NewTryLock("count-test", client)
	wg := sync.WaitGroup{}
	countResult := 0
	count := 0

	for i := 0; i < 100; i++ {
		wg.Add(1)
		countResult++
		go func() {
			for {
				success := locker.Lock(context.TODO())
				if success {
					count++
					locker.Unlock(context.TODO())
					wg.Done()
					return
				}
			}
		}()
	}

	wg.Wait()

	if count != countResult {
		t.Errorf("Not Match count is %d with actually %d", countResult, count)
	}
}
