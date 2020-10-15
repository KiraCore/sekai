package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sekaiapp "github.com/KiraCore/sekai/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/sr25519"
	"github.com/tyler-smith/go-bip39"
)

func readConfig() InterxConfig {
	sekaiapp.SetConfig()

	type ConfigFromFile struct {
		Mnemonic string `json:"mnemonic"`
		Faucet   struct {
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
	}

	config := InterxConfig{}

	// Interx Main Configuration
	config.Mnemonic = configFromFile.Mnemonic
	config.PrivKey = secp256k1.GenPrivKeyFromSecret(bip39.NewSeed(config.Mnemonic, ""))
	config.PubKey = config.PrivKey.PubKey()
	config.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), config.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Interx Mnemonic   : ", config.Mnemonic)
	fmt.Println("Interx Address    : ", config.Address)
	fmt.Println("Interx Public Key : ", sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, config.PubKey))

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
