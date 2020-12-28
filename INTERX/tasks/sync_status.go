package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/database"
)

func getStatus(rpcAddr string) {
	resp, err := http.Get(fmt.Sprintf("%s/block", rpcAddr))
	if err != nil {
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
		return
	}

	common.Mutex.Lock()
	common.NodeStatus.Chainid = result.Result.Block.Header.Chainid
	common.NodeStatus.Block, _ = strconv.ParseInt(result.Result.Block.Header.Height, 10, 64)
	common.NodeStatus.Blocktime = result.Result.Block.Header.Time
	common.Mutex.Unlock()

	// save block height/time
	t, _ := time.Parse(time.RFC3339, common.NodeStatus.Blocktime)
	database.AddBlockTime(common.NodeStatus.Block, t.Unix())
}

// SyncStatus is a function for syncing sekaid status.
func SyncStatus(rpcAddr string, isLog bool) {
	for {
		getStatus(rpcAddr)

		if isLog {
			fmt.Println("\nsync node status	: ")
			fmt.Println("	chain id	: ", common.NodeStatus.Chainid)
			fmt.Println("	block		: ", common.NodeStatus.Block)
			fmt.Println("	blocktime	: ", common.NodeStatus.Blocktime)
		}

		time.Sleep(time.Duration(interx.Config.StatusSync) * time.Second)
	}
}
