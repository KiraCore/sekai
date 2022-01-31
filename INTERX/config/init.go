package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getPostMethods() []string {
	return []string{
		PostTransaction,
		EncodeTransaction,
	}
}

func getRPCSettings() RPCConfig {
	config := RPCConfig{}

	defaultRPCSetting := RPCSetting{
		Disable:              false,
		RateLimit:            0,
		AuthRateLimit:        0,
		CachingDisable:       false,
		CachingDuration:      30,
		CachingBlockDuration: 5,
	}

	config.API = make(map[string]map[string]RPCSetting)
	config.API["GET"] = make(map[string]RPCSetting)
	config.API["POST"] = make(map[string]RPCSetting)

	// endpoints that can change within 1 block time
	defaultRPCSetting.CachingDuration = -1
	defaultRPCSetting.CachingBlockDuration = 1
	config.API["GET"][QueryAccounts] = defaultRPCSetting
	config.API["GET"][QueryTotalSupply] = defaultRPCSetting
	config.API["GET"][QueryBalances] = defaultRPCSetting
	config.API["GET"][QueryAccounts] = defaultRPCSetting
	config.API["GET"][QueryDataReferenceKeys] = defaultRPCSetting
	config.API["GET"][QueryDataReference] = defaultRPCSetting
	config.API["GET"][QueryKiraStatus] = defaultRPCSetting
	config.API["GET"][QueryWithdraws] = defaultRPCSetting
	config.API["GET"][QueryDeposits] = defaultRPCSetting
	config.API["GET"][QueryStatus] = defaultRPCSetting
	config.API["GET"][QueryValidators] = defaultRPCSetting
	config.API["GET"][QueryValidatorInfos] = defaultRPCSetting
	config.API["GET"][QueryRoles] = defaultRPCSetting
	config.API["GET"][QueryRolesByAddress] = defaultRPCSetting
	config.API["GET"][QueryPermissionsByAddress] = defaultRPCSetting
	config.API["GET"][QueryProposals] = defaultRPCSetting
	config.API["GET"][QueryProposal] = defaultRPCSetting
	config.API["GET"][QueryKiraTokensAliases] = defaultRPCSetting
	config.API["GET"][QueryKiraTokensRates] = defaultRPCSetting
	config.API["GET"][QueryVoters] = defaultRPCSetting
	config.API["GET"][QueryVotes] = defaultRPCSetting
	config.API["GET"][QueryKiraTokensAliases] = defaultRPCSetting
	config.API["GET"][QueryKiraTokensRates] = defaultRPCSetting
	config.API["GET"][QueryNetworkProperties] = defaultRPCSetting
	config.API["GET"][QueryRosettaNetworkList] = defaultRPCSetting
	config.API["GET"][QueryRosettaNetworkOptions] = defaultRPCSetting
	config.API["GET"][QueryRosettaNetworkStatus] = defaultRPCSetting
	config.API["GET"][QueryRosettaAccountBalance] = defaultRPCSetting
	config.API["GET"][QueryCurrentPlan] = defaultRPCSetting
	config.API["GET"][QueryNextPlan] = defaultRPCSetting
	config.API["GET"][QueryAddrBook] = defaultRPCSetting
	config.API["GET"][QueryNetInfo] = defaultRPCSetting
	config.API["GET"][QueryIdentityRecord] = defaultRPCSetting
	config.API["GET"][QueryIdentityRecordsByAddress] = defaultRPCSetting
	config.API["GET"][QueryAllIdentityRecords] = defaultRPCSetting
	config.API["GET"][QueryIdentityRecordVerifyRequest] = defaultRPCSetting
	config.API["GET"][QueryIdentityRecordVerifyRequestsByRequester] = defaultRPCSetting
	config.API["GET"][QueryIdentityRecordVerifyRequestsByApprover] = defaultRPCSetting
	config.API["GET"][QueryAllIdentityRecordVerifyRequests] = defaultRPCSetting

	// endpoints that never change
	defaultRPCSetting.CachingDuration = -1
	defaultRPCSetting.CachingBlockDuration = -1
	config.API["GET"][QueryTransactionHash] = defaultRPCSetting
	config.API["GET"][QueryBlocks] = defaultRPCSetting
	config.API["GET"][QueryBlockByHeightOrHash] = defaultRPCSetting
	config.API["GET"][QueryTransactionResult] = defaultRPCSetting

	// endpoints that change constantly
	defaultRPCSetting.CachingDuration = 1
	defaultRPCSetting.CachingBlockDuration = -1
	config.API["GET"][QueryConsensus] = defaultRPCSetting
	config.API["GET"][QueryPrivP2PList] = defaultRPCSetting
	config.API["GET"][QueryPubP2PList] = defaultRPCSetting
	config.API["GET"][QueryInterxList] = defaultRPCSetting
	config.API["GET"][QuerySnapList] = defaultRPCSetting

	for _, item := range getPostMethods() {
		config.API["POST"][item] = defaultRPCSetting
	}

	return config
}

func defaultConfig() InterxConfigFromFile {
	configFromFile := InterxConfigFromFile{}

	configFromFile.Version = InterxVersion
	configFromFile.ServeHTTPS = false
	configFromFile.GRPC = "dns:///0.0.0.0:9090"
	configFromFile.RPC = "http://0.0.0.0:26657"
	configFromFile.PORT = "11000"

	configFromFile.Node.NodeType = "seed"
	configFromFile.Node.SentryNodeID = ""
	configFromFile.Node.SnapshotNodeID = ""
	configFromFile.Node.ValidatorNodeID = ""
	configFromFile.Node.SeedNodeID = ""

	configFromFile.MnemonicFile = "interx.mnemonic"

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

	configFromFile.Faucet.MnemonicFile = configFromFile.MnemonicFile

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

	configFromFile.MnemonicFile = signingMnemonic

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

	configFromFile.Faucet.MnemonicFile = faucetMnemonic
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
