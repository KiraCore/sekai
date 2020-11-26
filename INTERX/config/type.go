package config

import (
	"github.com/tendermint/tendermint/crypto"
)

// FaucetConfig is a struct to be used for Faucet configuration
type FaucetConfig struct {
	Mnemonic             string            `json:"mnemonic"`
	FaucetAmounts        map[string]int64  `json:"faucet_amounts"`
	FaucetMinimumAmounts map[string]int64  `json:"faucet_minimum_amounts"`
	FeeAmounts           map[string]string `json:"fee_amounts"`
	TimeLimit            int64             `json:"time_limit"`
	PrivKey              crypto.PrivKey    `json:"privkey"`
	PubKey               crypto.PubKey     `json:"pubkey"`
	Address              string            `json:"address"`
}

// RPCSetting is a struct to be used for endpoint setting
type RPCSetting struct {
	Disable         bool    `json:"disable"`
	RateLimit       float64 `json:"rate_limit,omitempty"`
	AuthRateLimit   float64 `json:"auth_rate_limit,omitempty"`
	CachingDisable  bool    `json:"caching_disable"`
	CachingDuration int64   `json:"caching_duration,omitempty"`
}

// RPCConfig is a struct to be used for PRC configuration
type RPCConfig struct {
	API map[string]map[string]RPCSetting `json:"API"`
}

// InterxConfig is a struct to be used for interx configuration
type InterxConfig struct {
	Mnemonic        string         `json:"mnemonic"`
	PrivKey         crypto.PrivKey `json:"privkey"`
	PubKey          crypto.PubKey  `json:"pubkey"`
	Address         string         `json:"address"`
	StatusSync      int64          `json:"status_sync"`
	CacheDir        string         `json:"cache_dir"`
	MaxCacheSize    int64          `json:"max_cache_size"`
	CachingDuration int64          `json:"caching_duration"`
	Faucet          FaucetConfig   `json:"faucet"`
	RPC             RPCConfig      `json:"rpc"`
}
