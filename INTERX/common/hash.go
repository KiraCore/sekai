package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

// GetBlake2bHashFromBytes is a function to get hash
func GetBlake2bHashFromBytes(data []byte) string {
	hash := blake2b.Sum256(data)
	return fmt.Sprintf("%X", hash)
}

// GetBlake2bHash is a function to get hash
func GetBlake2bHash(request interface{}) string {
	// Calculate blake2b hash
	requestJSON, err := json.Marshal(request)
	if err != nil {
		GetLogger().Error("[blake2b-hash] Unable to marshal request: ", err)
	}

	return GetBlake2bHashFromBytes(requestJSON)
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
