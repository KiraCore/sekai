package functionmeta

import (
	"encoding/json"

	"github.com/KiraCore/sekai/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/iancoleman/strcase"
	"github.com/tendermint/go-amino"
)

var (
	functionList types.FunctionList = make(types.FunctionList)
)

// RegisterMsg is a wrapper function to register messages
func RegisterMsg(cdc *codec.LegacyAmino, msg sdk.Msg, name string, copts *amino.ConcreteOptions, meta string) {
	cdc.RegisterConcrete(msg, name, copts)

	AddNewFunction(msg.Type(), meta)
}

// AddNewFunction is a function to add a function
func AddNewFunction(functionType string, meta string) {
	metadata := types.FunctionMeta{}
	if err := json.Unmarshal([]byte(meta), &metadata); err != nil {
		panic(err)
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
	registerCrisisMsgs()
	registerDistributionMsgs()
	registerEvidenceMsgs()
	registerGovMsgs()
	registerStakingMsgs()
}

func registerBankMsgs() {
	AddNewFunction(
		(bank.MsgSend{}).Type(),
		`{
			"function_id": 0,
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
			"function_id": 1,
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

func registerCrisisMsgs() {
	AddNewFunction(
		(crisis.MsgVerifyInvariant{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgVerifyInvariant represents a message to verify a particular invariance.",
			"parameters": {
				"sender": {
					"type":        "byte[]",
					"description": "Sender address"
				},
				"invariant_module_name": {
					"type":        "string",
					"description": "Invariant module name"
				},
				"invariant_route": {
					"type":        "string",
					"description": "Invariant route"
				}
			}
		}`,
	)
}

func registerDistributionMsgs() {
	AddNewFunction(
		(distribution.MsgSetWithdrawAddress{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgSetWithdrawAddress sets the withdraw address for a delegator (or validator self-delegation).",
			"parameters": {
				"delegator_address": {
					"type":        "byte[]",
					"description": "Delegator address"
				},
				"withdraw_address": {
					"type":        "byte[]",
					"description": "Withdraw address"
				}
			}
		}`,
	)

	AddNewFunction(
		(distribution.MsgWithdrawDelegatorReward{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator from a single validator.",
			"parameters": {
				"delegator_address": {
					"type":        "byte[]",
					"description": "Delegator address"
				},
				"withdraw_address": {
					"type":        "byte[]",
					"description": "Withdraw address"
				}
			}
		}`,
	)

	AddNewFunction(
		(distribution.MsgWithdrawValidatorCommission{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgWithdrawValidatorCommission withdraws the full commission to the validator address.",
			"parameters": {
				"validator_address": {
					"type":        "byte[]",
					"description": "Validator address"
				}
			}
		}`,
	)

	AddNewFunction(
		(distribution.MsgFundCommunityPool{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgFundCommunityPool allows an account to directly fund the community pool.",
			"parameters": {
				"amount": {
					"type":        "Coin",
					"description": "Fund amount",
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
				},
				"depositor": {
					"type":        "byte[]",
					"description": "Depositor address"
				}
			}
		}`,
	)
}

func registerEvidenceMsgs() {
	AddNewFunction(
		(evidence.MsgSubmitEvidence{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgSubmitEvidence represents a message that supports submitting arbitrary Evidence of misbehavior such as equivocation or counterfactual signing.",
			"parameters": {
				"submitter": {
					"type":        "byte[]",
					"description": "Submitter address"
				},
				"evidence": {
					"type":        "any",
					"description": "Evidence"
				}
			}
		}`,
	)
}

func registerGovMsgs() {
	AddNewFunction(
		(gov.MsgSubmitProposal{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgSubmitProposal defines an sdk.Msg type that supports submitting arbitrary proposal Content.",
			"parameters": {
				"content": {
					"type":        "any",
					"description": "Content"
				},
				"initial_deposit": {
					"type":        "array<Coin>",
					"description": "Initial deposit amount",
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
				},
				"proposer": {
					"type":        "byte[]",
					"description": "Proposer address"
				}
			}
		}`,
	)

	AddNewFunction(
		(gov.MsgDeposit{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgDeposit defines a message to submit a deposit to an existing proposal.",
			"parameters": {
				"proposal_id": {
					"type":        "Uint64",
					"description": "Proposal ID"
				},
				"depositor": {
					"type":        "byte[]",
					"description": "Depositor address"
				},
				"amount": {
					"type":        "array<Coin>",
					"description": "Deposit amount",
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
		(gov.MsgVote{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgVote defines a message to cast a vote.",
			"parameters": {
				"proposal_id": {
					"type":        "Uint64",
					"description": "Proposal ID"
				},
				"voter": {
					"type":        "byte[]",
					"description": "Voter address"
				},
				"option": {
					"type":        "Int32",
					"description": "Vote option"
				}
			}
		}`,
	)
}

func registerStakingMsgs() {
	AddNewFunction(
		(staking.MsgCreateValidator{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgCreateValidator defines an SDK message for creating a new validator.",
			"parameters": {
				"description": {
					"type":        "Description",
					"description": "Description defines a validator description.",
					"fields": {
						"moniker": {
							"type":        "string",
							"description": "Moniker"
						},
						"identity": {
							"type":        "string",
							"description": "Identity"
						},
						"website": {
							"type":        "string",
							"description": "Website"
						},
						"security_contact": {
							"type":        "string",
							"description": "Security Contact"
						},
						"details": {
							"type":        "string",
							"description": "Details"
						}
					}
				},
				"commission": {
					"type":        "CommissionRates",
					"description": "CommissionRates defines the initial commission rates to be used for creating a validator",
					"fields": {
						"rate": {
							"type":        "BigInt",
							"description": "Normal rate"
						},
						"max_rate": {
							"type":        "BigInt",
							"description": "Maximum rate"
						},
						"max_change_rate": {
							"type":        "BigInt",
							"description": "Maximum change rate"
						}
					}
				},
				"min_self_delegation": {
					"type":        "Int",
					"description": "Minimum self delegation"
				},
				"delegator_address": {
					"type":        "byte[]",
					"description": "Delegator address"
				},
				"validator_address": {
					"type":        "byte[]",
					"description": "Validator address"
				},
				"pubkey": {
					"type":        "string",
					"description": "Public key"
				},
				"value": {
					"type":        "Coin",
					"description": "Self delegation amount",
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
		(staking.MsgEditValidator{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgEditValidator defines an SDK message for editing an existing validator.",
			"parameters": {
				"description": {
					"type":        "Description",
					"description": "Description defines a validator description.",
					"fields": {
						"moniker": {
							"type":        "string",
							"description": "Moniker"
						},
						"identity": {
							"type":        "string",
							"description": "Identity"
						},
						"website": {
							"type":        "string",
							"description": "Website"
						},
						"security_contact": {
							"type":        "string",
							"description": "Security Contact"
						},
						"details": {
							"type":        "string",
							"description": "Details"
						}
					}
				},
				"validator_address": {
					"type":        "byte[]",
					"description": "Validator address"
				},
				"commission_rate": {
					"type":        "BigInt",
					"description": "Commission rate"
				},
				"min_self_delegation": {
					"type":        "Int",
					"description": "Minimum self delegation"
				}
			}
		}`,
	)

	AddNewFunction(
		(staking.MsgDelegate{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgDelegate defines an SDK message for performing a delegation of coins from a delegator to a validator.",
			"parameters": {
				"delegator_address": {
					"type":        "byte[]",
					"description": "Delegator address"
				},
				"validator_address": {
					"type":        "byte[]",
					"description": "Validator address"
				},
				"amount": {
					"type":        "Coin",
					"description": "Delegation amount",
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
		(staking.MsgBeginRedelegate{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgBeginRedelegate defines an SDK message for performing a redelegation of coins from a delegator and source validator to a destination validator.",
			"parameters": {
				"delegator_address": {
					"type":        "byte[]",
					"description": "Delegator address"
				},
				"validator_src_address": {
					"type":        "byte[]",
					"description": "Validator address from"
				},
				"validator_dst_address": {
					"type":        "byte[]",
					"description": "Validator address to"
				},
				"amount": {
					"type":        "Coin",
					"description": "Redelegation amount",
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
		(staking.MsgUndelegate{}).Type(),
		`{
			"function_id": 0,
			"description": "MsgUndelegate defines an SDK message for performing an undelegation from a delegate and a validator.",
			"parameters": {
				"delegator_address": {
					"type":        "byte[]",
					"description": "Delegator address"
				},
				"validator_address": {
					"type":        "byte[]",
					"description": "Validator address"
				},
				"amount": {
					"type":        "Coin",
					"description": "Undelegation amount",
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
}
