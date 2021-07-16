package config

import (
	"github.com/KiraCore/sekai/INTERX/types"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/tendermint/p2p"
)

// FaucetConfig is a struct to be used for Faucet configuration
type FaucetConfig struct {
	Mnemonic             string            `json:"mnemonic"`
	FaucetAmounts        map[string]string `json:"faucet_amounts"`
	FaucetMinimumAmounts map[string]string `json:"faucet_minimum_amounts"`
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

type BlockConfig struct {
	StatusSync          int64 `json:"status_sync"`
	HaltedAvgBlockTimes int64 `json:"halted_avg_block_times"`
}

// CacheConfig is a struct to be used for cache configuration
type CacheConfig struct {
	CacheDir                   string `json:"cache_dir"`
	MaxCacheSize               int64  `json:"max_cache_size"`
	CachingDuration            int64  `json:"caching_duration"`
	DownloadFileSizeLimitation int64  `json:"download_file_size_limitation"`
}

type NodeDiscoveryConfig struct {
	UseHttps              bool   `json:"use_https"`
	DefaultInterxPort     string `json:"default_interx_port"`
	DefaultTendermintPort string `json:"default_tendermint_port"`
	ConnectionTimeout     string `json:"connection_timeout"`
}

// InterxConfig is a struct to be used for interx configuration
type InterxConfig struct {
	Version       string              `json:"version"`
	ServeHTTPS    bool                `json:"serve_https"`
	GRPC          string              `json:"grpc"`
	RPC           string              `json:"rpc"`
	PORT          string              `json:"port"`
	Node          types.NodeConfig    `json:"node"`
	Mnemonic      string              `json:"mnemonic"`
	AddrBooks     []string            `json:"addrbooks"`
	NodeKey       *p2p.NodeKey        `json:"node_key"`
	TxModes       []string            `json:"tx_modes"`
	PrivKey       crypto.PrivKey      `json:"privkey"`
	PubKey        crypto.PubKey       `json:"pubkey"`
	Address       string              `json:"address"`
	NodeDiscovery NodeDiscoveryConfig `json:"node_discovery"`
	Block         BlockConfig         `json:"block"`
	Cache         CacheConfig         `json:"cache"`
	Faucet        FaucetConfig        `json:"faucet"`
	RPCMethods    RPCConfig           `json:"rpc_methods"`
}

// InterxConfigFromFile is a struct to be used for interx configuration file
type InterxConfigFromFile struct {
	Version       string              `json:"version"`
	ServeHTTPS    bool                `json:"serve_https"`
	GRPC          string              `json:"grpc"`
	RPC           string              `json:"rpc"`
	PORT          string              `json:"port"`
	Node          types.NodeConfig    `json:"node"`
	MnemonicFile  string              `json:"mnemonic"`
	AddrBooks     string              `json:"addrbooks"`
	NodeKey       string              `json:"node_key"`
	TxModes       string              `json:"tx_modes"`
	Block         BlockConfig         `json:"block"`
	NodeDiscovery NodeDiscoveryConfig `json:"node_discovery"`
	Cache         struct {
		CacheDir                   string `json:"cache_dir"`
		MaxCacheSize               string `json:"max_cache_size"`
		CachingDuration            int64  `json:"caching_duration"`
		DownloadFileSizeLimitation string `json:"download_file_size_limitation"`
	} `json:"cache"`
	Faucet struct {
		MnemonicFile         string            `json:"mnemonic"`
		FaucetAmounts        map[string]string `json:"faucet_amounts"`
		FaucetMinimumAmounts map[string]string `json:"faucet_minimum_amounts"`
		FeeAmounts           map[string]string `json:"fee_amounts"`
		TimeLimit            int64             `json:"time_limit"`
	} `json:"faucet"`
	RPCMethods RPCConfig `json:"rpc_methods"`
}
