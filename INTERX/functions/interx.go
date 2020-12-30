package functions

import (
	"encoding/json"

	sekaitypes "github.com/KiraCore/sekai/types"
	"github.com/iancoleman/strcase"
)

var (
	interxFunctions sekaitypes.FunctionList = make(sekaitypes.FunctionList)
)

// AddInterxFunction is a function to add a function
func AddInterxFunction(functionType string, meta string) {
	metadata := sekaitypes.FunctionMeta{}
	if err := json.Unmarshal([]byte(meta), &metadata); err != nil {
		panic(err)
	}

	interxFunctions[strcase.ToCamel(functionType)] = metadata
}

// RegisterInterxFunctions is a function to register all interx functions
func RegisterInterxFunctions() {
	AddInterxFunction(
		"Faucet",
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
		"Withdraws",
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
		"QueryStatus",
		`{
			"description": "QueryStatus is a function to query the node status"
		}`,
	)

	AddInterxFunction(
		"QueryAccount",
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
		"QueryBalance",
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
		"Download",
		`{
			"description": "Download is a function to download a data reference or arbitrary data.",
			"parameters": {
				"path": {
					"type":        "string",
					"description": "This represents the path to download the reference."
				}
			}
		}`,
	)
}
