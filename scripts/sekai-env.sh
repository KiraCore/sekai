#!/usr/bin/env bash

###################################### permissions #################################################

# no-op permission
PermZero=0

# the permission that allows to Set Permissions to other actors
PermSetPermissions=1

# permission that allows to Claim a validator Seat
PermClaimValidator=2

# permission that allows to Claim a Councilor Seat
PermClaimCouncilor=3

# permission to create proposals to whitelist account permission
PermWhitelistAccountPermissionProposal=4

# permission to vote on a proposal to whitelist account permission
PermVoteWhitelistAccountPermissionProposal=5

# permission to change transaction fees - execution fee and fee range
PermChangeTxFee=7

# permission to upsert token rates
PermUpsertTokenInfo=8

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

# permission to create proposals for setting poor network messages
PermCreateSetPoorNetworkMessagesProposal=16

# permission to vote proposals to set poor network messages
PermVoteSetPoorNetworkMessagesProposal=17

# permission to create proposals to upsert token rate
PermCreateUpsertTokenInfoProposal=18

# permission to vote propsals to upsert token rate
PermVoteUpsertTokenInfoProposal=19

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

# permission that allows to Set ClaimValidatorPermission to other actors
PermSetClaimValidatorPermission=30

# permission needed to create a proposal to set proposal duration
PermCreateSetProposalDurationProposal=31

# permission needed to vote a proposal to set proposal duration
PermVoteSetProposalDurationProposal=32

# permission needed to create proposals for blacklisting an account permission.
PermBlacklistAccountPermissionProposal=33

# permission that an actor must have in order to vote a Proposal to blacklist account permission.
PermVoteBlacklistAccountPermissionProposal=34

# permission needed to create proposals for removing whitelisted permission from an account.
PermRemoveWhitelistedAccountPermissionProposal=35

# permission that an actor must have in order to vote a proposal to remove a whitelisted account permission
PermVoteRemoveWhitelistedAccountPermissionProposal=36

# permission needed to create proposals for removing blacklisted permission from an account.
PermRemoveBlacklistedAccountPermissionProposal=37

# permission that an actor must have in order to vote a proposal to remove a blacklisted account permission.
PermVoteRemoveBlacklistedAccountPermissionProposal=38

# permission needed to create proposals for whitelisting an role permission.
PermWhitelistRolePermissionProposal=39

#permission that an actor must have in order to vote a proposal to whitelist role permission.
PermVoteWhitelistRolePermissionProposal=40

#permission needed to create proposals for blacklisting an role permission.
PermBlacklistRolePermissionProposal=41

# permission that an actor must have in order to vote a proposal to blacklist role permission.
PermVoteBlacklistRolePermissionProposal=42

# permission needed to create proposals for removing whitelisted permission from a role.
PermRemoveWhitelistedRolePermissionProposal=43

# permission that an actor must have in order to vote a proposal to remove a whitelisted role permission.
PermVoteRemoveWhitelistedRolePermissionProposal=44;

# permission needed to create proposals for removing blacklisted permission from a role.
PermRemoveBlacklistedRolePermissionProposal=45

# permission that an actor must have in order to vote a proposal to remove a blacklisted role permission.
PermVoteRemoveBlacklistedRolePermissionProposal=46;

# permission needed to create proposals to assign role to an account
PermAssignRoleToAccountProposal=47

# permission that an actor must have in order to vote a proposal to assign role to an account
PermVoteAssignRoleToAccountProposal=48

# permission needed to create proposals to unassign role from an account
PermUnassignRoleFromAccountProposal=49

# permission that an actor must have in order to vote a proposal to unassign role from an account
PermVoteUnassignRoleFromAccountProposal=50

# permission needed to create a proposal to remove a role.
PermRemoveRoleProposal=51

# permission needed to vote a proposal to remove a role.
PermVoteRemoveRoleProposal=52

# permission needed to create proposals to upsert ubi
PermCreateUpsertUBIProposal=53

# permission that an actor must have in order to vote a proposal to upsert ubi
PermVoteUpsertUBIProposal=54

# permission needed to create a proposal to remove ubi.
PermCreateRemoveUBIProposal=55

# permission needed to vote a proposal to remove ubi.
PermVoteRemoveUBIProposal=56

# permission needed to create a proposal to slash validator.
PermCreateSlashValidatorProposal=57

# permission needed to vote a proposal to slash validator.
PermVoteSlashValidatorProposal=58

# permission needed to create a proposal related to basket.
PermCreateBasketProposal=59

# permission needed to vote a proposal related to basket.
PermVoteBasketProposal=60

# permission needed to handle emergency issues on basket.
PermHandleBasketEmergency=61

# permission needed to create a proposal to reset whole councilor rank
PermCreateResetWholeCouncilorRankProposal=62

# permission needed to vote on reset whole councilor rank proposal
PermVoteResetWholeCouncilorRankProposal=63

# permission needed to create a proposal to jail councilors
PermCreateJailCouncilorProposal=64

# permission needed to vote on jail councilors proposal
PermVoteJailCouncilorProposal=65

# permission needed to create a poll proposal
PermCreatePollProposal=66

# permission needed to create a dapp proposal without bond
PermCreateDappProposalWithoutBond=67

# permission needed to create a proposal to set execution fees
PermCreateSetExecutionFeesProposal=68

# permission needed to vote on set execution fees proposal
PermVoteSetExecutionFeesProposal=69

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
MsgTypeUnassignRole="unassign-role"
MsgTypeWhitelistRolePermission="whitelist-role-permission"
MsgTypeBlacklistRolePermission="blacklist-role-permission"
MsgTypeRemoveWhitelistRolePermission="remove-whitelist-role-permission"
MsgTypeRemoveBlacklistRolePermission="remove-blacklist-role-permission"
MsgTypeClaimValidator="claim-validator"
MsgTypeUpsertTokenInfo="upsert-token-rate"

###################################### function IDs ######################################
FuncIDMsgSend=1
FuncIDMultiSend=2

FuncIDMsgSubmitProposal=10
FuncIDMsgVoteProposal=11
FuncIDMsgRegisterIdentityRecords=12
FuncIDMsgDeleteIdentityRecords=13
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
FuncIDMsgUnassignRole=27
FuncIDMsgWhitelistRolePermission=28
FuncIDMsgBlacklistRolePermission=29
FuncIDMsgRemoveWhitelistRolePermission=30
FuncIDMsgRemoveBlacklistRolePermission=31
FuncIDMsgClaimValidator=32
FuncIDMsgUpsertTokenInfo=34
FuncIDMsgActivate=35
FuncIDMsgPause=36
FuncIDMsgUnpause=37