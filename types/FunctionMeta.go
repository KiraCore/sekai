package types

import (
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidence "github.com/cosmos/cosmos-sdk/x/evidence/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc-transfer/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type FunctionParameter struct {
	Type        string              `json:"type"`
	Optional    bool                `json:"optional"`
	Description string              `json:"description"`
	Fields      *FunctionParameters `json:"fields,omitempty"`
}

type FunctionParameters = map[string]FunctionParameter

type FunctionMeta struct {
	FunctionID  int64              `json:"function_id"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

type FunctionList = map[string]FunctionMeta

var MsgFuncIDMapping = map[string]uint32{
	(bank.MsgSend{}).Type():                                     1,
	(bank.MsgMultiSend{}).Type():                                2,
	(crisis.MsgVerifyInvariant{}).Type():                        3,
	(distribution.MsgSetWithdrawAddress{}).Type():               4,
	(distribution.MsgWithdrawDelegatorReward{}).Type():          5,
	(distribution.MsgWithdrawValidatorCommission{}).Type():      6,
	(distribution.MsgFundCommunityPool{}).Type():                7,
	(evidence.MsgSubmitEvidence{}).Type():                       8,
	(gov.MsgSubmitProposal{}).Type():                            9,
	(gov.MsgDeposit{}).Type():                                   10,
	(gov.MsgVote{}).Type():                                      11,
	(ibc.MsgTransfer{}).Type():                                  12,
	(slashing.MsgUnjail{}).Type():                               13,
	(staking.MsgCreateValidator{}).Type():                       14,
	(staking.MsgEditValidator{}).Type():                         15,
	(staking.MsgDelegate{}).Type():                              16,
	(staking.MsgBeginRedelegate{}).Type():                       17,
	(staking.MsgUndelegate{}).Type():                            18,
	(&customgovtypes.MsgSetNetworkProperties{}).Type():          19,
	(&customgovtypes.MsgSetExecutionFee{}).Type():               20,
	(&customgovtypes.MsgProposalAssignPermission{}).Type():      21,
	(&customgovtypes.MsgProposalSetNetworkProperty{}).Type():    22,
	(&customgovtypes.MsgProposalUpsertDataRegistry{}).Type():    23,
	(&customgovtypes.MsgVoteProposal{}).Type():                  24,
	(&customgovtypes.MsgClaimCouncilor{}).Type():                25,
	(&customgovtypes.MsgWhitelistPermissions{}).Type():          26,
	(&customgovtypes.MsgBlacklistPermissions{}).Type():          27,
	(&customgovtypes.MsgCreateRole{}).Type():                    28,
	(&customgovtypes.MsgAssignRole{}).Type():                    29,
	(&customgovtypes.MsgRemoveRole{}).Type():                    30,
	(&customgovtypes.MsgWhitelistRolePermission{}).Type():       31,
	(&customgovtypes.MsgBlacklistRolePermission{}).Type():       32,
	(&customgovtypes.MsgRemoveWhitelistRolePermission{}).Type(): 33,
	(&customgovtypes.MsgRemoveBlacklistRolePermission{}).Type(): 34,
	(&customstakingtypes.MsgClaimValidator{}).Type():            35,
	(&tokenstypes.MsgUpsertTokenAlias{}).Type():                 36,
	(&tokenstypes.MsgUpsertTokenRate{}).Type():                  37,
}
