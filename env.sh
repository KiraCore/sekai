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

# permission needed to create a proposal to reset whole validator rank
PermCreateResetWholeValidatorRankProposal=26

# permission needed to vote on reset whole validator rank proposal
PermVoteResetWholeValidatorRankProposal=27

# permission needed to create a proposal for software upgrade
PermCreateSoftwareUpgradeProposal=28

# permission needed to vote on software upgrade proposal
PermVoteSoftwareUpgradeProposal=29

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

FuncIDMsgSubmitProposal=10
FuncIDMsgVoteProposal=11
FuncIDMsgRegisterIdentityRecords=12
FuncIDMsgEditIdentityRecord=13
FuncIDMsgRequestIdentityRecordsVerify=14
FuncIDMsgHandleIdentityRecordsVerifyRequest=15
FuncIDMsgCancelIdentityRecordsVerifyRequest=16

FuncIDMsgSetNetworkProperties=20
FuncIDMsgSetExecutionFee=21
FuncIDMsgClaimCouncilor=22
FuncIDMsgWhitelistPermissions=23
FuncIDMsgBlacklistPermissions=24
FuncIDMsgCreateRole=25
FuncIDMsgAssignRole=26
FuncIDMsgRemoveRole=27
FuncIDMsgWhitelistRolePermission=28
FuncIDMsgBlacklistRolePermission=29
FuncIDMsgRemoveWhitelistRolePermission=30
FuncIDMsgRemoveBlacklistRolePermission=31
FuncIDMsgClaimValidator=32
FuncIDMsgUpsertTokenAlias=33
FuncIDMsgUpsertTokenRate=34
FuncIDMsgActivate=35
FuncIDMsgPause=36
FuncIDMsgUnpause=37