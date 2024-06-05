package core

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewSessionInfo(timeout time.Duration, session *data.Session, mapRule map[string]Rule) *data.Info {
	now := time.Now()
	var (
		id  uint64 = session.Id
		err error
	)
	if id == 0 {
		for {
			id, err = util.NextID()
			if err != nil {
				log.Error(err)
				continue
			}
			break
		}
	}

	endCount := 0
	infinite := false
	for _, rule := range session.MapNode {
		if rule.End {
			endCount++
		}
		rulInfo, ok := mapRule[rule.RuleId]
		if ok && rulInfo.Infinity() {
			infinite = true
		}
	}

	return &data.Info{
		Id:       id,
		Time:     &now,
		Infinite: infinite,
		ChanId:   session.ChanId,
		Timeout:  timeout,
		MapNode:  session.MapNode,
		EndCount: endCount,
		RootNode: session.MapNode[session.RootNode],
		Data:     session.Data,
	}
}

func (d *DipperEngine) StartSession(ctx context.Context, sessionId uint64) error {
	if d.store.Has(sessionId) {
		sessionInfo := d.store.Get(sessionId)
		if sessionInfo.RootNode != nil {
			node := sessionInfo.RootNode
			ruleQueue, ok := d.mapQueueInputRule[node.RuleId]
			if ok {
				err := ruleQueue.Publish(ctx, &data.InputEngine{
					SessionId:  sessionInfo.Id,
					ChanId:     sessionInfo.ChanId,
					FromEngine: node.NodeId,
					ToEngine:   "",
					Node:       node,
					Data:       sessionInfo.Data,
					Time:       sessionInfo.Time,
					Type:       data.TypeOutputEngineSuccess,
					Error:      nil,
				})
				if err != nil {
					log.Error("Publish have error ", err)
					return err
				}
			}
		}
	}

	return nil
}

func (d *DipperEngine) Add(ctx context.Context, sessionData *data.Session) error {
	sessionInfo := NewSessionInfo(time.Duration(d.config.TimeoutSession), sessionData, d.mapRule)
	d.store.Add(sessionInfo)
	return d.StartSession(ctx, sessionInfo.Id)
}

func (d *DipperEngine) SessionInputQueue(factoryQueueName FactoryQueueName[*data.Session]) {
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

func (d *DipperEngine) SessionOutputQueue(factoryQueueOutputName FactoryQueueName[*data.ResultSession]) {
	defaultOutputTopic := "session-output"
	topic, ok := d.config.BusMap[defaultOutputTopic]
	if !ok {
		topic = defaultOutputTopic
	}
	d.queueOutput = factoryQueueOutputName(topic)
}

func (d *DipperEngine) OutputSubscribe(ctx context.Context, callback queue.SubscribeFunction[*data.ResultSession]) {
	d.queueOutput.Subscribe(ctx, callback)
}
