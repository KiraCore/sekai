package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sekaiapp "github.com/KiraCore/sekai/app"
	middleware "github.com/KiraCore/sekai/middleware"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bytesize "github.com/inhies/go-bytesize"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/sr25519"
	"github.com/tyler-smith/go-bip39"
)

func readConfig() InterxConfig {
	sekaiapp.SetConfig()
	middleware.RegisterStdMsgs()

	type ConfigFromFile struct {
		Mnemonic        string `json:"mnemonic"`
		StatusSync      int64  `json:"status_sync"`
		CacheDir        string `json:"cache_dir"`
		MaxCacheSize    string `json:"max_cache_size"`
		CachingDuration int64  `json:"caching_duration"`
		Faucet          struct {
			Mnemonic             string           `json:"mnemonic"`
			FaucetAmounts        map[string]int64 `json:"faucet_amounts"`
			FaucetMinimumAmounts map[string]int64 `json:"faucet_minimum_amounts"`
			TimeLimit            int64            `json:"time_limit"`
		} `json:"faucet"`
		RPC RPCConfig `json:"rpc"`
	}

	file, _ := ioutil.ReadFile("./config.json")

	configFromFile := ConfigFromFile{}

	err := json.Unmarshal([]byte(file), &configFromFile)
	if err != nil {
		fmt.Println("Invalid configuration: {}", err)
		panic(err)
	}

	config := InterxConfig{}

	// Interx Main Configuration
	config.Mnemonic = configFromFile.Mnemonic
	config.StatusSync = configFromFile.StatusSync
	config.CacheDir = configFromFile.CacheDir
	b, _ := bytesize.Parse(configFromFile.MaxCacheSize)
	config.MaxCacheSize = int64(b)
	config.CachingDuration = configFromFile.CachingDuration
	config.PrivKey = secp256k1.GenPrivKeyFromSecret(bip39.NewSeed(config.Mnemonic, ""))
	config.PubKey = config.PrivKey.PubKey()
	config.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), config.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Interx Mnemonic   : ", config.Mnemonic)
	fmt.Println("Interx Address    : ", config.Address)
	fmt.Println("Interx Public Key : ", sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, config.PubKey))
	fmt.Println("Max Cache Size    : ", config.MaxCacheSize)
	fmt.Println("Caching Duration  : ", config.CachingDuration)

	// Faucet Configuration
	config.Faucet = FaucetConfig{
		Mnemonic:             configFromFile.Faucet.Mnemonic,
		FaucetAmounts:        configFromFile.Faucet.FaucetAmounts,
		FaucetMinimumAmounts: configFromFile.Faucet.FaucetMinimumAmounts,
		TimeLimit:            configFromFile.Faucet.TimeLimit,
	}

	config.Faucet.PrivKey = secp256k1.GenPrivKeyFromSecret(bip39.NewSeed(config.Faucet.Mnemonic, ""))
	config.Faucet.PubKey = config.Faucet.PrivKey.PubKey()
	config.Faucet.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), config.Faucet.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Faucet Mnemonic   : ", config.Faucet.Mnemonic)
	fmt.Println("Faucet Address    : ", config.Faucet.Address)
	fmt.Println("Faucet Public Key : ", sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, config.Faucet.PubKey))

	// RPC Configuration
	config.RPC = configFromFile.RPC

	return config
}

// GenPrivKey is a function to generate a privKey
func GenPrivKey() crypto.PrivKey {
	return sr25519.GenPrivKey()
}

var (
	// Config is a configuration for interx
	Config = readConfig()
	// EncodingCg is a configuration for Amino Encoding
	EncodingCg = sekaiapp.MakeEncodingConfig()
)
