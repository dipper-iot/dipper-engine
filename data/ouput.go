package data

import (
	"github.com/dipper-iot/dipper-engine/errors"
	"time"
)

type TypeOutputEngine string

const (
	TypeOutputEngineError   TypeOutputEngine = "error"
	TypeOutputEngineSuccess TypeOutputEngine = "success"
)

type OutputEngine struct {
	SessionId  uint64                 `json:"session_id"`
	ChanId     string                 `json:"chan_id"`
	IdNode     string                 `json:"id_node"`
	FromEngine string                 `json:"from_engine"`
	MetaData   map[string]interface{} `json:"meta_data"`
	Data       map[string]interface{} `json:"data"`
	BranchMain string                 `json:"branch_main"`
	Next       []string               `json:"next"`
	Time       *time.Time             `json:"time"`
	Type       TypeOutputEngine       `json:"type"`
	Error      *errors.ErrorEngine    `json:"error"`
	Debug      bool
}

func (o OutputEngine) Clone() *OutputEngine {
	return &OutputEngine{
		SessionId:  o.SessionId,
		ChanId:     o.ChanId,
		IdNode:     o.IdNode,
		FromEngine: o.FromEngine,
		Data:       o.Data,
		BranchMain: o.BranchMain,
		Next:       o.Next,
		Time:       o.Time,
		Type:       o.Type,
		Error:      o.Error,
		Debug:      o.Debug,
	}
}
