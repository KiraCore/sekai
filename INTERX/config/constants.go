package config

const (
	QueryAccounts        = "/api/cosmos/auth/accounts/{address}"
	QueryTotalSupply     = "/api/cosmos/bank/supply"
	QueryBalances        = "/api/cosmos/bank/balances/{address}"
	PostTransaction      = "/api/cosmos/txs"
	QueryTransactionHash = "/api/cosmos/txs/{hash}"
	EncodeTransaction    = "/api/cosmos/txs/encode"

	QueryProposals         = "/api/kira/gov/proposals"
	QueryProposal          = "/api/kira/gov/proposals/{proposal_id}"
	QueryVoters            = "/api/kira/gov/voters/{proposal_id}"
	QueryVotes             = "/api/kira/gov/votes/{proposal_id}"
	QueryDataReferenceKeys = "/api/kira/gov/data_keys"
	QueryDataReference     = "/api/kira/gov/data/{key}"
	QueryKiraTokensAliases = "/api/kira/tokens/aliases"
	QueryKiraTokensRates   = "/api/kira/tokens/rates"
	QueryKiraFunctions     = "/api/kira/metadata"
	QueryKiraStatus        = "/api/kira/status"

	QueryInterxFunctions = "/api/metadata"

	FaucetRequestURL         = "/api/faucet"
	QueryRPCMethods          = "/api/rpc_methods"
	QueryWithdraws           = "/api/withdraws"
	QueryDeposits            = "/api/deposits"
	QueryBlocks              = "/api/blocks"
	QueryBlockByHeightOrHash = "/api/blocks/{height}"
	QueryBlockTransactions   = "/api/blocks/{height}/transactions"
	QueryTransactionResult   = "/api/transactions/{txHash}"
	QueryStatus              = "/api/status"
	QueryValidators          = "/api/valopers"
	QueryValidatorInfos      = "/api/valoperinfos"
	QueryGenesis             = "/api/genesis"
	QueryGenesisSum          = "/api/gensum"

	Download = "/download"

	DataReferenceRegistry = "DRR"
)
