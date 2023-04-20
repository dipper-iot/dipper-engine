package conditional

import (
	"context"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/queue"
	"testing"
)

func TestConditionalRule_Run(t *testing.T) {
	a := NewConditionalRule()
	a.Infinity()
	a.Initialize(context.TODO(), map[string]interface{}{})
	qsub := queue.NewDefaultQueue[*data.InputEngine]("qsub")
	qpub := queue.NewDefaultQueue[*data.OutputEngine]("qpub")

	qsub.Publish(context.TODO(), &data.InputEngine{
		SessionId:  1,
		ChanId:     "1",
		IdNode:     "noed1",
		FromEngine: "test",
		ToEngine:   "test",
		Node: &data.NodeRule{
			Option: map[string]interface{}{
				"conditional":         "a > b",
				"set_param_result_to": "cond_a_b",
				"next_error":          "4",
				"next_true":           "4",
				"next_false":          "4",
			},
			End:    false,
			Debug:  false,
			RuleId: "1",
			NodeId: "1",
		},
		Data: map[string]interface{}{
			"default": map[string]interface{}{
				"a": 2,
				"b": 2,
				"x": 3,
			},
		},
		BranchMain: "default",
	})
	qsub.Publish(context.TODO(), &data.InputEngine{
		SessionId:  1,
		ChanId:     "1",
		IdNode:     "noed1",
		FromEngine: "test",
		ToEngine:   "test",
		Node: &data.NodeRule{
			Option: map[string]interface{}{
				"conditional":         "a < b",
				"set_param_result_to": "cond_a_b",
				"next_error":          "4",
				"next_true":           "4",
				"next_false":          "4",
			},
			End:    false,
			Debug:  false,
			RuleId: "1",
			NodeId: "1",
		},
		Data: map[string]interface{}{
			"default": map[string]interface{}{
				"a": 2,
				"b": 2,
				"x": 3,
			},
		},
		BranchMain: "default",
	})
	a.Run(context.TODO(), qsub.Subscribe, qpub.Publish)
}
