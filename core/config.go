package core

type RuleOption struct {
	Enable  bool                   `json:"enable"`
	Worker  int                    `json:"worker"`
	Options map[string]interface{} `json:"options"`
}

type ConfigEngine struct {
	BusMap         map[string]string      `json:"bus_map"`
	Rules          map[string]*RuleOption `json:"rules"`
	TimeoutSession int                    `json:"timeout_session"`
	Log            ConfigLog              `json:"log"`
}

type ConfigLog struct {
	Level         string `json:"level"`
	Output        string `json:"output"`
	FileName      string `json:"file_name"`
	LogMethodName bool   `json:"log_method_name"`
}
