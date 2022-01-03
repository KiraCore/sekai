package common

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// GetSha256SumFromBytes is a function to get hash
func GetSha256SumFromBytes(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// GetBlake2bHash is a function to get hash
func GetBlake2bHash(request interface{}) string {
	// Calculate blake2b hash
	requestJSON, err := json.Marshal(request)
	if err != nil {
		GetLogger().Error("[blake2b-hash] Unable to marshal request: ", err)
	}

	return GetSha256SumFromBytes(requestJSON)
}

// GetMD5Hash is a function to get hash
func GetMD5Hash(request interface{}) string {
	// Calculate md5 hash
	requestJSON, err := json.Marshal(request)
	if err != nil {
		GetLogger().Error("[md5-hash] Unable to marshal request: ", err)
	}

	hash := md5.Sum([]byte(requestJSON))
	return fmt.Sprintf("%X", hash)
}
