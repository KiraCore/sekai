package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	interx "github.com/KiraCore/sekai/INTERX/config"
)

// PutCache is a function to save value to cache
func PutCache(chainIDHash string, endpointHash string, requestHash string, value InterxResponse) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	folderPath := fmt.Sprintf("%s/%s/%s", interx.Config.RPC.CacheDir, chainIDHash, endpointHash)
	filePath := fmt.Sprintf("%s/%s", folderPath, requestHash)

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetCache is a function to get value from cache
func GetCache(chainIDHash string, endpointHash string, requestHash string) (InterxResponse, error) {
	filePath := fmt.Sprintf("%s/%s/%s/%s", interx.Config.RPC.CacheDir, chainIDHash, endpointHash, requestHash)

	data, _ := ioutil.ReadFile(filePath)

	response := InterxResponse{}
	err := json.Unmarshal([]byte(data), &response)

	return response, err
}

// RemoveCache is a function to get value from cache
func RemoveCache(chainIDHash string, endpointHash string, requestHash string) {
	filePath := fmt.Sprintf("%s/%s/%s/%s", interx.Config.RPC.CacheDir, chainIDHash, endpointHash, requestHash)

	_ = os.Remove(filePath)
}
