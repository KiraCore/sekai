package config

import (
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
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
	RateLimit       float64 `json:"rate_limit"`
	AuthRateLimit   float64 `json:"auth_rate_limit"`
	CachingDisable  bool    `json:"caching_disable"`
	CachingDuration int64   `json:"caching_duration"`
}

// RPCConfig is a struct to be used for PRC configuration
type RPCConfig struct {
	API map[string]map[string]RPCSetting `json:"API"`
}

// CacheConfig is a struct to be used for cache configuration
type CacheConfig struct {
	StatusSync                 int64  `json:"status_sync"`
	CacheDir                   string `json:"cache_dir"`
	MaxCacheSize               int64  `json:"max_cache_size"`
	CachingDuration            int64  `json:"caching_duration"`
	DownloadFileSizeLimitation int64  `json:"download_file_size_limitation"`
}

// InterxConfig is a struct to be used for interx configuration
type InterxConfig struct {
	ServeHTTPS bool           `json:"serve_https"`
	GRPC       string         `json:"grpc"`
	RPC        string         `json:"rpc"`
	PORT       string         `json:"port"`
	Mnemonic   string         `json:"mnemonic"`
	PrivKey    crypto.PrivKey `json:"privkey"`
	PubKey     crypto.PubKey  `json:"pubkey"`
	Address    string         `json:"address"`
	Cache      CacheConfig    `json:"cache"`
	Faucet     FaucetConfig   `json:"faucet"`
	RPCMethods RPCConfig      `json:"rpc_methods"`
}

// InterxConfigFromFile is a struct to be used for interx configuration file
type InterxConfigFromFile struct {
	ServeHTTPS   bool   `json:"serve_https"`
	GRPC         string `json:"grpc"`
	RPC          string `json:"rpc"`
	PORT         string `json:"port"`
	MnemonicFile string `json:"mnemonic"`
	Cache        struct {
		StatusSync                 int64  `json:"status_sync"`
		CacheDir                   string `json:"cache_dir"`
		MaxCacheSize               string `json:"max_cache_size"`
		CachingDuration            int64  `json:"caching_duration"`
		DownloadFileSizeLimitation string `json:"download_file_size_limitation"`
	} `json:"cache"`
	Faucet struct {
		MnemonicFile         string            `json:"mnemonic"`
		FaucetAmounts        map[string]int64  `json:"faucet_amounts"`
		FaucetMinimumAmounts map[string]int64  `json:"faucet_minimum_amounts"`
		FeeAmounts           map[string]string `json:"fee_amounts"`
		TimeLimit            int64             `json:"time_limit"`
	} `json:"faucet"`
	RPCMethods RPCConfig `json:"rpc_methods"`
}
