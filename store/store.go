package store

import "github.com/dipper-iot/dipper-engine/data"

type Store interface {
	Add(session *data.Info) error
	Get(sessionId uint64) *data.Info
	Has(sessionId uint64) bool
	Done(sessionId uint64, result *data.OutputEngine) (session *data.ResultSession, success bool)
}
