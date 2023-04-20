package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/engine"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})

	err := client.Ping(context.TODO()).Err()
	if err != nil {
		log.Println(err)
		return
	}

	c := engine.NewWithConfig(getConfig())

	c.Hook(engine.AfterStart, func(dipper *core.DipperEngine, c *cli.Context) error {

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
						"operators": map[string]string{
							"c": "a+b",
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
						"next_success": []string{"3", "4", "5"},
					},
					NodeId: "2",
					RuleId: "fork",
					End:    false,
				},
				"5": {
					Debug:  true,
					Option: map[string]interface{}{},
					NodeId: "5",
					RuleId: "output-redis-queue",
					End:    true,
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
						"conditional":         "a!=b",
						"set_param_result_to": "cond_a_b",
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

		dipper.Add(context.Background(), session)

		return nil
	})

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
		}
	}()

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
