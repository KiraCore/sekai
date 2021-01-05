package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	functions "github.com/KiraCore/sekai/INTERX/functions"
	sekaiapp "github.com/KiraCore/sekai/app"
	functionmeta "github.com/KiraCore/sekai/function_meta"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bytesize "github.com/inhies/go-bytesize"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/sr25519"
	"github.com/tyler-smith/go-bip39"
)

func parseSizeString(size string) int64 {
	b, _ := bytesize.Parse(size)
	return int64(b)
}

func mnemonicFromFile(filename string) string {
	if len(filename) == 0 {
		return ""
	}

	mnemonic, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	return string(mnemonic)
}

func loadMnemonic(mnemonic string) string {
	if bip39.IsMnemonicValid(mnemonic) {
		return mnemonic
	}

	return mnemonicFromFile(mnemonic)
}

func readConfig() InterxConfig {
	functions.RegisterInterxFunctions()
	functionmeta.RegisterStdMsgs()
	sekaiapp.SetConfig()

	type ConfigFromFile struct {
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
	config.Mnemonic = loadMnemonic(configFromFile.MnemonicFile)
	config.Cache.StatusSync = configFromFile.Cache.StatusSync
	config.Cache.CacheDir = configFromFile.Cache.CacheDir
	config.Cache.MaxCacheSize = parseSizeString(configFromFile.Cache.MaxCacheSize)
	config.Cache.CachingDuration = configFromFile.Cache.CachingDuration
	config.Cache.DownloadFileSizeLimitation = parseSizeString(configFromFile.Cache.DownloadFileSizeLimitation)

	if !bip39.IsMnemonicValid(config.Mnemonic) {
		fmt.Println("Invalid Interx Mnemonic: ", config.Mnemonic)
		panic("Invalid Interx Mnemonic")
	}
	config.PrivKey = secp256k1.GenPrivKeyFromSecret(bip39.NewSeed(config.Mnemonic, ""))
	config.PubKey = config.PrivKey.PubKey()
	config.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), config.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Interx Mnemonic   : ", config.Mnemonic)
	fmt.Println("Interx Address    : ", config.Address)
	fmt.Println("Interx Public Key : ", sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, config.PubKey))

	fmt.Println("Max Cache Size    : ", config.Cache.MaxCacheSize)
	fmt.Println("Caching Duration  : ", config.Cache.CachingDuration)
	fmt.Println("Download File Size Limitation  : ", config.Cache.DownloadFileSizeLimitation)

	// Faucet Configuration
	config.Faucet = FaucetConfig{
		Mnemonic:             loadMnemonic(configFromFile.Faucet.MnemonicFile),
		FaucetAmounts:        configFromFile.Faucet.FaucetAmounts,
		FaucetMinimumAmounts: configFromFile.Faucet.FaucetMinimumAmounts,
		FeeAmounts:           configFromFile.Faucet.FeeAmounts,
		TimeLimit:            configFromFile.Faucet.TimeLimit,
	}

	if !bip39.IsMnemonicValid(config.Faucet.Mnemonic) {
		fmt.Println("Invalid Faucet Mnemonic: ", config.Faucet.Mnemonic)
		panic("Invalid Faucet Mnemonic")
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

	if _, err := os.Stat(config.Cache.CacheDir); os.IsNotExist(err) {
		os.Mkdir(config.Cache.CacheDir, os.ModePerm)
	}
	if _, err := os.Stat(config.Cache.CacheDir + "/reference/"); os.IsNotExist(err) {
		os.Mkdir(config.Cache.CacheDir+"/reference/", os.ModePerm)
	}
	if _, err := os.Stat(config.Cache.CacheDir + "/response/"); os.IsNotExist(err) {
		os.Mkdir(config.Cache.CacheDir+"/response/", os.ModePerm)
	}
	if _, err := os.Stat(config.Cache.CacheDir + "/db/"); os.IsNotExist(err) {
		os.Mkdir(config.Cache.CacheDir+"/db/", os.ModePerm)
	}

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

// GetReferenceCacheDir is a function to get reference directory
func GetReferenceCacheDir() string {
	return Config.Cache.CacheDir + "/reference/"
}

// GetResponseCacheDir is a function to get reference directory
func GetResponseCacheDir() string {
	return Config.Cache.CacheDir + "/response/"
}

// GetDbCacheDir is a function to get db directory
func GetDbCacheDir() string {
	return Config.Cache.CacheDir + "/db/"
}
