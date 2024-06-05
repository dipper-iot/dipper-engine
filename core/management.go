package core

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/debug"
	"github.com/dipper-iot/dipper-engine/internal/util"
	"github.com/dipper-iot/dipper-engine/queue"
	log "github.com/sirupsen/logrus"
)

func (d *DipperEngine) registerOutput() {

	err := d.queueOutputRule.Subscribe(d.ctx, func(deliver *queue.Deliver[*data.OutputEngine]) {

		err := d.handlerOutput(deliver.Context, deliver.Data)
		if err != nil {
			log.Error(err)
			deliver.Reject()
			return
		}

		deliver.Ack()
	})
	if err != nil {
		log.Error(err)
	}
}

func (d *DipperEngine) publishBus(name string, dataOutput interface{}) {
	topic, ok := d.config.BusMap[name]
	if !ok {
		topic = name
	}
	d.bus.Pushlish(context.TODO(), topic, dataOutput)
}

func (d *DipperEngine) handlerOutput(ctx context.Context, dataOutput *data.OutputEngine) error {

	dataOutput = dataOutput.Clone()
	if dataOutput.Debug {
		debug.PrintJson(dataOutput, "Debug-output => ChanId: %s | Form: %s ", dataOutput.ChanId, dataOutput.FromEngine)
		d.publishBus("debug-output", dataOutput)
		//return nil
	}

	if !d.store.Has(dataOutput.SessionId) {
		return nil
	}

	dataOutput.Next = util.ValidateNext(dataOutput.Next)

	log.Debug("Next => ", len(dataOutput.Next), "FromEngine => ", dataOutput.FromEngine)
	if len(dataOutput.Next) == 0 {
		// finish
		session, success := d.store.Done(dataOutput.SessionId, dataOutput)
		if success {
			d.publishBus("session-finish", session)
			if d.queueOutput != nil {
				d.queueOutput.Publish(ctx, session)
			}
		}

		return nil
	}

	if d.store.Has(dataOutput.SessionId) {
		sessionInfo := d.store.Get(dataOutput.SessionId)

		for _, nextId := range dataOutput.Next {
			node, ok := sessionInfo.MapNode[nextId]
			if ok {
				if node.Option == nil {
					node.Option = map[string]interface{}{}
				}
				node.Option["debug"] = node.Debug
				ruleQueue, ok := d.mapQueueInputRule[node.RuleId]
				if ok {
					err := ruleQueue.Publish(ctx, &data.InputEngine{
						SessionId:  sessionInfo.Id,
						ChanId:     sessionInfo.ChanId,
						IdNode:     nextId,
						MetaData:   dataOutput.MetaData,
						BranchMain: dataOutput.BranchMain,
						FromEngine: node.RuleId,
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
				} else {
					log.Errorf("Not found Rule Id: %s", node.NodeId)
				}
			} else {
				log.Errorf("Not found next Id Id: %s", nextId)
			}
		}
	}

	return nil
}
