package functions

import (
	"encoding/json"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/iancoleman/strcase"
)

type InterxFunctionParameter struct {
	Type        string                    `json:"type"`
	Optional    bool                      `json:"optional"`
	Description string                    `json:"description"`
	Fields      *InterxFunctionParameters `json:"fields,omitempty"`
}

type InterxFunctionParameters = map[string]InterxFunctionParameter

type InterxFunctionMeta struct {
	Endpoint    string                   `json:"endpoint"`
	Description string                   `json:"description"`
	Parameters  InterxFunctionParameters `json:"parameters"`
}

type InterxFunctionList = map[string]InterxFunctionMeta

var (
	interxFunctions InterxFunctionList = make(InterxFunctionList)
)

// AddInterxFunction is a function to add a function
func AddInterxFunction(functionType string, endpoint string, meta string) {
	metadata := InterxFunctionMeta{}
	if err := json.Unmarshal([]byte(meta), &metadata); err != nil {
		panic(err)
	}

	metadata.Endpoint = endpoint

	interxFunctions[strcase.ToCamel(functionType)] = metadata
}

// RegisterInterxFunctions is a function to register all interx functions
func RegisterInterxFunctions() {
	AddInterxFunction(
		"QueryNodeStatus",
		config.QueryKiraStatus,
		`{
			"description": "QueryNodeStatus is a function to query the node status"
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
			}
		}`,
	)

	AddInterxFunction(
		"QueryTotalSupply",
		config.QueryTotalSupply,
		`{
			"description": "QueryTotalSupply is a function to query total supply."
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
			"description": "QueryProposals is a function to query all proposals."
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
				"max": {
					"type":        "int",
					"description": "This represents the maximum number of results. (1 ~ 1000)",
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
				"max": {
					"type":        "int",
					"description": "This represents the maximum number of results. (1 ~ 1000)",
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
