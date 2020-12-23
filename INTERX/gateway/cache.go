package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
)

// PutCache is a function to save value to cache
func PutCache(chainIDHash string, endpointHash string, requestHash string, value common.InterxResponse) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	folderPath := fmt.Sprintf("%s/%s/%s", interx.GetResponseCacheDir(), chainIDHash, endpointHash)
	filePath := fmt.Sprintf("%s/%s", folderPath, requestHash)

	common.Mutex.Lock()
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		common.Mutex.Unlock()
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	common.Mutex.Unlock()

	return err
}

// GetCache is a function to get value from cache
func GetCache(chainIDHash string, endpointHash string, requestHash string) (common.InterxResponse, error) {
	filePath := fmt.Sprintf("%s/%s/%s/%s", interx.GetResponseCacheDir(), chainIDHash, endpointHash, requestHash)

	common.Mutex.Lock()
	data, _ := ioutil.ReadFile(filePath)
	common.Mutex.Unlock()

	response := common.InterxResponse{}
	err := json.Unmarshal([]byte(data), &response)

	return response, err
}
