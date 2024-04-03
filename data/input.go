package data

import (
	"github.com/dipper-iot/dipper-engine/errors"
	"time"
)

type InputEngine struct {
	SessionId  uint64                 `json:"session_id"`
	ChanId     string                 `json:"chan_id"`
	IdNode     string                 `json:"id_node"`
	FromEngine string                 `json:"from_engine"`
	ToEngine   string                 `json:"to_engine"`
	Node       *NodeRule              `json:"node"`
	MetaData   map[string]interface{} `json:"meta_data"`
	Data       map[string]interface{} `json:"data"`
	BranchMain string                 `json:"branch_main"`
	Type       TypeOutputEngine       `json:"type"`
	Error      *errors.ErrorEngine    `json:"error"`
	Time       *time.Time             `json:"time"`
}
