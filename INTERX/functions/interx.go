package functions

import (
	"encoding/json"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/iancoleman/strcase"
)

type InterxFunctionParameter struct {
	Type        string                    `json:"type,omitempty"`
	Optional    bool                      `json:"optional,omitempty"`
	Description string                    `json:"description"`
	Fields      *InterxFunctionParameters `json:"fields,omitempty"`
}

type InterxFunctionParameters = map[string]InterxFunctionParameter

type InterxFunctionMeta struct {
	Endpoint    string                   `json:"endpoint"`
	Description string                   `json:"description"`
	Parameters  InterxFunctionParameters `json:"parameters"`
	Response    InterxFunctionParameters `json:"response"`
}

type InterxMetadata struct {
	Functions      map[string]InterxFunctionMeta `json:"functions"`
	ResponseHeader InterxFunctionParameters      `json:"response_header"`
}

var (
	interxMetadata InterxMetadata = InterxMetadata{
		Functions:      make(map[string]InterxFunctionMeta),
		ResponseHeader: InterxFunctionParameters{},
	}
)

// AddInterxFunction is a function to add a function
func AddInterxFunction(functionType string, endpoint string, meta string) {
	metadata := InterxFunctionMeta{}
	if err := json.Unmarshal([]byte(meta), &metadata); err != nil {
		panic(err)
	}
	metadata.Endpoint = endpoint

	interxMetadata.Functions[strcase.ToCamel(functionType)] = metadata
}

// RegisterInterxFunctions is a function to register all interx functions
func RegisterInterxFunctions() {
	interxMetadata.ResponseHeader["Interx_chain_id"] = InterxFunctionParameter{
		Type:        "number",
		Description: "This represents the current chain id.",
	}
	interxMetadata.ResponseHeader["Interx_block"] = InterxFunctionParameter{
		Type:        "number",
		Description: "This represents the current block number.",
	}
	interxMetadata.ResponseHeader["Interx_blocktime"] = InterxFunctionParameter{
		Type:        "number",
		Description: "This represents the current block timestamp.",
	}
	interxMetadata.ResponseHeader["Interx_timestamp"] = InterxFunctionParameter{
		Type:        "string",
		Description: "This represents the current interx timestamp.",
	}
	interxMetadata.ResponseHeader["Interx_request_hash"] = InterxFunctionParameter{
		Type:        "string",
		Description: "This represents the hash of request parameters.",
	}
	interxMetadata.ResponseHeader["Interx_signature"] = InterxFunctionParameter{
		Type:        "string",
		Description: "This represents the interx response signature.",
	}
	interxMetadata.ResponseHeader["Interx_hash"] = InterxFunctionParameter{
		Type:        "string",
		Description: "This represents the interx response hash.",
	}
	interxMetadata.ResponseHeader["Interx_ref"] = InterxFunctionParameter{
		Type:        "string",
		Description: "This represents link to download the data reference.",
	}

	AddInterxFunction(
		"QueryKiraStatus",
		config.QueryKiraStatus,
		`{
			"description": "QueryKiraStatus is a function to query the node status",
			"response": {
				"node_info": {
					"description": "The connected node information"
				},
				"sync_info": {
					"description": "The sync status of connected node"
				},
				"validator_info": {
					"description": "The validator information of connect node"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryAccount",
		config.QueryAccounts,
		`{
			"description": "QueryAccount is a function to query the account info.",
			"parameters": {
				"address": {
					"type":        "string",
					"description": "This represents the account address."
				}
			},
			"response": {
				"account": {
					"description": "The account info with address, pubkey and sequence."
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryTotalSupply",
		config.QueryTotalSupply,
		`{
			"description": "QueryTotalSupply is a function to query total supply.",
			"response": {
				"supply": {
					"type": "Coin[]",
					"description": "The total supply of the network"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryBalance",
		config.QueryBalances,
		`{
			"description": "QueryBalance is a function to query the account balances.",
			"parameters": {
				"address": {
					"type":        "string",
					"description": "This represents the account address."
				},
				"limit": {
					"type":        "number",
					"description": "This represents the page size"
				},
				"offset": {
					"type":        "number",
					"description": "This represents the page number"
				},
				"count_total": {
					"type":        "number",
					"description": "This represents the option to return total count of data reference keys.",
					"optional": true
				}
			},
			"response": {
				"balances": {
					"type": "Coin[]",
					"description": "The account balances with pagination"
				},
				"pagination": {
					"description": "The pagination response information like total and next_key"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryTransactionHash",
		config.QueryTransactionHash,
		`{
			"description": "QueryTransactionHash is a function to query transaction details from transaction hash.",
			"parameters": {
				"hash": {
					"type":        "string",
					"description": "This represents the transaction hash. (e.g. 0x20.....)"
				}
			},
			"response": {
				"hash": {
					"description": "The transaction hash"
				},
				"height": {
					"description": "The block height of transation"
				},
				"tx": {
					"description": "The base-64 encoded transaction"
				},
				"tx_result": {
					"description": "The result of transaction with events, gas info, logs and error code"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryDataReferenceKeys",
		config.QueryDataReferenceKeys,
		`{
			"description": "QueryDataReferenceKeys is a function to query data reference keys with pagination.",
			"parameters": {
				"limit": {
					"type":        "number",
					"description": "This represents the page size"
				},
				"offset": {
					"type":        "number",
					"description": "This represents the page number"
				},
				"count_total": {
					"type":        "number",
					"description": "This represents the option to return total count of data reference keys.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryDataReference",
		config.QueryDataReference,
		`{
			"description": "QueryDataReference is a function to query data reference by a key.",
			"parameters": {
				"key": {
					"type":        "string",
					"description": "This represents data reference key."
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryProposals",
		config.QueryProposals,
		`{
			"description": "QueryProposals is a function to query all proposals.",
			"parameters": {
				"voter": {
					"type":        "string",
					"description": "This represents the kira account address.",
					"optional": true
				},
				"all": {
					"type":        "bool",
					"description": "This an option to query all proposals.",
					"optional": true
				},
				"reverse": {
					"type":        "bool",
					"description": "This an option to sort proposals.",
					"optional": true
				},
				"key": {
					"type":        "string",
					"description": "This is an option to validators pagination. key is a value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set.",
					"optional": true
				},
				"offset": {
					"type":        "string",
					"description": "This is an option to validators pagination. offset is a numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set.",
					"optional": true
				},
				"limit": {
					"type":        "string",
					"description": "This is an option to validators pagination. limit is the total number of results to be returned in the result page. If left empty it will default to a value to be set by each app.",
					"optional": true
				},
				"countTotal": {
					"type":        "string",
					"description": "This is an option to validators pagination. count_total is set to true  to indicate that the result set should include a count of the total number of items available for pagination in UIs. count_total is only respected when offset is used. It is ignored when key is set.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryProposal",
		config.QueryProposal,
		`{
			"description": "QueryProposal is a function to query a proposal by a given proposal_id.",
			"parameters": {
				"proposal_id": {
					"type":        "number",
					"description": "This is an option of a proposal id"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryVoters",
		config.QueryVoters,
		`{
			"description": "QueryVoters is a function to query voters by a given proposal id.",
			"parameters": {
				"proposal_id": {
					"type":        "number",
					"description": "This is an option of a proposal id"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryVotes",
		config.QueryVotes,
		`{
			"description": "QueryVotes is a function to query votes by a given proposal id.",
			"parameters": {
				"proposal_id": {
					"type":        "number",
					"description": "This is an option of a proposal id"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryNetworkProperties",
		config.QueryNetworkProperties,
		`{
			"description": "QueryNetworkProperties is a function to query network properties."
		}`,
	)

	AddInterxFunction(
		"QueryKiraTokensAliases",
		config.QueryKiraTokensAliases,
		`{
			"description": "QueryKiraTokensAliases is a function to query all registered tokens."
		}`,
	)

	AddInterxFunction(
		"QueryKiraTokensRates",
		config.QueryKiraTokensRates,
		`{
			"description": "QueryKiraTokensRates is a function to query all registered token rates."
		}`,
	)

	AddInterxFunction(
		"QueryKiraTokensAliases",
		config.QueryKiraTokensAliases,
		`{
			"description": "QueryKiraTokensAliases is a function to query all tokens aliases."
		}`,
	)

	AddInterxFunction(
		"QueryKiraTokensRates",
		config.QueryKiraTokensRates,
		`{
			"description": "QueryKiraTokensRates is a function to query all tokens rates."
		}`,
	)

	AddInterxFunction(
		"Withdraws",
		config.QueryWithdraws,
		`{
			"description": "Withdraws is a function to query withdraw transactions of the account.",
			"parameters": {
				"account": {
					"type":        "string",
					"description": "This represents the kira account address."
				},
				"type": {
					"type":        "string",
					"description": "This represents the transaction type.",
					"optional": true
				},
				"page": {
					"type":        "int",
					"description": "This represents the page number of results.",
					"optional": true
				},
				"pageSize": {
					"type":        "int",
					"description": "This represents the pageSize number of results. (1 ~ 1000)",
					"optional": true
				},
				"last": {
					"type":        "string",
					"description": "This represents the last transaction hash.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"Deposits",
		config.QueryDeposits,
		`{
			"description": "Deposits is a function to query deposit transactions of the account.",
			"parameters": {
				"account": {
					"type":        "string",
					"description": "This represents the kira account address."
				},
				"type": {
					"type":        "string",
					"description": "This represents the transaction type.",
					"optional": true
				},
				"page": {
					"type":        "int",
					"description": "This represents the page number of results.",
					"optional": true
				},
				"pageSize": {
					"type":        "int",
					"description": "This represents the pageSize number of results. (1 ~ 1000)",
					"optional": true
				},
				"last": {
					"type":        "string",
					"description": "This represents the last transaction hash.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryUnconfirmedTxs",
		config.QueryUnconfirmedTxs,
		`{
			"description": "QueryUnconfirmedTxs is a function to query unconfirmed transactions.",
			"parameters": {
				"limit": {
					"type":        "int",
					"description": "This represents the limit of the transaction. (1 ~ 1000)",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"Broadcast",
		config.PostTransaction,
		`{
			"description": "Broadcast is a function to broadcast signed transaction.",
			"parameters": {
				"tx": {
					"type":        "byte[]",
					"description": "This represents the transaction bytes."
				},
				"mode": {
					"type":        "string",
					"description": "This represents the broadcast mode. (block, sync, async)",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"Faucet",
		config.FaucetRequestURL,
		`{
			"description": "Faucet is a function to claim tokens to the account for free. Returns the available faucet amount when 'claim' and 'token' is unset.",
			"parameters": {
				"claim": {
					"type":        "string",
					"description": "This represents the kira account address.",
					"optional": true
				},
				"token": {
					"type":        "string",
					"description": "This represents the token name.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"Download",
		config.Download,
		`{
			"description": "Download is a function to download a data reference or arbitrary data.",
			"parameters": {
				"module": {
					"type":        "string",
					"description": "This represents the module name. (e.g. DRR for data reference registry.)"
				},
				"key": {
					"type":        "string",
					"description": "This represents the reference key. (It saves reference data with hashed name. e.g. 2CEE6B1689EDDDD6F08EB1EAEC7D3C4E.)"
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryValidators",
		config.QueryValidators,
		`{
			"description": "QueryValidators is a function to query validators.",
			"parameters": {
				"all": {
					"type":        "string",
					"description": "This is an option to query all validators.",
					"optional": true
				},
				"address": {
					"type":        "string",
					"description": "This is an option to query validator by a given kira address",
					"optional": true
				},
				"valkey": {
					"type":        "string",
					"description": "This is an option to query validator by a given valoper address",
					"optional": true
				},
				"pubkey": {
					"type":        "string",
					"description": "This is an option to query validator by a given pubkey",
					"optional": true
				},
				"moniker": {
					"type":        "string",
					"description": "This is an option to query validator by a given moniker",
					"optional": true
				},
				"status": {
					"type":        "string",
					"description": "This is an option to query validators by a given status",
					"optional": true
				},
				"proposer": {
					"type":        "string",
					"description": "This is an option to query validators by a given proposer address",
					"optional": true
				},
				"key": {
					"type":        "string",
					"description": "This is an option to validators pagination. key is a value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set.",
					"optional": true
				},
				"offset": {
					"type":        "string",
					"description": "This is an option to validators pagination. offset is a numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set.",
					"optional": true
				},
				"limit": {
					"type":        "string",
					"description": "This is an option to validators pagination. limit is the total number of results to be returned in the result page. If left empty it will default to a value to be set by each app.",
					"optional": true
				},
				"countTotal": {
					"type":        "string",
					"description": "This is an option to validators pagination. count_total is set to true  to indicate that the result set should include a count of the total number of items available for pagination in UIs. count_total is only respected when offset is used. It is ignored when key is set.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryValidatorInfos",
		config.QueryValidatorInfos,
		`{
			"description": "QueryValidatorInfos is a function to query validator infos.",
			"parameters": {
				"key": {
					"type":        "string",
					"description": "This is an option to validators pagination. key is a value returned in PageResponse.next_key to begin querying the next page most efficiently. Only one of offset or key should be set.",
					"optional": true
				},
				"offset": {
					"type":        "string",
					"description": "This is an option to validators pagination. offset is a numeric offset that can be used when key is unavailable. It is less efficient than using key. Only one of offset or key should be set.",
					"optional": true
				},
				"limit": {
					"type":        "string",
					"description": "This is an option to validators pagination. limit is the total number of results to be returned in the result page. If left empty it will default to a value to be set by each app.",
					"optional": true
				},
				"countTotal": {
					"type":        "string",
					"description": "This is an option to validators pagination. count_total is set to true  to indicate that the result set should include a count of the total number of items available for pagination in UIs. count_total is only respected when offset is used. It is ignored when key is set.",
					"optional": true
				},
				"all": {
					"type":        "boolean",
					"description": "This is an option to validators pagination. all is set to true  to indicate that all the results should be returned.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryBlocks",
		config.QueryBlocks,
		`{
			"description": "QueryBlocks is a function to query blocks with pagination.",
			"parameters": {
				"minHeight": {
					"type":        "string",
					"description": "This is the option of the minimum block height.",
					"optional": true
				},
				"maxHeight": {
					"type":        "string",
					"description": "This is the option of the maximum block height.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryBlockByHeightOrHash",
		config.QueryBlockByHeightOrHash,
		`{
			"description": "QueryBlockByHeightOrHash is a function to query block by height or hash.",
			"parameters": {
				"height": {
					"type":        "string",
					"description": "This is an option of block height or hash.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryBlockTransactions",
		config.QueryBlockTransactions,
		`{
			"description": "QueryBlockTransactions is a function to query block transactions by height.",
			"parameters": {
				"height": {
					"type":        "string",
					"description": "This is an option of block height.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryTransactionResult",
		config.QueryTransactionResult,
		`{
			"description": "QueryTransactionResult is a function to query transaction result by hash.",
			"parameters": {
				"txHash": {
					"type":        "string",
					"description": "This is an option of a transaction hash.",
					"optional": true
				}
			}
		}`,
	)

	AddInterxFunction(
		"QueryGenesis",
		config.QueryGenesis,
		`{
			"description": "QueryGenesis is a function to query genesis."
		}`,
	)

	AddInterxFunction(
		"QueryGenesisSum",
		config.QueryGenesisSum,
		`{
			"description": "QueryGenesisSum is a function to query genesis checksum."
		}`,
	)

	AddInterxFunction(
		"QueryInterxStatus",
		config.QueryStatus,
		`{
			"description": "QueryInterxStatus is a function to query interx informations."
		}`,
	)

	AddInterxFunction(
		"QueryRPCMethods",
		config.QueryRPCMethods,
		`{
			"description": "QueryRPCMethods is a function to query all rpc methods available."
		}`,
	)

	AddInterxFunction(
		"QueryKiraFunctions",
		config.QueryKiraFunctions,
		`{
			"description": "QueryKiraFunctions is a function to query kira functions."
		}`,
	)

	AddInterxFunction(
		"QueryInterxFunctions",
		config.QueryInterxFunctions,
		`{
			"description": "QueryInterxFunctions is a function to query interx functions."
		}`,
	)
}
