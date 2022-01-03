package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/KiraCore/sekai/INTERX/types"
	sekaiapp "github.com/KiraCore/sekai/app"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bytesize "github.com/inhies/go-bytesize"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/sr25519"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tyler-smith/go-bip39"
)

var (
	// Config is a configuration for interx
	Config = InterxConfig{}
	// EncodingCg is a configuration for Amino Encoding
	EncodingCg = sekaiapp.MakeEncodingConfig()
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

func TrimMnemonic(mnemonic string) string {
	return strings.Join(strings.Fields(mnemonic), " ")
}

// LoadMnemonic is a function to extract mnemonic
func LoadMnemonic(mnemonic string) string {
	if bip39.IsMnemonicValid(mnemonic) {
		return TrimMnemonic(mnemonic)
	}

	return TrimMnemonic(mnemonicFromFile(mnemonic))
}

// LoadConfig is a function to load interx configurations from a given file
func LoadConfig(configFilePath string) {
	sekaiapp.SetConfig()

	Config = InterxConfig{}

	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Invalid configuration: {}", err)
		panic(err)
	}

	configFromFile := InterxConfigFromFile{}

	err = json.Unmarshal([]byte(file), &configFromFile)
	if err != nil {
		fmt.Println("Invalid configuration: {}", err)
		panic(err)
	}

	// Interx Main Configuration
	Config.Version = configFromFile.Version
	Config.ServeHTTPS = configFromFile.ServeHTTPS
	Config.GRPC = configFromFile.GRPC
	Config.RPC = configFromFile.RPC
	Config.PORT = configFromFile.PORT
	Config.Mnemonic = LoadMnemonic(configFromFile.MnemonicFile)

	Config.Node = configFromFile.Node

	fmt.Println("Interx GRPC: ", Config.GRPC)
	fmt.Println("Interx RPC : ", Config.RPC)
	fmt.Println("Interx PORT: ", Config.PORT)

	if !bip39.IsMnemonicValid(Config.Mnemonic) {
		fmt.Println("Invalid Interx Mnemonic: ", Config.Mnemonic)
		panic("Invalid Interx Mnemonic")
	}

	seed, err := bip39.NewSeedWithErrorChecking(Config.Mnemonic, "")
	if err != nil {
		panic(err)
	}
	master, ch := hd.ComputeMastersFromSeed(seed)
	priv, err := hd.DerivePrivateKeyForPath(master, ch, "44'/118'/0'/0/0")
	if err != nil {
		panic(err)
	}

	Config.PrivKey = &secp256k1.PrivKey{Key: priv}
	Config.PubKey = Config.PrivKey.PubKey()
	Config.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), Config.PubKey.Address())

	Config.AddrBooks = strings.Split(configFromFile.AddrBooks, ",")
	Config.NodeKey, err = p2p.LoadOrGenNodeKey(configFromFile.NodeKey)
	if err != nil {
		panic(err)
	}

	Config.TxModes = strings.Split(configFromFile.TxModes, ",")
	if len(Config.TxModes) == 0 {
		Config.TxModes = strings.Split("sync,async,block", ",")
	}

	// Display mnemonic and keys
	fmt.Println("Interx Mnemonic  : ", Config.Mnemonic)
	fmt.Println("Interx Address   : ", Config.Address)
	fmt.Println("Interx Public Key: ", Config.PubKey.String())

	Config.NodeDiscovery = configFromFile.NodeDiscovery

	Config.Block.StatusSync = configFromFile.Block.StatusSync
	Config.Block.HaltedAvgBlockTimes = configFromFile.Block.HaltedAvgBlockTimes

	Config.Cache.CacheDir = configFromFile.Cache.CacheDir
	Config.Cache.MaxCacheSize = parseSizeString(configFromFile.Cache.MaxCacheSize)
	Config.Cache.CachingDuration = configFromFile.Cache.CachingDuration
	Config.Cache.DownloadFileSizeLimitation = parseSizeString(configFromFile.Cache.DownloadFileSizeLimitation)

	// Display cache configurations
	fmt.Println("Interx Block StatusSync                : ", Config.Block.StatusSync)
	fmt.Println("Halted Avg Block Times                 : ", Config.Block.HaltedAvgBlockTimes)

	fmt.Println("Interx Cache CacheDir                  : ", Config.Cache.CacheDir)
	fmt.Println("Interx Cache MaxCacheSize              : ", Config.Cache.MaxCacheSize)
	fmt.Println("Interx Cache CachingDuration           : ", Config.Cache.CachingDuration)
	fmt.Println("Interx Cache DownloadFileSizeLimitation: ", Config.Cache.DownloadFileSizeLimitation)

	// Faucet Configuration
	Config.Faucet = FaucetConfig{
		Mnemonic:             LoadMnemonic(configFromFile.Faucet.MnemonicFile),
		FaucetAmounts:        configFromFile.Faucet.FaucetAmounts,
		FaucetMinimumAmounts: configFromFile.Faucet.FaucetMinimumAmounts,
		FeeAmounts:           configFromFile.Faucet.FeeAmounts,
		TimeLimit:            configFromFile.Faucet.TimeLimit,
	}

	if !bip39.IsMnemonicValid(Config.Faucet.Mnemonic) {
		fmt.Println("Invalid Faucet Mnemonic: ", Config.Faucet.Mnemonic)
		panic("Invalid Faucet Mnemonic")
	}

	seed, err = bip39.NewSeedWithErrorChecking(Config.Faucet.Mnemonic, "")
	if err != nil {
		panic(err)
	}
	master, ch = hd.ComputeMastersFromSeed(seed)
	priv, err = hd.DerivePrivateKeyForPath(master, ch, "44'/118'/0'/0/0")
	if err != nil {
		panic(err)
	}

	Config.Faucet.PrivKey = &secp256k1.PrivKey{Key: priv}
	Config.Faucet.PubKey = Config.Faucet.PrivKey.PubKey()
	Config.Faucet.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), Config.Faucet.PubKey.Address())

	// Display mnemonic and keys
	fmt.Println("Faucet Mnemonic  : ", Config.Faucet.Mnemonic)
	fmt.Println("Faucet Address   : ", Config.Faucet.Address)
	fmt.Println("Faucet Public Key: ", Config.Faucet.PubKey.String())

	// Faucet configurations
	fmt.Println("Interx Faucet FaucetAmounts       : ", Config.Faucet.FaucetAmounts)
	fmt.Println("Interx Faucet FaucetMinimumAmounts: ", Config.Faucet.FaucetMinimumAmounts)
	fmt.Println("Interx Faucet FeeAmounts          : ", Config.Faucet.FeeAmounts)
	fmt.Println("Interx Faucet TimeLimit           : ", Config.Faucet.TimeLimit)

	// RPC Configuration
	Config.RPCMethods = configFromFile.RPCMethods

	if _, err := os.Stat(Config.Cache.CacheDir); os.IsNotExist(err) {
		os.Mkdir(Config.Cache.CacheDir, os.ModePerm)
	}
	if _, err := os.Stat(Config.Cache.CacheDir + "/reference/"); os.IsNotExist(err) {
		os.Mkdir(Config.Cache.CacheDir+"/reference/", os.ModePerm)
	}
	if _, err := os.Stat(Config.Cache.CacheDir + "/response/"); os.IsNotExist(err) {
		os.Mkdir(Config.Cache.CacheDir+"/response/", os.ModePerm)
	}
	if _, err := os.Stat(Config.Cache.CacheDir + "/db/"); os.IsNotExist(err) {
		os.Mkdir(Config.Cache.CacheDir+"/db/", os.ModePerm)
	}
	if _, err := os.Stat(GetReferenceCacheDir() + "/genesis.json"); !os.IsNotExist(err) {
		os.Remove(GetReferenceCacheDir() + "/genesis.json")
	}

	if _, err := os.Stat(GetDbCacheDir() + "/token-aliases.json"); !os.IsNotExist(err) {
		os.Remove(GetDbCacheDir() + "/token-aliases.json")
	}
}

// GenPrivKey is a function to generate a privKey
func GenPrivKey() crypto.PrivKey {
	return sr25519.GenPrivKey()
}

// GetReferenceCacheDir is a function to get reference directory
func GetReferenceCacheDir() string {
	return Config.Cache.CacheDir + "/reference"
}

// GetResponseCacheDir is a function to get reference directory
func GetResponseCacheDir() string {
	return Config.Cache.CacheDir + "/response"
}

// GetDbCacheDir is a function to get db directory
func GetDbCacheDir() string {
	return Config.Cache.CacheDir + "/db"
}

func LoadAddressBooks() []types.AddrBookJSON {
	addrBooks := make([]types.AddrBookJSON, 0)
	for _, addrFile := range Config.AddrBooks {
		file, _ := ioutil.ReadFile(addrFile)

		book := types.AddrBookJSON{}

		err := json.Unmarshal([]byte(file), &book)

		if err != nil {
			fmt.Println("Failed to load addrBook: ", addrFile)
		}

		addrBooks = append(addrBooks, book)
	}

	return addrBooks
}

func LoadUniqueIPAddresses() []string {
	ipAddresses := make([]string, 0)

	flag := make(map[string]bool)
	for _, addrFile := range Config.AddrBooks {
		file, _ := ioutil.ReadFile(addrFile)

		book := types.AddrBookJSON{}

		err := json.Unmarshal([]byte(file), &book)

		if err != nil {
			fmt.Println("Failed to load addrBook: ", addrFile)
		}

		for _, addr := range book.Addrs {
			if _, ok := flag[addr.Addr.IP]; !ok {
				ipAddresses = append(ipAddresses, addr.Addr.IP)
			}
			flag[addr.Addr.IP] = true
		}
	}

	return ipAddresses
}
