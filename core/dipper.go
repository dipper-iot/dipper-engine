package core

import (
	"context"
	bus2 "github.com/dipper-iot/dipper-engine/bus"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/errors"
	"github.com/dipper-iot/dipper-engine/queue"
	"github.com/dipper-iot/dipper-engine/store"
	log "github.com/sirupsen/logrus"
	"time"
)

type DipperEngine struct {
	ctx                context.Context
	cancel             context.CancelFunc
	config             *ConfigEngine
	mapRule            map[string]Rule
	mapQueueInputRule  map[string]queue.QueueEngine[*data.InputEngine]
	queueOutputRule    queue.QueueEngine[*data.OutputEngine]
	factoryQueue       FactoryQueue[*data.InputEngine]
	factoryQueueOutput FactoryQueueName[*data.OutputEngine]
	queueInput         queue.QueueEngine[*data.Session]
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
	}
}

func (d *DipperEngine) WithContext(ctx context.Context) *DipperEngine {
	ctx, cancel := context.WithCancel(ctx)
	d.ctx = ctx
	d.cancel = cancel
	return d
}

func (d *DipperEngine) LoadRulePlugin() {

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

func (d *DipperEngine) Add(ctx context.Context, sessionData *data.Session) error {
	sessionInfo := data.NewSessionInfo(time.Duration(d.config.TimeoutSession), sessionData)
	d.store.Add(sessionInfo)
	return d.startSession(ctx, sessionInfo.Id)
}

func (d *DipperEngine) SessionFromQueue(factoryQueueName FactoryQueueName[*data.Session]) {
	defaultTopic := "session-input"
	topic, ok := d.config.BusMap[defaultTopic]
	if !ok {
		topic = defaultTopic
	}

	d.queueInput = factoryQueueName(topic)

	d.queueInput.Subscribe(context.TODO(), func(sessionDeliver *queue.Deliver[*data.Session]) {
		err := d.Add(context.TODO(), sessionDeliver.Data)
		if err != nil {
			sessionDeliver.Reject()
			return
		}
		sessionDeliver.Ack()
	})
}

func (d *DipperEngine) Start() error {
	log.Debug("Start Dipper Engine")
	d.queueOutputRule = d.factoryQueueOutput("output")

	// init Rule
	for name, rule := range d.mapRule {
		option, ok := d.config.Rules[name]
		if ok && option.Enable {
			err := rule.Initialize(d.ctx, map[string]interface{}{})
			if err != nil {
				return err
			}
		}

	}

	// Run Rule
	for name, rule := range d.mapRule {
		queueInput, ok := d.mapQueueInputRule[name]
		if !ok {
			return errors.ErrorNotFoundQueue
		}
		option, ok := d.config.Rules[name]
		if ok && option.Enable {
			for i := 0; i < option.Worker; i++ {
				go rule.Run(d.ctx, queueInput.Subscribe, d.queueOutputRule.Publish)
			}
		}
	}

	go d.registerOutput()

	return nil
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
