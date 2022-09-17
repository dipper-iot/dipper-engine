package core

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
)

func (d DipperEngine) startSession(ctx context.Context, sessionId uint64) error {
	if d.store.Has(sessionId) {
		sessionInfo := d.store.Get(sessionId)
		if sessionInfo.RootNode != nil {
			node := sessionInfo.RootNode
			ruleQueue, ok := d.mapQueueInputRule[node.RuleId]
			if ok {
				err := ruleQueue.Pushlish(ctx, &data.InputEngine{
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
					log.Error(err)
					return err
				}
			}
		}
	}

	return nil
}

func (d DipperEngine) registerOutput() {

	d.queueOutputRule.Subscribe(d.ctx, func(deliver *queue.Deliver[*data.OutputEngine]) {

		err := d.handlerOutput(deliver.Context, deliver.Data)
		if err != nil {
			log.Error(err)
			deliver.Reject()
			return
		}

		deliver.Ack()
	})
}

func (d DipperEngine) pushlishBus(name string, dataOutput interface{}) {
	topic, ok := d.config.BusMap[name]
	if !ok {
		topic = name
	}
	d.bus.Pushlish(context.TODO(), topic, dataOutput)
}

func (d DipperEngine) handlerOutput(ctx context.Context, dataOutput *data.OutputEngine) error {

	if dataOutput.Debug {
		d.pushlishBus("debug-output", dataOutput)
		return nil
	}

	if d.store.Has(dataOutput.SessionId) {
		return nil
	}

	if len(dataOutput.Next) == 0 {
		// finish
		session, success := d.store.Done(dataOutput.SessionId, dataOutput)
		if success {
			d.pushlishBus("session-finish", session)
		}
		return nil
	}

	if d.store.Has(dataOutput.SessionId) {
		sessionInfo := d.store.Get(dataOutput.SessionId)

		for _, nextId := range dataOutput.Next {
			node, ok := sessionInfo.MapNode[nextId]
			if ok {
				ruleQueue, ok := d.mapQueueInputRule[node.RuleId]
				if ok {
					err := ruleQueue.Pushlish(ctx, &data.InputEngine{
						SessionId:  sessionInfo.Id,
						ChanId:     sessionInfo.ChanId,
						FromEngine: node.NodeId,
						ToEngine:   dataOutput.FromEngine,
						Node:       node,
						Data:       dataOutput.Data,
						Time:       dataOutput.Time,
						Type:       dataOutput.Type,
						Error:      dataOutput.Error,
					})
					if err != nil {
						log.Error(err)
						return err
					}
				}
			}
		}
	}

	return nil
}
