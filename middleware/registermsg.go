package middleware

import (
	"github.com/KiraCore/sekai/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	amino "github.com/tendermint/go-amino"
)

var (
	functionList types.FunctionList = make(types.FunctionList)
)

// RegisterMsg is a wrapper function to register messages
func RegisterMsg(cdc *codec.LegacyAmino, msg sdk.Msg, name string, copts *amino.ConcreteOptions, meta types.FunctionMeta) {
	cdc.RegisterConcrete(msg, name, copts)

	AddNewFunction(msg.Type(), meta)
}

// AddNewFunction is a function to add a function
func AddNewFunction(functionType string, meta types.FunctionMeta) {
	functionList[functionType] = meta
}

// GetFunctionList is a function to get functions list
func GetFunctionList() types.FunctionList {
	return functionList
}

// RegisterStdMsgs is a function to register std msgs
func RegisterStdMsgs() {
	AddNewFunction((bank.MsgSend{}).Type(), types.NewFunctionMeta(
		0,
		"MsgSend represents a message to send coins from one account to another.",
		types.FunctionParameters{
			"from_address": types.FunctionParameter{
				Type:        "byte[]",
				Description: "This is the address that will send coins.",
			},
			"to_address": types.FunctionParameter{
				Type:        "byte[]",
				Description: "This is the address that will receive coins.",
			},
			"amount": types.FunctionParameter{
				Type:        "array<Coin>",
				Description: "This is the amount to be sent.",
				Fields: &types.FunctionParameters{
					"denom": types.FunctionParameter{
						Type:        "string",
						Description: "This is the denomination of each coin",
					},
					"amount": types.FunctionParameter{
						Type:        "int",
						Description: "This is the amount of each coin",
					},
				},
			},
		},
	))

	AddNewFunction((bank.MsgMultiSend{}).Type(), types.NewFunctionMeta(
		1,
		"MsgMultiSend represents an arbitrary multi-in, multi-out send message.",
		types.FunctionParameters{
			"inputs": types.FunctionParameter{
				Type:        "array<Input>",
				Description: "This is the inputs that will send coins.",
				Fields: &types.FunctionParameters{
					"address": types.FunctionParameter{
						Type:        "byte[]",
						Description: "This is the address that will send coins.",
					},
					"coins": types.FunctionParameter{
						Type:        "array<Coin>",
						Description: "This is the amount of coins the account will send.",
						Fields: &types.FunctionParameters{
							"denom": types.FunctionParameter{
								Type:        "string",
								Description: "This is the denomination of each coin",
							},
							"amount": types.FunctionParameter{
								Type:        "int",
								Description: "This is the amount of each coin",
							},
						},
					},
				},
			},
			"outputs": types.FunctionParameter{
				Type:        "array<Output>",
				Description: "This is the inputs that will receive coins.",
				Fields: &types.FunctionParameters{
					"address": types.FunctionParameter{
						Type:        "byte[]",
						Description: "This is the address that will receive coins.",
					},
					"coins": types.FunctionParameter{
						Type:        "array<Coin>",
						Description: "This is the amount of coins the account will receive.",
						Fields: &types.FunctionParameters{
							"denom": types.FunctionParameter{
								Type:        "string",
								Description: "This is the denomination of each coin",
							},
							"amount": types.FunctionParameter{
								Type:        "int",
								Description: "This is the amount of each coin",
							},
						},
					},
				},
			},
		},
	))
}
