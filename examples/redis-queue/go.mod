module github.com/dipper-iot/dipper-engine/examples/redis-queue

go 1.18

require (
	github.com/dipper-iot/dipper-engine v0.0.0-20221022050351-f7e9b09096a6
	github.com/go-redis/redis/v8 v8.11.5
	github.com/sirupsen/logrus v1.9.0
	github.com/urfave/cli/v2 v2.11.2
)

require (
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sony/sonyflake v1.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)

replace github.com/dipper-iot/dipper-engine => ../../
