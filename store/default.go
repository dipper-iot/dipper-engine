package store

import (
	"github.com/dipper-iot/dipper-engine/data"
	"sync"
)

type defaultStore struct {
	mapData sync.Map
}

func NewDefaultStore() Store {
	return &defaultStore{
		mapData: sync.Map{},
	}
}

func (d *defaultStore) Add(session *data.Info) error {
	d.mapData.Store(session.Id, session)
	return nil
}

func (d *defaultStore) Get(sessionId uint64) *data.Info {
	raw, ok := d.mapData.Load(sessionId)
	if ok && raw != nil {
		return raw.(*data.Info)
	}
	return nil
}

func (d *defaultStore) Has(sessionId uint64) bool {
	raw, ok := d.mapData.Load(sessionId)
	return ok && raw != nil
}

func (d *defaultStore) Done(sessionId uint64, result *data.OutputEngine) (session *data.ResultSession, success bool) {
	success = false
	if !d.Has(sessionId) {
		return
	}

	sessionInfo := d.Get(sessionId)

	if sessionInfo.Infinite {
		return nil, false
	}

	sessionInfo.EndCount -= 1
	if sessionInfo.Result == nil {
		sessionInfo.Result = map[string]*data.OutputEngine{}
	}
	sessionInfo.Result[result.IdNode] = result

	if sessionInfo.EndCount > 0 {
		return
	}
	success = true
	// delete store
	d.mapData.Delete(sessionId)
	// result
	session = &data.ResultSession{
		Id:     sessionInfo.Id,
		Data:   sessionInfo.Data,
		ChanId: sessionInfo.ChanId,
		Result: sessionInfo.Result,
	}

	return
}
