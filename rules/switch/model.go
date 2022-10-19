package _switch

type Option struct {
	Key       string            `json:"key"`
	MapSwitch map[string]string `json:"map_switch"`
	NextError string            `json:"next_error"`
	Debug     bool              `json:"debug"`
}
