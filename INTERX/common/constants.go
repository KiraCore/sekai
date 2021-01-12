package common

const (
	QueryAccounts          = "/api/cosmos/auth/accounts/{address}"
	QueryTotalSupply       = "/api/cosmos/bank/supply"
	QueryBalances          = "/api/cosmos/bank/balances"
	PostTransaction        = "/api/cosmos/txs"
	QueryTransactionHash   = "/api/cosmos/txs"
	EncodeTransaction      = "/api/cosmos/txs/encode"
	FaucetRequestURL       = "/api/faucet"
	QueryRPCMethods        = "/api/rpc_methods"
	QueryInterxFunctions   = "/api/metadata"
	QueryStatus            = "/api/status"
	QueryWithdraws         = "/api/withdraws"
	QueryDeposits          = "/api/deposits"
	QueryDataReferenceKeys = "/api/kira/gov/data_keys"
	QueryDataReference     = "/api/kira/gov/data"
	QueryKiraFunctions     = "/api/kira/metadata"
	Download               = "/download"
	DataReferenceRegistry  = "DRR"
	QueryValidators        = "/api/validators"
)
