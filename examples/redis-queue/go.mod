module github.com/dipper-iot/dipper-engine/examples/redis-queue

go 1.18

require (
	github.com/dipper-iot/dipper-engine v0.0.0-20221022050351-f7e9b09096a6
	github.com/go-redis/redis/v9 v9.0.0-rc.1
	github.com/sirupsen/logrus v1.9.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/sony/sonyflake v1.1.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)

replace ithub.com/dipper-iot/dipper-engine => ../../
