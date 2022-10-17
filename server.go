package main

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/bus"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/util"
	rs "github.com/dipper-iot/dipper-engine/redis"
	"github.com/dipper-iot/dipper-engine/rules/arithmetic"
	"github.com/dipper-iot/dipper-engine/rules/fork"
	log2 "github.com/dipper-iot/dipper-engine/rules/log"
	"github.com/dipper-iot/dipper-engine/rules/relational"
	_switch "github.com/dipper-iot/dipper-engine/rules/switch"
	"github.com/dipper-iot/dipper-engine/store"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"time"
)

func main() {
	signalStop := make(chan os.Signal)
	signal.Notify(signalStop, os.Interrupt, os.Kill)

	var (
		dipper *core.DipperEngine
	)

	app := &cli.App{
		Name:  "Dipper Engine",
		Usage: "Rule Engine",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Config file",
				Value:   "config.json",
				Aliases: []string{"c"},
			},
			&cli.BoolFlag{
				Name:    "session-from-queue",
				Usage:   "Session from queue",
				Value:   false,
				Aliases: []string{"sq"},
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
		},
		Action: func(c *cli.Context) error {
			fmt.Println("Start Dipper Engine")

			var (
				usingRedis       bool = false
				client           *redis.Client
				storeSession     store.Store
				factoryQueue     core.FactoryQueue[*data.InputEngine]
				factoryQueueName core.FactoryQueueName[*data.OutputEngine]
				busData          bus.Bus
				config           core.ConfigEngine
				timeout          time.Duration
			)

			configFile := c.String("config")
			sessionFromQueue := c.Bool("session-from-queue")
			busType := c.String("bus")
			queueType := c.String("queue")
			storeType := c.String("store")
			redisHost := c.String("redis-host")
			redisPass := c.String("redis-pass")
			timeoutData := c.Int("redis-pass")
			timeout = time.Duration(timeoutData)

			if queueType == "redis" || busType == "redis" || storeType == "redis" {
				usingRedis = true

				client = redis.NewClient(&redis.Options{
					Addr:     redisHost,
					Password: redisPass,
				})

				err := client.Ping(c.Context).Err()
				if err != nil {
					log.Println(err)
					return err
				}
			}

			switch busType {
			case "redis":
				busData = rs.NewRedisBus(client)
				break
			default:
				busData = bus.NewDefaultBus()
				break
			}

			switch queueType {
			case "redis":
				factoryQueue = rs.FactoryQueueRedis[*data.InputEngine](client, &data.InputEngine{})
				factoryQueueName = rs.FactoryQueueNameRedis[*data.OutputEngine](client, &data.OutputEngine{})
				break
			default:
				factoryQueue = core.FactoryQueueDefault[*data.InputEngine]()
				factoryQueueName = core.FactoryQueueNameDefault[*data.OutputEngine]()
				break
			}

			switch storeType {
			case "redis":
				storeSession = rs.NewRedisStore(client, timeout)
				break
			default:
				storeSession = store.NewDefaultStore()
				break
			}

			err := util.ReadFile(&config, configFile)
			if err != nil {
				log.Println(err)
				return err
			}

			dipper = core.NewDipperEngine(
				&config,
				factoryQueue,
				factoryQueueName,
				storeSession,
				busData,
			)

			dipper.AddRule(
				&log2.LogRule{},
				&arithmetic.Arithmetic{},
				&fork.ForkRule{},
				&relational.RelationalRule{},
				&_switch.SwitchRule{},
			)

			err = dipper.Start()
			if err != nil {
				log.Println(err)
				return err
			}

			if sessionFromQueue && usingRedis {
				dipper.SessionFromQueue(rs.FactoryQueueNameRedis[*data.Session](client, &data.Session{}))
			}

			<-signalStop
			fmt.Println("Stopping Dipper Engine")
			return nil
		},
		After: func(context *cli.Context) error {
			if dipper != nil {
				return dipper.Stop()
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
