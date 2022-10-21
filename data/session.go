package data

import (
	"github.com/dipper-iot/dipper-engine/internal/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type NodeRule struct {
	NodeId string                 `json:"node_id"`
	RuleId string                 `json:"rule_id"`
	Option map[string]interface{} `json:"option"`
	Debug  bool                   `json:"debug"`
	End    bool                   `json:"end"`
}

type Session struct {
	ChanId   string                   `json:"chan_id"`
	MapNode  map[string]*NodeRule     `json:"map_node"`
	RootNode string                   `json:"root_node"`
	Data     map[string]interface{}   `json:"data"`
	Result   map[string]*OutputEngine `json:"result"`
}

type ResultSession struct {
	Id     uint64                   `json:"id"`
	ChanId string                   `json:"chan_id"`
	Data   map[string]interface{}   `json:"data"`
	Result map[string]*OutputEngine `json:"result"`
}

type Info struct {
	Id       uint64                   `json:"id"`
	Time     *time.Time               `json:"time"`
	ChanId   string                   `json:"chan_id"`
	Timeout  time.Duration            `json:"timeout"`
	MapNode  map[string]*NodeRule     `json:"map_node"`
	RootNode *NodeRule                `json:"root_node"`
	Data     map[string]interface{}   `json:"data"`
	Result   map[string]*OutputEngine `json:"result"`
	EndCount int                      `json:"end_count"`
}

func NewSessionInfo(timeout time.Duration, data *Session) *Info {
	now := time.Now()
	var (
		id  uint64
		err error
	)
	for {
		id, err = util.NextID()
		if err != nil {
			log.Error(err)
			continue
		}
		break
	}

	return &Info{
		Id:       id,
		Time:     &now,
		ChanId:   data.ChanId,
		Timeout:  timeout,
		MapNode:  data.MapNode,
		RootNode: data.MapNode[data.RootNode],
		Data:     data.Data,
	}
}
