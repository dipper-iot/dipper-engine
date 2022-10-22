package engine

import (
	"github.com/dipper-iot/dipper-engine/bus"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/dipper-iot/dipper-engine/data"
	"github.com/dipper-iot/dipper-engine/internal/util"
	rs "github.com/dipper-iot/dipper-engine/redis"
	"github.com/dipper-iot/dipper-engine/store"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"time"
)

func (a *App) newEngine(c *cli.Context) error {
	var (
		usingRedis       = false
		client           *redis.Client
		storeSession     store.Store
		factoryQueue     core.FactoryQueue[*data.InputEngine]
		factoryQueueName core.FactoryQueueName[*data.OutputEngine]
		busData          bus.Bus
		config           core.ConfigEngine
		timeout          time.Duration
	)

	configFile := c.String("config")
	pluginEnable := c.Bool("plugin")
	sessionInputmQueue := c.Bool("session-input-queue")
	sessionOutputQueue := c.Bool("session-output-queue")
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

		a.clientRedis = client
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

	a.dipper = core.NewDipperEngine(
		&config,
		factoryQueue,
		factoryQueueName,
		storeSession,
		busData,
	)

	a.beforeStartHooks = append(a.beforeStartHooks, func(dipper *core.DipperEngine, c *cli.Context) error {
		if sessionOutputQueue && usingRedis {
			a.dipper.SessionOutputQueue(rs.FactoryQueueNameRedis[*data.ResultSession](client, &data.ResultSession{}))
		}
		if sessionInputmQueue && usingRedis {
			a.dipper.SessionInputQueue(rs.FactoryQueueNameRedis[*data.Session](client, &data.Session{}))
		}
		return nil
	})

	if pluginEnable {
		a.dipper.LoadPlugin()
	}

	return nil
}
