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
			"function_id": 1,
			"description": "Faucet is a function to claim tokens to the account for free.",
			"parameters": {
				"claim": {
					"type":        "string",
					"description": "This represents the kira account address."
				},
				"token": {
					"type":        "string",
					"description": "This represents the token name."
				}
			}
		}`,
	)

	AddInterxFunction(
		"Withdraws",
		`{
			"function_id": 2,
			"description": "Withdraws is a function to query withdraw transactions of the account",
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
			"function_id": 3,
			"description": "Deposits is a function to query deposit transactions of the account",
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
}
