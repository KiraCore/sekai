package config

const (
	QueryAccounts        = "/api/cosmos/auth/accounts"
	QueryTotalSupply     = "/api/cosmos/bank/supply"
	QueryBalances        = "/api/cosmos/bank/balances"
	PostTransaction      = "/api/cosmos/txs"
	QueryTransactionHash = "/api/cosmos/txs"
	EncodeTransaction    = "/api/cosmos/txs/encode"

	QueryDataReferenceKeys = "/api/kira/gov/data_keys"
	QueryDataReference     = "/api/kira/gov/data"
	QueryKiraFunctions     = "/api/kira/metadata"
	QueryKiraStatus        = "/api/kira/status"

	FaucetRequestURL     = "/api/faucet"
	QueryRPCMethods      = "/api/rpc_methods"
	QueryInterxFunctions = "/api/metadata"
	QueryWithdraws       = "/api/withdraws"
	QueryDeposits        = "/api/deposits"
	QueryStatus          = "/api/status"
	QueryValidators      = "/api/valopers"
	QueryGenesis         = "/api/genesis"
	QueryGenesisSum      = "/api/gensum"

	Download = "/download"

	DataReferenceRegistry = "DRR"
)
