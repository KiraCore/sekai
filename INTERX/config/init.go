package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getGetMethods() []string {
	return []string{
		QueryAccounts,
		QueryTotalSupply,
		QueryBalances,
		QueryTransactionHash,
		QueryDataReferenceKeys,
		QueryDataReference,
		QueryKiraStatus,
		QueryWithdraws,
		QueryDeposits,
		QueryStatus,
		QueryConsensus,
		QueryValidators,
		QueryValidatorInfos,
		QueryBlocks,
		QueryBlockByHeightOrHash,
		QueryBlockTransactions,
		QueryTransactionResult,
		QueryProposals,
		QueryProposal,
		QueryKiraTokensAliases,
		QueryKiraTokensRates,
		QueryVoters,
		QueryVotes,
		QueryKiraTokensAliases,
		QueryKiraTokensRates,

		QueryRosettaNetworkList,
		QueryRosettaNetworkOptions,
		QueryRosettaNetworkStatus,
		QueryRosettaAccountBalance,
	}
}

func getPostMethods() []string {
	return []string{
		PostTransaction,
		EncodeTransaction,
	}
}

func defaultConfig() InterxConfigFromFile {
	configFromFile := InterxConfigFromFile{}

	configFromFile.ServeHTTPS = false
	configFromFile.GRPC = "dns:///0.0.0.0:9090"
	configFromFile.RPC = "http://0.0.0.0:26657"
	configFromFile.PORT = "11000"

	configFromFile.MnemonicFile = LoadMnemonic("swap exercise equip shoot mad inside floor wheel loan visual stereo build frozen always bulb naive subway foster marine erosion shuffle flee action there")

	configFromFile.Cache.StatusSync = 5
	configFromFile.Cache.CacheDir = "cache"
	configFromFile.Cache.MaxCacheSize = "2GB"
	configFromFile.Cache.CachingDuration = 5
	configFromFile.Cache.DownloadFileSizeLimitation = "10MB"

	configFromFile.Faucet.MnemonicFile = LoadMnemonic("equip exercise shoot mad inside floor wheel loan visual stereo build frozen potato always bulb naive subway foster marine erosion shuffle flee action there")
	configFromFile.Faucet.FaucetAmounts = make(map[string]int64)
	configFromFile.Faucet.FaucetAmounts["stake"] = 100000
	configFromFile.Faucet.FaucetAmounts["validatortoken"] = 100000
	configFromFile.Faucet.FaucetAmounts["ukex"] = 100000
	configFromFile.Faucet.FaucetMinimumAmounts = make(map[string]int64)
	configFromFile.Faucet.FaucetMinimumAmounts["stake"] = 100
	configFromFile.Faucet.FaucetMinimumAmounts["validatortoken"] = 100
	configFromFile.Faucet.FaucetMinimumAmounts["ukex"] = 100
	configFromFile.Faucet.FeeAmounts = make(map[string]string)
	configFromFile.Faucet.FeeAmounts["stake"] = "1000ukex"
	configFromFile.Faucet.FeeAmounts["validatortoken"] = "1000ukex"
	configFromFile.Faucet.FeeAmounts["ukex"] = "1000ukex"
	configFromFile.Faucet.TimeLimit = 20

	defaultRPCSetting := RPCSetting{
		Disable:         false,
		RateLimit:       0,
		AuthRateLimit:   0,
		CachingDisable:  false,
		CachingDuration: 30,
	}

	configFromFile.RPCMethods.API = make(map[string]map[string]RPCSetting)
	configFromFile.RPCMethods.API["GET"] = make(map[string]RPCSetting)
	configFromFile.RPCMethods.API["POST"] = make(map[string]RPCSetting)
	for _, item := range getGetMethods() {
		configFromFile.RPCMethods.API["GET"][item] = defaultRPCSetting
	}
	for _, item := range getPostMethods() {
		configFromFile.RPCMethods.API["POST"][item] = defaultRPCSetting
	}

	return configFromFile
}

// InitConfig is a function to load interx configurations from a given file
func InitConfig(
	configFilePath string,
	serveHTTPS bool,
	grpc string,
	rpc string,
	port string,
	signingMnemonic string,
	syncStatus int64,
	cacheDir string,
	maxCacheSize string,
	cachingDuration int64,
	maxDownloadSize string,
	faucetMnemonic string,
	faucetTimeLimit int64,
	faucetAmounts string,
	faucetMinimumAmounts string,
	feeAmounts string,
) {
	configFromFile := defaultConfig()

	configFromFile.ServeHTTPS = serveHTTPS
	configFromFile.GRPC = grpc
	configFromFile.RPC = rpc
	configFromFile.PORT = port
	configFromFile.MnemonicFile = LoadMnemonic(signingMnemonic)

	configFromFile.Cache.StatusSync = syncStatus
	configFromFile.Cache.CacheDir = cacheDir
	configFromFile.Cache.MaxCacheSize = maxCacheSize
	configFromFile.Cache.CachingDuration = cachingDuration
	configFromFile.Cache.DownloadFileSizeLimitation = maxDownloadSize

	configFromFile.Faucet.MnemonicFile = LoadMnemonic(faucetMnemonic)
	configFromFile.Faucet.TimeLimit = faucetTimeLimit

	configFromFile.Faucet.FaucetAmounts = make(map[string]int64)
	for _, amount := range strings.Split(faucetAmounts, ",") {
		coin, err := sdk.ParseCoinNormalized(amount)
		if err == nil {
			configFromFile.Faucet.FaucetAmounts[coin.Denom] = coin.Amount.Int64()
		}
	}

	configFromFile.Faucet.FaucetMinimumAmounts = make(map[string]int64)
	for _, amount := range strings.Split(faucetMinimumAmounts, ",") {
		coin, err := sdk.ParseCoinNormalized(amount)
		if err == nil {
			configFromFile.Faucet.FaucetMinimumAmounts[coin.Denom] = coin.Amount.Int64()
		}
	}

	configFromFile.Faucet.FeeAmounts = make(map[string]string)
	for _, denom_amount := range strings.Split(feeAmounts, ",") {
		denom := strings.Split(denom_amount, " ")[0]
		amount := strings.Split(denom_amount, " ")[1]
		configFromFile.Faucet.FeeAmounts[denom] = amount
	}

	bytes, err := json.MarshalIndent(&configFromFile, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(configFilePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
