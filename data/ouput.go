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
	FromEngine string                 `json:"from_engine"`
	Data       map[string]interface{} `json:"data"`
	BranchMain string                 `json:"branch_main"`
	Next       []string               `json:"next"`
	Time       *time.Time             `json:"time"`
	Type       TypeOutputEngine       `json:"type"`
	Error      *errors.ErrorEngine    `json:"error"`
	Debug      bool
}
