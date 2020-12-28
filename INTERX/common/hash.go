package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

// GetBlake2bHash is a function to get hash
func GetBlake2bHash(request interface{}) string {
	// Calculate blake2b hash
	requestJSON, _ := json.Marshal(request)
	hash := blake2b.Sum256([]byte(requestJSON))
	return fmt.Sprintf("%X", hash)
}

// GetMD5Hash is a function to get hash
func GetMD5Hash(request interface{}) string {
	// Calculate md5 hash
	requestJSON, _ := json.Marshal(request)
	hash := md5.Sum([]byte(requestJSON))
	return fmt.Sprintf("%X", hash)
}
