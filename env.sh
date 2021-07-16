#!/bin/bash

###################################### permissions #################################################

# no-op permission
PermZero=0

# the permission that allows to Set Permissions to other actors
PermSetPermissions=1

# permission that allows to Claim a validator Seat
PermClaimValidator=2

# permission that allows to Claim a Councilor Seat
PermClaimCouncilor=3

# permission to create proposals for setting permissions
PermCreateSetPermissionsProposal=4

# permission to vote on a proposal to set permissions
PermVoteSetPermissionProposal=5

# permission to upsert token alias
PermUpsertTokenAlias=6

# permission to change transaction fees - execution fee and fee range
PermChangeTxFee=7

# permission to upsert token rates
PermUpsertTokenRate=8

# permission to add, modify and assign roles
PermUpsertRole=9

# permission to create a proposal to change the Data Registry
PermCreateUpsertDataRegistryProposal=10

# permission to vote on a proposal to change the Data Registry
PermVoteUpsertDataRegistryProposal=11

# permission to create proposals for setting network property
PermCreateSetNetworkPropertyProposal=12

# permission to vote a proposal to set network property
PermVoteSetNetworkPropertyProposal=13

# permission to create proposals to upsert token alias
PermCreateUpsertTokenAliasProposal=14

# permission to vote proposals to upsert token alias
PermVoteUpsertTokenAliasProposal=15

# permission to create proposals for setting poor network messages
PermCreateSetPoorNetworkMessagesProposal=16

# permission to vote proposals to set poor network messages
PermVoteSetPoorNetworkMessagesProposal=17

# permission to create proposals to upsert token rate
PermCreateUpsertTokenRateProposal=18

# permission to vote propsals to upsert token rate
PermVoteUpsertTokenRateProposal=19

# permission to create a proposal to unjail a validator
PermCreateUnjailValidatorProposal=20

# permission to vote a proposal to unjail a validator
PermVoteUnjailValidatorProposal=21

# permission to create a proposal to create a role
PermCreateRoleProposal=22

# permission to vote a proposal to create a role
PermVoteCreateRoleProposal=23

# permission to create a proposal to change blacklist/whitelisted tokens
PermCreateTokensWhiteBlackChangeProposal=24

# permission to vote a proposal to change blacklist/whitelisted tokens
PermVoteTokensWhiteBlackChangeProposal=25

###################################### transaction_types ######################################
TypeMsgSend="send"
TypeMsgMultiSend="multisend"
MsgTypeVoteProposal="vote-proposal"
MsgTypeWhitelistPermissions="whitelist-permissions"
MsgTypeBlacklistPermissions="blacklist-permissions"
MsgTypeClaimCouncilor="claim-councilor"
MsgTypeSetNetworkProperties="set-network-properties"
MsgTypeSetExecutionFee="set-execution-fee"
MsgTypeCreateRole="create-role"
MsgTypeAssignRole="assign-role"
MsgTypeRemoveRole="remove-role"
MsgTypeWhitelistRolePermission="whitelist-role-permission"
MsgTypeBlacklistRolePermission="blacklist-role-permission"
MsgTypeRemoveWhitelistRolePermission="remove-whitelist-role-permission"
MsgTypeRemoveBlacklistRolePermission="remove-blacklist-role-permission"
MsgTypeClaimValidator="claim-validator"
MsgTypeUpsertTokenAlias="upsert-token-alias"
MsgTypeUpsertTokenRate="upsert-token-rate"

###################################### function IDs ######################################
FuncIDMsgSend=1
FuncIDMultiSend=2
FuncIDMsgVoteProposal=6
FuncIDMsgWhitelistPermissions=7
FuncIDMsgBlacklistPermissions=8
FuncIDMsgClaimCouncilor=9
FuncIDMsgSetNetworkProperties=10
FuncIDMsgSetExecutionFee=11
FuncIDMsgCreateRole=12
FuncIDMsgAssignRole=13
FuncIDMsgRemoveRole=14
FuncIDMsgWhitelistRolePermission=15
FuncIDMsgBlacklistRolePermission=16
FuncIDMsgRemoveWhitelistRolePermission=17
FuncIDMsgRemoveBlacklistRolePermission=18
FuncIDMsgClaimValidator=19
FuncIDMsgUpsertTokenAlias=20
FuncIDMsgUpsertTokenRate=21
