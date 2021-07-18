package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/global"
)

func getStatus(rpcAddr string, isLog bool) {
	url := fmt.Sprintf("%s/block", rpcAddr)
	resp, err := http.Get(url)
	if err != nil {
		common.GetLogger().Error("[node-status] Unable to connect to ", url)
		return
	}
	defer resp.Body.Close()

	type RPCTempResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Block struct {
				Header struct {
					Chainid string `json:"chain_id"`
					Height  string `json:"height"`
					Time    string `json:"time"`
				} `json:"header"`
			} `json:"block"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}

	result := new(RPCTempResponse)
	if json.NewDecoder(resp.Body).Decode(result) != nil {
		common.GetLogger().Error("[node-status] Unexpected response: ", url)
		return
	}

	global.Mutex.Lock()
	common.NodeStatus.Chainid = result.Result.Block.Header.Chainid
	common.NodeStatus.Block, _ = strconv.ParseInt(result.Result.Block.Header.Height, 10, 64)
	common.NodeStatus.Blocktime = result.Result.Block.Header.Time
	global.Mutex.Unlock()

	if isLog {
		common.GetLogger().Info("[node-status] (new block) height: ", common.NodeStatus.Block, " time: ", common.NodeStatus.Blocktime)
	}

	// save block height/time
	blockTime, _ := time.Parse(time.RFC3339, result.Result.Block.Header.Time)
	database.AddBlockTime(common.NodeStatus.Block, blockTime.Unix())
	database.AddBlockNanoTime(common.NodeStatus.Block, blockTime.UnixNano())
	common.AddNewBlock(common.NodeStatus.Block, blockTime.UnixNano())
}

// SyncStatus is a function for syncing sekaid status.
func SyncStatus(rpcAddr string, isLog bool) {
	for {
		getStatus(rpcAddr, isLog)

		if isLog {
			common.GetLogger().Info("[node-status] Syncing node status")
			common.GetLogger().Info("[node-status] Chain_id = ", common.NodeStatus.Chainid)
			common.GetLogger().Info("[node-status] Block = ", common.NodeStatus.Block)
			common.GetLogger().Info("[node-status] Blocktime = ", common.NodeStatus.Blocktime)
		}

		time.Sleep(time.Duration(config.Config.Block.StatusSync) * time.Second)
	}
}
