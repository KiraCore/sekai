package functionmeta

import (
	"encoding/json"

	"github.com/KiraCore/sekai/types"
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/iancoleman/strcase"
	"github.com/tendermint/go-amino"
)

var (
	functionList types.FunctionList = make(types.FunctionList)
)

// RegisterMsg is a wrapper function to register messages
func RegisterMsg(cdc *codec.LegacyAmino, msg sdk.Msg, name string, copts *amino.ConcreteOptions, meta string) {
	cdc.RegisterConcrete(msg, name, copts)

	AddNewFunction(kiratypes.MsgType(msg), meta)
}

// AddNewFunction is a function to add a function
func AddNewFunction(functionType string, meta string) {
	metadata := types.FunctionMeta{}
	if err := json.Unmarshal([]byte(meta), &metadata); err != nil {
		panic(err)
	}
	metadata.FunctionID = types.MsgFuncIDMapping[functionType]
	if metadata.FunctionID == 0 { // error if not exist
		panic("function id should be available for all the function types: " + functionType)
	}
	functionList[strcase.ToCamel(functionType)] = metadata
}

// GetFunctionList is a function to get functions list
func GetFunctionList() types.FunctionList {
	return functionList
}

// RegisterStdMsgs is a function to register std msgs
func RegisterStdMsgs() {
	registerBankMsgs()
}

func registerBankMsgs() {
	AddNewFunction(
		(bank.MsgSend{}).Type(),
		`{
			"description": "MsgSend represents a message to send coins from one account to another.",
			"parameters": {
				"from_address": {
					"type":        "byte[]",
					"description": "This is the address that will send coins."
				},
				"to_address": {
					"type":        "byte[]",
					"description": "This is the address that will receive coins."
				},
				"amount": {
					"type":        "array<Coin>",
					"description": "This is the amount to be sent.",
					"fields": {
						"denom": {
							"type":        "string",
							"description": "This is the denomination of each coin"
						},
						"amount": {
							"type":        "int",
							"description": "This is the amount of each coin"
						}
					}
				}
			}
		}`,
	)

	AddNewFunction(
		(bank.MsgMultiSend{}).Type(),
		`{
			"description": "MsgMultiSend represents an arbitrary multi-in, multi-out send message.",
			"parameters": {
				"inputs": {
					"type":        "array<Input>",
					"description": "This is the inputs that will send coins.",
					"fields": {
						"address": {
							"type":        "byte[]",
							"description": "This is the address that will send coins."
						},
						"coins": {
							"type":        "array<Coin>",
							"description": "This is the amount of coins the account will send.",
							"fields": {
								"denom": {
									"type":        "string",
									"description": "This is the denomination of each coin"
								},
								"amount": {
									"type":        "int",
									"description": "This is the amount of each coin"
								}
							}
						}
					}
				},
				"outputs": {
					"type":        "array<Output>",
					"description": "This is the inputs that will receive coins.",
					"fields": {
						"address": {
							"type":        "byte[]",
							"description": "This is the address that will receive coins."
						},
						"coins": {
							"type":        "array<Coin>",
							"description": "This is the amount of coins the account will receive.",
							"fields": {
								"denom": {
									"type":        "string",
									"description": "This is the denomination of each coin"
								},
								"amount": {
									"type":        "int",
									"description": "This is the amount of each coin"
								}
							}
						}
					}
				}
			}
		}`,
	)
}
