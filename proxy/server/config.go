package server

import (
	"time"

	"github.com/mason-leap-lab/infinicache/proxy/lambdastore"
)

const LambdaMaxDeployments = 400
const NumLambdaClusters = 400
const LambdaStoreName = "LambdaStore"
const LambdaPrefix = "Store1VPCNode"
const InstanceWarmTimout = 1 * time.Minute
const InstanceCapacity = 1024 * 1000000    // MB
const InstanceOverhead = 100 * 1000000     // MB

func init() {
	lambdastore.WarmTimout = InstanceWarmTimout
}
