package main

import (
	"context"
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/bus"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/dipper-iot/dipper-engine/rules/arithmetic"
	"github.com/dipper-iot/dipper-engine/rules/fork"
	log2 "github.com/dipper-iot/dipper-engine/rules/log"
	"github.com/dipper-iot/dipper-engine/rules/relational"
	_switch "github.com/dipper-iot/dipper-engine/rules/switch"
	"github.com/dipper-iot/dipper-engine/store"
	"log"
	"sync"
	"testing"
)

func Test_Run(t *testing.T) {
	var (
		storeSession     store.Store
		factoryQueue     core.FactoryQueue[*data.InputEngine]
		factoryQueueName core.FactoryQueueName[*data.OutputEngine]
		busData          bus.Bus
		config           core.ConfigEngine
	)
	err := util.ReadFile(&config, "config.json")
	if err != nil {
		log.Println(err)
		return
	}

	busData = bus.NewDefaultBus()
	factoryQueue = core.FactoryQueueDefault[*data.InputEngine]()
	factoryQueueName = core.FactoryQueueNameDefault[*data.OutputEngine]()
	storeSession = store.NewDefaultStore()

	dipper := core.NewDipperEngine(
		&config,
		factoryQueue,
		factoryQueueName,
		storeSession,
		busData,
	)

	wg := sync.WaitGroup{}
	wg.Add(1)

	dipper.AddRule(
		&log2.LogRule{},
		&arithmetic.Arithmetic{},
		&fork.ForkRule{},
		&relational.RelationalRule{},
		&_switch.SwitchRule{},
		&LogTest{
			&wg,
		},
	)

	err = dipper.Start()
	if err != nil {
		log.Println(err)
	}

	err = dipper.Add(context.Background(), &data.Session{
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
				Debug: true,
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
					"next_success": "3",
					"debug":        false,
				},
				NodeId: "1",
				RuleId: "arithmetic",
				End:    false,
			},
			"2": {
				Debug:  true,
				Option: map[string]interface{}{},
				NodeId: "2",
				RuleId: "logger",
				End:    true,
			},
			"3": {
				Debug: true,
				Option: map[string]interface{}{
					"next_success": []string{"4", "2"},
					"debug":        true,
				},
				NodeId: "3",
				RuleId: "fork",
				End:    false,
			},
			"4": {
				Debug:  true,
				Option: map[string]interface{}{},
				NodeId: "4",
				RuleId: "test",
				End:    true,
			},
		},
	})
	if err != nil {
		log.Println(err)
	}

	wg.Wait()
}

type LogTest struct {
	wg *sync.WaitGroup
}

func (l *LogTest) Id() string {
	return "test"
}

func (l *LogTest) Initialize(ctx context.Context, option map[string]interface{}) error {

	return nil
}

func (l *LogTest) Run(ctx context.Context, subscribeQueueInput func(ctx context.Context, callback queue.SubscribeFunction[*data.InputEngine]) error, pushQueueOutput func(ctx context.Context, input *data.OutputEngine) error) {
	err := subscribeQueueInput(ctx, func(deliver *queue.Deliver[*data.InputEngine]) {
		defer deliver.Ack()
		dataByte, err := json.Marshal(deliver.Data.Data)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(dataByte))
		l.wg.Done()
	})
	if err != nil {
		log.Println(err)
		return
	}

}

func (l *LogTest) Stop(ctx context.Context) error {

	return nil
}
