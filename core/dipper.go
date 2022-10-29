package core

import (
	"context"
	bus2 "github.com/dipper-iot/dipper-engine/bus"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/dipper-iot/dipper-engine/store"
	log "github.com/sirupsen/logrus"
)

type DipperEngine struct {
	ctx                context.Context
	cancel             context.CancelFunc
	config             *ConfigEngine
	mapRule            map[string]Rule
	mapSessionControl  map[string]SessionControl
	mapQueueInputRule  map[string]queue.QueueEngine[*data.InputEngine]
	queueOutputRule    queue.QueueEngine[*data.OutputEngine]
	factoryQueue       FactoryQueue[*data.InputEngine]
	factoryQueueOutput FactoryQueueName[*data.OutputEngine]
	queueInput         queue.QueueEngine[*data.Session]
	queueOutput        queue.QueueEngine[*data.ResultSession]
	store              store.Store
	bus                bus2.Bus
}

func NewDipperEngine(
	config *ConfigEngine,
	factoryQueue FactoryQueue[*data.InputEngine],
	factoryQueueOutput FactoryQueueName[*data.OutputEngine],
	store store.Store,
	bus bus2.Bus,
) *DipperEngine {
	ctx, cancel := context.WithCancel(context.TODO())
	return &DipperEngine{
		ctx:                ctx,
		cancel:             cancel,
		config:             config,
		factoryQueue:       factoryQueue,
		factoryQueueOutput: factoryQueueOutput,
		store:              store,
		bus:                bus,
		mapRule:            map[string]Rule{},
		mapQueueInputRule:  map[string]queue.QueueEngine[*data.InputEngine]{},
		mapSessionControl:  map[string]SessionControl{},
	}
}

func (d *DipperEngine) SetConfig(conf *ConfigEngine) {
	d.config = conf
}

func (d *DipperEngine) SetContext(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	d.ctx = ctx
	d.cancel = cancel
}

func (d *DipperEngine) AddRule(rules ...Rule) {
	for _, rule := range rules {
		if rule != nil {
			d.addRule(rule)
		}
	}
}

func (d *DipperEngine) addRule(rule Rule) {
	log.Tracef("Add Rule: %s", rule.Id())
	d.mapRule[rule.Id()] = rule

	queue := d.factoryQueue(rule)
	log.Tracef("Add Queue: %s", queue.Name())
	d.mapQueueInputRule[rule.Id()] = queue
}

func (d *DipperEngine) Stop() error {
	d.cancel()
	return nil
}

func (d *DipperEngine) RuleEnable() []string {
	list := make([]string, 0)
	for name, option := range d.config.Rules {
		if option.Enable {
			list = append(list, name)
		}
	}
	return list
}
