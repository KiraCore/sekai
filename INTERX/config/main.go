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

// InterxConfigFromFile is a struct to be used for interx configuration
type InterxConfigFromFile struct {
	Mnemonic string `json:"mnemonic"`
}

// InterxConfig is a struct to be used for interx configuration
type InterxConfig struct {
	Mnemonic string         `json:"mnemonic"`
	PrivKey  crypto.PrivKey `json:"privkey"`
	PubKey   crypto.PubKey  `json:"pubkey"`
	Address  string         `json:"address"`
}

func readInterxConfig() InterxConfig {
	file, _ := ioutil.ReadFile("./config.json")

	configFromFile := InterxConfigFromFile{}

	err := json.Unmarshal([]byte(file), &configFromFile)
	if err != nil {
		fmt.Println("Invalid configuration: {}", err)
	}

	config := InterxConfig{
		Mnemonic: configFromFile.Mnemonic,
	}

	seed := bip39.NewSeed(config.Mnemonic, "")
	config.PrivKey = secp256k1.GenPrivKeyFromSecret(seed)
	config.PubKey = config.PrivKey.PubKey()
	config.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), config.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Interx Mnemonic   : ", config.Mnemonic)
	fmt.Println("Interx Address    : ", config.Address)
	fmt.Println("Interx Public Key : ", sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, config.PubKey))

	return config
}

// GenPrivKey is a function to generate a privKey
func GenPrivKey() crypto.PrivKey {
	return sr25519.GenPrivKey()
}

var (
	// InterxCg is a configuration for rpc whitelist
	InterxCg = readInterxConfig()
	// EncodingCg is a configuration for Amino Encoding
	EncodingCg = sekaiapp.MakeEncodingConfig()
)
