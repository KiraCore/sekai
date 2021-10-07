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
		QueryNetworkProperties,

		QueryRosettaNetworkList,
		QueryRosettaNetworkOptions,
		QueryRosettaNetworkStatus,
		QueryRosettaAccountBalance,

		QueryPrivP2PList,
		QueryPubP2PList,
		QueryInterxList,
		QuerySnapList,

		QueryCurrentPlan,
		QueryNextPlan,

		QueryAddrBook,
		QueryNetInfo,

		QueryIdentityRecord,
		QueryIdentityRecordsByAddress,
		QueryAllIdentityRecords,
		QueryIdentityRecordVerifyRequest,
		QueryIdentityRecordVerifyRequestsByRequester,
		QueryIdentityRecordVerifyRequestsByApprover,
		QueryAllIdentityRecordVerifyRequests,
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

	configFromFile.Version = "0.1.0"
	configFromFile.ServeHTTPS = false
	configFromFile.GRPC = "dns:///0.0.0.0:9090"
	configFromFile.RPC = "http://0.0.0.0:26657"
	configFromFile.PORT = "11000"

	configFromFile.Node.NodeType = "seed"
	configFromFile.Node.SentryNodeID = ""
	configFromFile.Node.SnapshotNodeID = ""
	configFromFile.Node.ValidatorNodeID = ""
	configFromFile.Node.SeedNodeID = ""

	configFromFile.MnemonicFile = LoadMnemonic("swap exercise equip shoot mad inside floor wheel loan visual stereo build frozen always bulb naive subway foster marine erosion shuffle flee action there")

	configFromFile.AddrBooks = "addrbook.json"
	configFromFile.NodeKey = "node_key.json"
	configFromFile.TxModes = "sync,async,block"

	configFromFile.Block.StatusSync = 5
	configFromFile.Block.HaltedAvgBlockTimes = 10

	configFromFile.NodeDiscovery.UseHttps = false
	configFromFile.NodeDiscovery.DefaultInterxPort = "11000"
	configFromFile.NodeDiscovery.DefaultTendermintPort = "26657"
	configFromFile.NodeDiscovery.ConnectionTimeout = "3s"

	configFromFile.Cache.CacheDir = "cache"
	configFromFile.Cache.MaxCacheSize = "2GB"
	configFromFile.Cache.CachingDuration = 5
	configFromFile.Cache.DownloadFileSizeLimitation = "10MB"

	configFromFile.Faucet.MnemonicFile = LoadMnemonic("equip exercise shoot mad inside floor wheel loan visual stereo build frozen potato always bulb naive subway foster marine erosion shuffle flee action there")
	configFromFile.Faucet.FaucetAmounts = make(map[string]string)
	configFromFile.Faucet.FaucetAmounts["stake"] = "100000"
	configFromFile.Faucet.FaucetAmounts["validatortoken"] = "100000"
	configFromFile.Faucet.FaucetAmounts["ukex"] = "100000"
	configFromFile.Faucet.FaucetMinimumAmounts = make(map[string]string)
	configFromFile.Faucet.FaucetMinimumAmounts["stake"] = "100"
	configFromFile.Faucet.FaucetMinimumAmounts["validatortoken"] = "100"
	configFromFile.Faucet.FaucetMinimumAmounts["ukex"] = "100"
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
	version string,
	configFilePath string,
	serveHTTPS bool,
	grpc string,
	rpc string,
	nodeType string,
	sentryNodeId string,
	snapshotNodeId string,
	validatorNodeId string,
	seedNodeId string,
	port string,
	signingMnemonic string,
	syncStatus int64,
	haltedAvgBlockTimes int64,
	cacheDir string,
	maxCacheSize string,
	cachingDuration int64,
	maxDownloadSize string,
	faucetMnemonic string,
	faucetTimeLimit int64,
	faucetAmounts string,
	faucetMinimumAmounts string,
	feeAmounts string,
	addrBooks string,
	txModes string,
	nodeDiscoveryUseHttps bool,
	nodeDiscoveryInterxPort string,
	nodeDiscoveryTendermintPort string,
	nodeDiscoveryTimeout string,
	nodeKey string,
) {
	configFromFile := defaultConfig()

	configFromFile.Version = version
	configFromFile.ServeHTTPS = serveHTTPS
	configFromFile.GRPC = grpc
	configFromFile.RPC = rpc
	configFromFile.PORT = port

	configFromFile.Node.NodeType = nodeType
	configFromFile.Node.SentryNodeID = sentryNodeId
	configFromFile.Node.SnapshotNodeID = snapshotNodeId
	configFromFile.Node.ValidatorNodeID = validatorNodeId
	configFromFile.Node.SeedNodeID = seedNodeId

	configFromFile.MnemonicFile = LoadMnemonic(signingMnemonic)

	configFromFile.AddrBooks = addrBooks
	configFromFile.NodeKey = nodeKey
	configFromFile.TxModes = txModes

	configFromFile.NodeDiscovery.UseHttps = nodeDiscoveryUseHttps
	configFromFile.NodeDiscovery.DefaultInterxPort = nodeDiscoveryInterxPort
	configFromFile.NodeDiscovery.DefaultTendermintPort = nodeDiscoveryTendermintPort
	configFromFile.NodeDiscovery.ConnectionTimeout = nodeDiscoveryTimeout

	configFromFile.Block.StatusSync = syncStatus
	configFromFile.Block.HaltedAvgBlockTimes = haltedAvgBlockTimes

	configFromFile.Cache.CacheDir = cacheDir
	configFromFile.Cache.MaxCacheSize = maxCacheSize
	configFromFile.Cache.CachingDuration = cachingDuration
	configFromFile.Cache.DownloadFileSizeLimitation = maxDownloadSize

	configFromFile.Faucet.MnemonicFile = LoadMnemonic(faucetMnemonic)
	configFromFile.Faucet.TimeLimit = faucetTimeLimit

	configFromFile.Faucet.FaucetAmounts = make(map[string]string)
	for _, amount := range strings.Split(faucetAmounts, ",") {
		coin, err := sdk.ParseCoinNormalized(amount)
		if err == nil {
			configFromFile.Faucet.FaucetAmounts[coin.Denom] = coin.Amount.String()
		}
	}

	configFromFile.Faucet.FaucetMinimumAmounts = make(map[string]string)
	for _, amount := range strings.Split(faucetMinimumAmounts, ",") {
		coin, err := sdk.ParseCoinNormalized(amount)
		if err == nil {
			configFromFile.Faucet.FaucetMinimumAmounts[coin.Denom] = coin.Amount.String()
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
