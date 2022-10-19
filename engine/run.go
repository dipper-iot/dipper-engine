package engine

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/rules/arithmetic"
	"github.com/dipper-iot/dipper-engine/rules/conditional"
	"github.com/dipper-iot/dipper-engine/rules/fork"
	log2 "github.com/dipper-iot/dipper-engine/rules/log"
	_switch "github.com/dipper-iot/dipper-engine/rules/switch"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func (a *App) Run(args []string) error {
	a.app = &cli.App{
		Name:  "Dipper Engine",
		Usage: "Rule Engine",
		Flags: a.flags,
		Action: func(c *cli.Context) error {
			fmt.Println("Start Dipper Engine")
			fmt.Println()

			if a.dipper == nil {
				err := a.newEngine(c)
				if err != nil {
					log.Println(err)
					return err
				}
			}

			err := a.runHooks(BeforeStart, c)
			if err != nil {
				log.Println(err)
				return err
			}

			a.dipper.AddRule(
				&log2.LogRule{},
				&arithmetic.Arithmetic{},
				&fork.ForkRule{},
				&conditional.ConditionalRule{},
				&_switch.SwitchRule{},
			)

			err = a.dipper.Start()
			if err != nil {
				log.Println(err)
				return err
			}

			err = a.runHooks(AfterStart, c)
			if err != nil {
				log.Println(err)
				return err
			}

			<-a.signalStop

			err = a.runHooks(BeforeStop, c)
			if err != nil {
				log.Println(err)
				return err
			}

			fmt.Println("Stopping Dipper Engine")
			return nil
		},
		After: func(context *cli.Context) error {
			if a.dipper != nil {
				return a.dipper.Stop()
			}
			return nil
		},
	}

	if err := a.app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
