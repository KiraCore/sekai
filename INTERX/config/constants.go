package config

const (
	QueryAccounts        = "/api/cosmos/auth/accounts/{address}"
	QueryTotalSupply     = "/api/cosmos/bank/supply"
	QueryBalances        = "/api/cosmos/bank/balances/{address}"
	PostTransaction      = "/api/cosmos/txs"
	QueryTransactionHash = "/api/cosmos/txs/{hash}"
	EncodeTransaction    = "/api/cosmos/txs/encode"

	QueryDataReferenceKeys = "/api/kira/gov/data_keys"
	QueryDataReference     = "/api/kira/gov/data/{key}"
	QueryKiraFunctions     = "/api/kira/metadata"
	QueryKiraStatus        = "/api/kira/status"

	FaucetRequestURL         = "/api/faucet"
	QueryRPCMethods          = "/api/rpc_methods"
	QueryInterxFunctions     = "/api/metadata"
	QueryWithdraws           = "/api/withdraws"
	QueryDeposits            = "/api/deposits"
	QueryBlocks              = "/api/blocks"
	QueryBlockByHeightOrHash = "/api/blocks/{height}"
	QueryBlockTransactions   = "/api/blocks/{height}/transactions"
	QueryTransactionResult   = "/api/transactions/{txHash}"
	QueryStatus              = "/api/status"
	QueryValidators          = "/api/valopers"
	QueryGenesis             = "/api/genesis"
	QueryGenesisSum          = "/api/gensum"

	Download = "/download"

	DataReferenceRegistry = "DRR"
)
