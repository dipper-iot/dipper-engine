package data

import (
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
	MetaData map[string]interface{}   `json:"meta_data"`
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
	Infinite bool                     `json:"infinite"`
	MapNode  map[string]*NodeRule     `json:"map_node"`
	RootNode *NodeRule                `json:"root_node"`
	Data     map[string]interface{}   `json:"data"`
	Result   map[string]*OutputEngine `json:"result"`
	EndCount int                      `json:"end_count"`
}
