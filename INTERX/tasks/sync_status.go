package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
)

// NodeStatus is a struct to be used for node status
var NodeStatus struct {
	Chainid   string `json:"chain_id"`
	Block     int64  `json:"block"`
	Blocktime string `json:"block_time"`
}

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
	NodeStatus.Chainid = result.Result.Block.Header.Chainid
	NodeStatus.Block, _ = strconv.ParseInt(result.Result.Block.Header.Height, 10, 64)
	NodeStatus.Blocktime = result.Result.Block.Header.Time
	common.Mutex.Unlock()
}

// SyncStatus is a function for syncing sekaid status.
func SyncStatus(rpcAddr string) {
	for {
		getStatus(rpcAddr)
		fmt.Println("\nsync node status	: ")
		fmt.Println("	chain id	: ", NodeStatus.Chainid)
		fmt.Println("	block		: ", NodeStatus.Block)
		fmt.Println("	blocktime	: ", NodeStatus.Blocktime)

		time.Sleep(time.Duration(interx.Config.StatusSync) * time.Second)
	}
}
