package main

import (
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/core"
	log "github.com/sirupsen/logrus"
)

var configMap = map[string]interface{}{
	"rules": map[string]interface{}{
		"log-core": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
		"arithmetic": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
		"fork": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
		"switch": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
		"conditional": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
		"input-redis-queue": map[string]interface{}{
			"enable": false,
			"worker": 1,
			"options": map[string]interface{}{
				"redis_address": "127.0.0.1:6379",
				"redis_db":      0,
			},
		},
		"input-redis-queue-extend": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
		"output-redis-queue": map[string]interface{}{
			"enable": false,
			"worker": 1,
			"options": map[string]interface{}{
				"redis_address": "127.0.0.1:6379",
				"redis_db":      0,
			},
		},
		"output-redis-queue-extend": map[string]interface{}{
			"enable": true,
			"worker": 1,
		},
	},
	"log": map[string]interface{}{
		"level":     "info",
		"out":       "console",
		"file_name": "dipper-engine.log",
	},
	"timeout_session": 30,
	"plugins":         []string{},
}

func getConfig() *core.ConfigEngine {
	var config core.ConfigEngine
	data, err := json.Marshal(configMap)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println(err)
	}
	return &config
}
