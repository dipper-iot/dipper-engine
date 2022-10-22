package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	err := client.Ping(context.TODO()).Err()
	if err != nil {
		log.Println(err)
		return
	}

	session := &data.Session{
		Data: map[string]interface{}{
			"default": map[string]interface{}{
				"a": 10,
				"b": 20,
				"d": 5,
			},
		},
		ChanId:   "test-1",
		RootNode: "1",
		MapNode: map[string]*data.NodeRule{
			"1": {
				Debug: false,
				Option: map[string]interface{}{
					"list": map[string]interface{}{
						"default.c": map[string]interface{}{
							"right": map[string]interface{}{
								"value": "default.a",
								"type":  "val",
							},
							"left": map[string]interface{}{
								"type":  "val",
								"value": "default.b",
							},
							"operator": "add",
							"type":     "operator",
						},
					},
					"next_error":   "2",
					"next_success": "2",
				},
				NodeId: "4",
				RuleId: "arithmetic",
				End:    false,
			},
			"2": {
				Debug: false,
				Option: map[string]interface{}{
					"next_success": []string{"3", "4"},
				},
				NodeId: "2",
				RuleId: "fork",
				End:    false,
			},
			"3": {
				Debug:  true,
				Option: map[string]interface{}{},
				NodeId: "3",
				RuleId: "log-core",
				End:    true,
			},
			"4": {
				Debug: true,
				Option: map[string]interface{}{
					"operator": map[string]interface{}{
						"right": map[string]interface{}{
							"value": "default.a",
							"type":  "val",
						},
						"left": map[string]interface{}{
							"type":  "val",
							"value": "default.b",
						},
						"operator": "<>",
						"type":     "operator",
					},
					"set_param_result_to": "default.cond_a_b",
					"next_error":          "2",
					"next_true":           "",
					"next_false":          "",
				},
				NodeId: "4",
				RuleId: "conditional",
				End:    true,
			},
		},
	}

	dataBye, err := json.MarshalIndent(session, " ", "  ")
	if err != nil {
		log.Error(err)
		return
	}

	err = client.RPush(context.Background(), "dipper-queue-session-input", dataBye).Err()
	if err != nil {
		log.Error(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			datab, err := client.RPop(context.Background(), "dipper-queue-session-output").Bytes()
			if err == io.EOF {
				return
			}
			if err == redis.Nil {
				continue
			}
			if err != nil {
				log.Error(err)
				return
			}

			var result data.ResultSession
			err = json.Unmarshal(datab, &result)
			if err != nil {
				log.Error(err)
				continue
			}

			dataResult, err := json.MarshalIndent(result, " ", "  ")
			if err != nil {
				log.Error(err)
				return
			}

			fmt.Println("Result To Queue Output: ")
			fmt.Println(string(dataResult))
			wg.Done()
		}
	}()
	wg.Wait()
}
