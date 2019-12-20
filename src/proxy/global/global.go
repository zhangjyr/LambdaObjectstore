package global

import (
	"github.com/cornelk/hashmap"
	"github.com/wangaoone/LambdaObjectstore/lib/logger"
	"sync"

	protocol "github.com/wangaoone/LambdaObjectstore/src/types"
	"github.com/wangaoone/LambdaObjectstore/src/proxy/types"
)

var (
	// Clients        = make([]chan interface{}, 1024*1024)
	DataCollected    sync.WaitGroup
	Log              logger.ILogger
	ReqMap           = hashmap.New(1024)
	Migrator         types.MigrationScheduler
	BasePort         = 6378
	BaseMigratorPort = 6380
	ServerIp         string
	Prefix           string
	Flags            uint64
)

func init() {
	Log = logger.NilLogger

	ip, err := GetPrivateIp()
	if err != nil {
		panic(err)
	}

	ServerIp = ip
	Flags = protocol.FLAG_WARMUP_FIXED_INTERVAL | protocol.FLAG_REPLICATE_ON_WARMUP
}

func IsWarmupWithFixedInterval() bool {
	return Flags & protocol.FLAG_WARMUP_FIXED_INTERVAL > 0
}
