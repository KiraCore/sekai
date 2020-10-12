package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sekaiapp "github.com/KiraCore/sekai/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
)

// FaucetConfigFromFile is a struct to be used for Faucet configuration
type FaucetConfigFromFile struct {
	Mnemonic             string           `json:"mnemonic"`
	FaucetAmounts        map[string]int64 `json:"faucet_amounts"`
	FaucetMinimumAmounts map[string]int64 `json:"faucet_minimum_amounts"`
	TimeLimit            int64            `json:"time_limit"`
}

// FaucetConfig is a struct to be used for Faucet configuration
type FaucetConfig struct {
	Mnemonic             string           `json:"mnemonic"`
	FaucetAmounts        map[string]int64 `json:"faucet_amounts"`
	FaucetMinimumAmounts map[string]int64 `json:"faucet_minimum_amounts"`
	TimeLimit            int64            `json:"time_limit"`
	PrivKey              crypto.PrivKey   `json:"privkey"`
	PubKey               crypto.PubKey    `json:"pubkey"`
	Address              string           `json:"address"`
}

func readFaucetConfig() FaucetConfig {
	sekaiapp.SetConfig()

	file, _ := ioutil.ReadFile("./faucet.json")

	configFromFile := FaucetConfigFromFile{}

	err := json.Unmarshal([]byte(file), &configFromFile)
	if err != nil {
		fmt.Println("Invalid configuration: {}", err)
	}

	config := FaucetConfig{
		Mnemonic:             configFromFile.Mnemonic,
		FaucetAmounts:        configFromFile.FaucetAmounts,
		FaucetMinimumAmounts: configFromFile.FaucetMinimumAmounts,
		TimeLimit:            configFromFile.TimeLimit,
	}

	seed := bip39.NewSeed(config.Mnemonic, "")
	config.PrivKey = secp256k1.GenPrivKeyFromSecret(seed)
	config.PubKey = config.PrivKey.PubKey()
	config.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), config.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Faucet Mnemonic   : ", config.Mnemonic)
	fmt.Println("Faucet Address    : ", config.Address)
	fmt.Println("Faucet Public Key : ", sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, config.PubKey))

	return config
}

var (
	// FaucetCg is a configuration for faucet server
	FaucetCg FaucetConfig = readFaucetConfig()
)
