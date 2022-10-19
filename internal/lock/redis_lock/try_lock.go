package redis_lock

import (
	"context"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

type TryLock struct {
	lockKey string
	client  *redis.Client
}

func NewTryLock(lockKey string, client *redis.Client) *TryLock {
	return &TryLock{lockKey: lockKey, client: client}
}

func (t *TryLock) Lock(ctx context.Context) bool {
	resp := t.client.SetNX(ctx, t.lockKey, 1, time.Second*5)
	lockSuccess, err := resp.Result()

	if err != nil {
		log.Error("lock failed ", err)
		return false
	}
	log.Trace("lock success!")
	return lockSuccess
}

func (t TryLock) Unlock(ctx context.Context) {
	delResp := t.client.Del(ctx, t.lockKey)

	unlockSuccess, err := delResp.Result()
	if err == nil && unlockSuccess > 0 {
		log.Trace("unlock success!")
	} else {
		log.Error("unlock failed ", err)
	}

}
