package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// WhitelistConfig is a struct to be used for PRC Whitelist configuration
type WhitelistConfig struct {
	Disable       bool    `json:"disable"`
	RateLimit     float64 `json:"rate_limit,omitempty"`
	AuthRateLimit float64 `json:"auth_rate_limit,omitempty"`
}

func readWhitelistConfig() map[string]map[string]WhitelistConfig {
	file, _ := ioutil.ReadFile("./whitelist.json")

	config := map[string]map[string]WhitelistConfig{}

	err := json.Unmarshal([]byte(file), &config)
	if err != nil {
		fmt.Println("Invalid configuration : {}", err)
	}

	return config
}

var (
	// WhitelistCg is a configuration for rpc whitelist
	WhitelistCg = readWhitelistConfig()
)
