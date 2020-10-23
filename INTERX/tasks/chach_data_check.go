package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
)

// CacheDataCheck is a function to check cache data if it's expired.
func CacheDataCheck(rpcAddr string) {
	for {
		err := filepath.Walk(interx.Config.RPC.CacheDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// check timestamp and block height in cache data and delete expired one
				delete := false

				if !info.IsDir() && info.Size() != 0 {
					common.Mutex.Lock()
					data, _ := ioutil.ReadFile(path)
					common.Mutex.Unlock()

					result := common.InterxResponse{}
					err := json.Unmarshal([]byte(data), &result)

					if err == nil && result.ExpireAt.Before(time.Now()) && result.Response.Block != NodeStatus.Block {
						delete = true
					}
				}

				if path != interx.Config.RPC.CacheDir && delete {
					fmt.Println("deleting file ... ", path)

					common.Mutex.Lock()
					err := os.Remove(path)
					common.Mutex.Unlock()

					if err != nil {
						fmt.Println("Error deleting file: ", err)
						return err
					}

					return nil
				}

				return nil
			})

		if err != nil {
			log.Println(err)
		}
	}
}
