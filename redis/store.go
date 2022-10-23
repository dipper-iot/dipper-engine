package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/lock"
	"github.com/dipper-iot/dipper-engine/internal/lock/redis_lock"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

type redisStore struct {
	client  *redis.Client
	locker  lock.TryLock
	timeout time.Duration
}

func NewRedisStore(client *redis.Client, timeout time.Duration) *redisStore {
	return &redisStore{
		client:  client,
		timeout: timeout,
		locker:  redis_lock.NewTryLock("dipper-store-session-locker", client),
	}
}

func (r redisStore) getKey(sessionId uint64) string {
	return fmt.Sprintf("dipper-store-session-%d", sessionId)
}

func (r redisStore) Add(sessionInfo *data.Info) error {
	key := r.getKey(sessionInfo.Id)

	data, err := util.ConvertToByte(sessionInfo)
	if err != nil {
		log.Error(err)
		return err
	}

	return r.client.Set(context.TODO(), key, data, r.timeout).Err()
}

func (r redisStore) Get(sessionId uint64) *data.Info {
	key := r.getKey(sessionId)

	dataStr, err := r.client.Get(context.TODO(), key).Result()
	if err != nil {
		log.Error(err)
		return nil
	}

	var data data.Info
	err = json.Unmarshal([]byte(dataStr), data)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &data
}

func (r redisStore) Has(sessionId uint64) bool {
	key := r.getKey(sessionId)

	dataStr, err := r.client.Get(context.TODO(), key).Result()
	if err != nil {
		log.Error(err)
		return false
	}
	var data data.Info
	err = json.Unmarshal([]byte(dataStr), data)
	if err != nil {
		log.Error(err)
		return false
	}

	return data.Id == sessionId
}

func (r redisStore) delete(sessionId uint64) {
	key := r.getKey(sessionId)

	err := r.client.Del(context.TODO(), key).Err()
	if err != nil {
		log.Error(err)

	}
}

func (r redisStore) Done(sessionId uint64, result *data.OutputEngine) (session *data.ResultSession, success bool) {
	ctx := context.TODO()
	for {
		ok := r.locker.Lock(ctx)
		if ok {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	defer r.locker.Unlock(ctx)

	success = false
	if !r.Has(sessionId) {
		return
	}

	sessionInfo := r.Get(sessionId)

	if sessionInfo.Infinite {
		return nil, false
	}

	sessionInfo.EndCount -= 1
	if sessionInfo.Result == nil {
		sessionInfo.Result = map[string]*data.OutputEngine{}
	}
	sessionInfo.Result[result.IdNode] = result.Clone()

	if sessionInfo.EndCount > 0 {
		return
	}
	success = true
	// delete store
	r.delete(sessionId)
	// result
	session = &data.ResultSession{
		Id:     sessionInfo.Id,
		Data:   sessionInfo.Data,
		ChanId: sessionInfo.ChanId,
		Result: sessionInfo.Result,
	}

	return
}
