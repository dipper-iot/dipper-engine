package engine

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/core"
	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
)

type App struct {
	config      *core.ConfigEngine
	flags       []cli.Flag
	app         *cli.App
	dipper      *core.DipperEngine
	signalStop  chan os.Signal
	clientRedis *redis.Client
	// hook
	beforeStartHooks []HookFunc
	beforeStopHooks  []HookFunc
	afterStartHooks  []HookFunc
}

func New(flags ...cli.Flag) *App {
	return NewWithConfig(nil, flags...)
}

func NewWithConfig(config *core.ConfigEngine, flags ...cli.Flag) *App {

	signalStop := make(chan os.Signal)
	signal.Notify(signalStop, os.Interrupt, os.Kill)

	flags = append(flags,
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Config file",
			Value:   "config.json",
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
	)
	return &App{
		config:           config,
		flags:            flags,
		signalStop:       signalStop,
		beforeStartHooks: []HookFunc{},
		afterStartHooks:  []HookFunc{},
		beforeStopHooks:  []HookFunc{},
	}
}

func NewWithEngine(engine *core.DipperEngine, flags ...cli.Flag) *App {

	signalStop := make(chan os.Signal)
	signal.Notify(signalStop, os.Interrupt, os.Kill)

	return &App{
		dipper:           engine,
		flags:            flags,
		signalStop:       signalStop,
		beforeStartHooks: []HookFunc{},
		afterStartHooks:  []HookFunc{},
		beforeStopHooks:  []HookFunc{},
	}
}

type HookType int8
type HookFunc func(dipper *core.DipperEngine, c *cli.Context) error

const (
	BeforeStart HookType = 1
	AfterStart  HookType = 2
	BeforeStop  HookType = 3
)

func (a *App) Engine() *core.DipperEngine {
	return a.dipper
}

func (a *App) Hook(typeHook HookType, callback HookFunc) {
	if callback == nil {
		return
	}
	switch typeHook {
	case AfterStart:
		a.afterStartHooks = append(a.afterStartHooks, callback)
		break
	case BeforeStart:
		a.beforeStartHooks = append(a.beforeStartHooks, callback)
		break
	case BeforeStop:
		a.beforeStopHooks = append(a.beforeStopHooks, callback)
		break
	}
}

func (a *App) Stop() {
	a.signalStop <- os.Interrupt
}

func (a *App) runHooks(typeHook HookType, c *cli.Context) error {
	var hooks []HookFunc

	switch typeHook {
	case AfterStart:
		hooks = a.afterStartHooks
		break
	case BeforeStart:
		hooks = a.beforeStartHooks
		break
	case BeforeStop:
		hooks = a.beforeStopHooks
		break
	default:
		return fmt.Errorf("not found hook type")
	}

	for _, hook := range hooks {
		err := hook(a.dipper, c)
		if err != nil {
			return err
		}
	}

	return nil
}
