package core

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

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

	list := make([]*viewRule, 0)
	index := 1
	// Run Rule
	for name, rule := range d.mapRule {
		infinity := rule.Infinity()
		infinityStr := "false"
		if infinity {
			infinityStr = "true"
			control, ok := rule.(SessionControl)
			if !ok {
				log.Error("Error: ", errors.ErrorNotControlEngine, rule.Id())
				return errors.ErrorNotControlEngine
			}
			d.mapSessionControl[rule.Id()] = control
		}

		queueInput, ok := d.mapQueueInputRule[name]
		if !ok {
			return errors.ErrorNotFoundQueue
		}
		status := ""
		worker := 0
		option, ok := d.config.Rules[name]
		if ok && option.Enable {
			for i := 0; i < option.Worker; i++ {
				go rule.Run(d.ctx, queueInput.Subscribe, d.queueOutputRule.Publish)
			}
			status = "enable"
			worker = option.Worker
		} else {
			status = "disable"
		}
		list = append(list, &viewRule{
			Name:     name,
			Infinity: infinityStr,
			Worker:   worker,
			Status:   status,
		})

		index++
	}
	viewListRule(list)
	go d.registerOutput()
	fmt.Println("Running Engine...")

	return nil
}

type viewRule struct {
	Name     string
	Worker   int
	Infinity string
	Status   string
}

func viewListRule(list []*viewRule) {
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i].Name, list[j].Name) < 0
	})
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Println(fmt.Sprintf("Rules: %d", len(list)))
	fmt.Println("-----------------------------------------------------------")
	fmt.Fprintln(w, "No\t\tRule Name\t\tWorker\t\tInfinity\t\tStatus\t")
	index := 1
	// Run Rule
	for _, rule := range list {
		fmt.Fprintln(w, fmt.Sprintf("%d\t\t%s\t\t%d\t\t%s\t\t%s\t", index, rule.Name, rule.Worker, rule.Infinity, rule.Status))
		index++
	}
	w.Flush()
	fmt.Println("-----------------------------------------------------------")
	fmt.Println()
}
