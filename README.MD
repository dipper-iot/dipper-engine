# Rule Engine Golang
[![Go](https://github.com/dipper-iot/dipper-engine/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/dipper-iot/dipper-engine/actions/workflows/go.yml) [![CodeQL](https://github.com/dipper-iot/dipper-engine/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/dipper-iot/dipper-engine/actions/workflows/codeql.yml) [![Coverage Status](https://coveralls.io/repos/github/dipper-iot/dipper-engine/badge.svg?branch=main)](https://coveralls.io/github/dipper-iot/dipper-engine?branch=main)

<!-- TOC -->

  * [Setup](#setup)
  * [Run](#run)
  * [Rule Engine](#rule-engine)
  * [Example Developer Test](#example-developer-test)
<!-- TOC -->

## Setup

```shell
$ go build
```

## Run
```shell
$ dipper-engine -h

NAME:
   Dipper Engine - Rule Engine

USAGE:
   Dipper Engine [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --bus value, -b value       Bus type: [chan, redis] (default: "chan")
   --config value, -c value    Config file (default: "config.json")
   --help, -h                  show help (default: false)
   --plugin, -p                Load Plugin (default: false)
   --queue value, -q value     Queue type: [chan, redis] (default: "chan")
   --redis-host value          Redis host (default: "127.0.0.1:6379")
   --redis-pass value          Redis pass
   --redis-timeout value       Redis Time out
   --session-from-queue, --sq  Session from queue (default: false)
   --store value, -s value     Store type: [memory, redis] (default: "memory")


```

```shell
$ dipper-engine

Start Dipper Engine

Rules: 5
-----------------------------------------------------------
No  Rule Name                  Worker  Infinity  Status  
1   arithmetic                 1       false     enable  
2   conditional                1       false     enable  
3   fork                       1       false     enable  
4   input-redis-queue          0       true      disable 
5   input-redis-queue-extend   1       true      enable  
6   log-core                   1       false     enable  
7   output-redis-queue         0       false     disable 
8   output-redis-queue-extend  1       false     enable  
9   switch                     1       false     enable
-----------------------------------------------------------

Running Engine...

```

## Rule Engine

|  No   | Rule                               | Description                      | Infinity | Doc |
|:-----:|:-----------------------------------|:---------------------------------|:--------:|-----|
|   1   | arithmetic                         | operator match rule              |  false   |     |
|   2   | conditional                        | compare data rule                |  false   |     |
|   3   | fork                               | fork to rules                    |  false   |     |
|   4   | input-redis-queue                  | input data from queue on config  |   true   |     |
|   5   | input-redis-queue-extend           | input data from queue on option  |   true   |     |
|   6   | log-core                           | log to console                   |  false   |     |
|   7   | output-redis-queue                 | output data from queue on config |  false   |     |
|   8   | output-redis-queue-extend          | output data from queue on option |  false   |     |
|   9   | switch                             | switch to rules                  |  false   |     |

## Example Developer Test

```golang

c := engine.New()

c.Hook(engine.AfterStart, func(dipper *core.DipperEngine, c *cli.Context) error {
    return dipper.Add(context.Background(), &data.Session{
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
                Debug: true,
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
                    "next_success": "3",
                    "debug":        true,
                },
                NodeId: "4",
                RuleId: "arithmetic",
                End:    false,
            },
            "2": {
                Debug: true,
                Option: map[string]interface{}{
                    "debug": true,
                },
                NodeId: "2",
                RuleId: "log-core",
                End:    true,
            },
            "3": {
                Debug: true,
                Option: map[string]interface{}{
                    "next_success": []string{"5", "2"},
                    "debug":        true,
                },
                NodeId: "3",
                RuleId: "fork",
                End:    false,
            },
            "5": {
                Debug: true,
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
                        "operator": "==",
                        "type":     "operator",
                    },
                    "set_param_result_to": "default.cond_a_b",
                    "next_error":          "2",
                    "next_true":           "2",
                    "next_false":          "2",
                    "debug":               true,
                },
                NodeId: "5",
                RuleId: "conditional",
                End:    false,
            },
		},
    })
})

if err := c.Run(os.Args); err != nil {
    log.Fatal(err)
}
```

