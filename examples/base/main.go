package main

import (
	"context"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/engine"
	"github.com/dipper-iot/dipper-engine/internal/debug"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {

	c := engine.NewWithConfig(getConfig())

	c.Hook(engine.AfterStart, func(dipper *core.DipperEngine, c *cli.Context) error {

		factoryResultSessionName := core.FactoryQueueNameDefault[*data.ResultSession]()
		dipper.SessionOutputQueue(factoryResultSessionName)

		dipper.OutputSubscribe(context.TODO(), func(sessionDeliver *queue.Deliver[*data.ResultSession]) {

			debug.PrintJson(sessionDeliver.Data, "Result: ")

			sessionDeliver.Ack()
		})

		return dipper.Add(context.Background(), &data.Session{
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
						"next_success": "3",
						"debug":        false,
					},
					NodeId: "4",
					RuleId: "arithmetic",
					End:    false,
				},
				"2": {
					Debug: false,
					Option: map[string]interface{}{
						"debug": false,
					},
					NodeId: "2",
					RuleId: "log-core",
					End:    true,
				},
				"3": {
					Debug: false,
					Option: map[string]interface{}{
						"next_success": []string{"5", "2"},
						"debug":        false,
					},
					NodeId: "3",
					RuleId: "fork",
					End:    false,
				},
				"5": {
					Debug: false,
					Option: map[string]interface{}{
						"conditional":         "a == b",
						"set_param_result_to": "cond_a_b",
						"next_error":          "2",
						"next_true":           "2",
						"next_false":          "2",
						"debug":               false,
					},
					NodeId: "5",
					RuleId: "conditional",
					End:    false,
				},
			},
		})
	})

	if err := c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
