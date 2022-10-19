package data

import (
	"github.com/dipper-iot/dipper-engine/internal/util"
	log "github.com/sirupsen/logrus"
	"time"
)

type NodeRule struct {
	NodeId string
	RuleId string
	Option map[string]interface{}
	Debug  bool
	End    bool
}

type Session struct {
	ChanId   string
	MapNode  map[string]*NodeRule
	RootNode string
	Data     map[string]interface{}
	Result   map[string]*OutputEngine
}

type ResultSession struct {
	Id     uint64
	ChanId string
	Data   map[string]interface{}
	Result map[string]*OutputEngine
}

type Info struct {
	Id       uint64
	Time     *time.Time
	ChanId   string
	Timeout  time.Duration
	MapNode  map[string]*NodeRule
	RootNode *NodeRule
	Data     map[string]interface{}
	Result   map[string]*OutputEngine
	EndCount int
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
