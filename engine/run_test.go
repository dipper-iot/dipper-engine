package engine

import (
	"context"
	"encoding/json"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/debug"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"sync"
	"testing"
)

func TestApp_Run(t *testing.T) {
	e := New()
	e.flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "test.v",
		},
		&cli.BoolFlag{
			Name: "test.paniconexit0",
		},
		&cli.BoolFlag{
			Name: "test.run",
		},
		&cli.IntFlag{
			Name: "test.timeout",
		},
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Config file",
			Value:   "../config.json",
			Aliases: []string{"c"},
		},
		&cli.BoolFlag{
			Name:    "plugin",
			Usage:   "Load Plugin",
			Value:   false,
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "session-input-queue",
			Usage:   "Session input queue",
			Value:   true,
			Aliases: []string{"iq"},
		},
		&cli.BoolFlag{
			Name:    "session-output-queue",
			Usage:   "Session output queue",
			Value:   true,
			Aliases: []string{"oq"},
		},
		&cli.StringFlag{
			Name:    "bus",
			Usage:   "Bus type: [chan, redis]",
			Value:   "redis",
			Aliases: []string{"b"},
		},
		&cli.StringFlag{
			Name:    "queue",
			Usage:   "Queue type: [chan, redis]",
			Value:   "redis",
			Aliases: []string{"q"},
		},
		&cli.StringFlag{
			Name:    "store",
			Usage:   "Store type: [memory, redis]",
			Value:   "memory",
			Aliases: []string{"s"},
		},
		&cli.StringFlag{
			Name:  "redis-host",
			Usage: "Redis host",
			Value: "127.0.0.1:6379",
		},
		&cli.StringFlag{
			Name:  "redis-pass",
			Usage: "Redis pass",
		},
		&cli.StringFlag{
			Name:  "redis-timeout",
			Usage: "Redis Time out",
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	e.Hook(AfterStart, func(dipper *core.DipperEngine, c *cli.Context) error {
		session := &data.Session{
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
						"next_success": "2",
					},
					NodeId: "4",
					RuleId: "arithmetic",
					End:    false,
				},
				"2": {
					Debug: false,
					Option: map[string]interface{}{
						"next_success": []string{"3", "4"},
					},
					NodeId: "2",
					RuleId: "fork",
					End:    false,
				},
				"3": {
					Debug:  false,
					Option: map[string]interface{}{},
					NodeId: "3",
					RuleId: "log-core",
					End:    true,
				},
				"4": {
					Debug: false,
					Option: map[string]interface{}{
						"operator": map[string]interface{}{
							"right": map[string]interface{}{
								"value": "default.a",
								"type":  "val",
							},
							"left": map[string]interface{}{
								"type":  "val",
								"value": "default.b",
							},
							"operator": "<>",
							"type":     "operator",
						},
						"set_param_result_to": "default.cond_a_b",
						"next_error":          "2",
						"next_true":           "",
						"next_false":          "",
					},
					NodeId: "4",
					RuleId: "conditional",
					End:    true,
				},
			},
		}

		dataBye, err := json.MarshalIndent(session, " ", "  ")
		if err != nil {
			log.Error(err)
			return err
		}

		if e.clientRedis == nil {
			return nil
		}
		err = e.clientRedis.RPush(context.Background(), "dipper-queue-session-input", dataBye).Err()
		if err != nil {
			log.Error(err)
			return err
		}

		go func() {
			for {
				datab, err := e.clientRedis.RPop(context.Background(), "dipper-queue-session-output").Bytes()
				if err == io.EOF {
					return
				}
				if err == redis.Nil {
					continue
				}
				if err != nil {
					log.Error(err)
					return
				}

				var result data.ResultSession
				err = json.Unmarshal(datab, &result)
				if err != nil {
					log.Error(err)
					continue
				}
				debug.PrintJson(result, "Result To Queue Output: ")
				wg.Done()
			}
		}()

		return nil
	})

	go func() {
		e.Run(os.Args)
		wg.Done()
	}()

	wg.Wait()
}

func Test_Run_Default(t *testing.T) {
	e := New()
	e.flags = []cli.Flag{
		&cli.BoolFlag{
			Name: "test.v",
		},
		&cli.BoolFlag{
			Name: "test.paniconexit0",
		},
		&cli.BoolFlag{
			Name: "test.run",
		}, &cli.IntFlag{
			Name: "test.timeout",
		},
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Config file",
			Value:   "../config.test.json",
			Aliases: []string{"c"},
		},
		&cli.BoolFlag{
			Name:    "plugin",
			Usage:   "Load Plugin",
			Value:   false,
			Aliases: []string{"p"},
		},
		&cli.BoolFlag{
			Name:    "session-input-queue",
			Usage:   "Session input queue",
			Value:   false,
			Aliases: []string{"iq"},
		},
		&cli.BoolFlag{
			Name:    "session-output-queue",
			Usage:   "Session output queue",
			Value:   false,
			Aliases: []string{"oq"},
		},
		&cli.StringFlag{
			Name:    "bus",
			Usage:   "Bus type: [chan, redis]",
			Value:   "chan",
			Aliases: []string{"b"},
		},
		&cli.StringFlag{
			Name:    "queue",
			Usage:   "Queue type: [chan, redis]",
			Value:   "chan",
			Aliases: []string{"q"},
		},
		&cli.StringFlag{
			Name:    "store",
			Usage:   "Store type: [memory, redis]",
			Value:   "memory",
			Aliases: []string{"s"},
		},
		&cli.StringFlag{
			Name:  "redis-host",
			Usage: "Redis host",
			Value: "127.0.0.1:6379",
		},
		&cli.StringFlag{
			Name:  "redis-pass",
			Usage: "Redis pass",
		},
		&cli.StringFlag{
			Name:  "redis-timeout",
			Usage: "Redis Time out",
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	e.Hook(AfterStart, func(dipper *core.DipperEngine, c *cli.Context) error {
		go func() {
			e.signalStop <- os.Interrupt
		}()
		wg.Done()
		return nil
	})

	go func() {
		e.Run(os.Args)
	}()

	wg.Wait()
}
