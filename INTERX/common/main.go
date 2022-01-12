package common

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	"google.golang.org/grpc/grpclog"
)

// RPCMethods is a variable for rpc methods
var RPCMethods = make(map[string]map[string]types.RPCMethod)

// AddRPCMethod is a function to add a RPC method
func AddRPCMethod(method string, url string, description string, canCache bool) {
	newMethod := types.RPCMethod{}
	newMethod.Description = description
	newMethod.Enabled = true
	newMethod.CachingEnabled = true

	if conf, ok := config.Config.RPCMethods.API[method][url]; ok {
		newMethod.Enabled = !conf.Disable
		newMethod.CachingEnabled = !conf.CachingDisable
		newMethod.RateLimit = conf.RateLimit
		newMethod.AuthRateLimit = conf.AuthRateLimit
		newMethod.CachingDuration = conf.CachingDuration
		newMethod.CachingBlockDuration = conf.CachingBlockDuration
	}

	if !canCache {
		newMethod.CachingEnabled = false
	}

	if _, ok := RPCMethods[method]; !ok {
		RPCMethods[method] = map[string]types.RPCMethod{}
	}
	RPCMethods[method][url] = newMethod
}

var logger = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)

// GetLogger is a function to get logger
func GetLogger() grpclog.LoggerV2 {
	return logger
}

// NodeStatus is a struct to be used for node status
var NodeStatus struct {
	Chainid   string `json:"chain_id"`
	Block     int64  `json:"block"`
	Blocktime string `json:"block_time"`
}

func IsCacheExpired(result types.InterxResponse) bool {
	isBlockExpire := false
	if result.CachingBlockDuration == 0 {
		isBlockExpire = true
	} else if result.CachingBlockDuration == -1 {
		isBlockExpire = false
	} else if result.Response.Block+result.CachingBlockDuration > NodeStatus.Block {
		isBlockExpire = false
	} else {
		isBlockExpire = true
	}

	isTimestampExpire := false
	if result.CachingDuration == 0 {
		isTimestampExpire = true
	} else if result.CachingDuration == -1 {
		isTimestampExpire = false
	} else if result.CacheTime.Add(time.Duration(result.CachingDuration) * time.Second).After(time.Now().UTC()) {
		isTimestampExpire = false
	} else {
		isTimestampExpire = true
	}

	return isBlockExpire || isTimestampExpire
}
