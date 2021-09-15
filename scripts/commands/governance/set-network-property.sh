#!/bin/bash

sekaid tx customgov proposal set-network-property MIN_TX_FEE 101  --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

### possible network property upgrade options ###
# MIN_TX_FEE: minimum transaction fee
# MAX_TX_FEE: maximum transaction fee
# VOTE_QUORUM: vote quorum to reach to move to enactment
# PROPOSAL_END_TIME: the duration to start processing the proposal
# PROPOSAL_ENACTMENT_TIME: the duration to wait for enactment after proposal processing
# MIN_PROPOSAL_END_BLOCKS: minimum blocks required for proposal voting
# MIN_PROPOSAL_ENACTMENT_BLOCKS: min blocks required for proposal enactment
# ENABLE_FOREIGN_FEE_PAYMENTS: flag to enable foreign tokens to be used as transaction fee
# MISCHANCE_RANK_DECREASE_AMOUNT: rank decrease amount per mischance increase (default 10)
# MAX_MISCHANCE: maximum mischance a validator could be in active status, default 110
# MISCHANCE_CONFIDENCE: the number of blocks validator miss to start counting mischance, default 10
# INACTIVE_RANK_DECREASE_PERCENT: percentage of decrease per status movement from active to inactive (default 50%)
# POOR_NETWORK_MAX_BANK_SEND: maximum amount of transfer on poor network, default 10000ukex
# MIN_VALIDATORS: minimum number of validators to perform full network actions - otherwise, it's called poor network
# JAIL_MAX_TIME: maximum jailed status duration in seconds to get back to the validator set again
# ENABLE_TOKEN_WHITELIST: TokenWhitelist is valid when enable_token_whitelist is set
# ENABLE_TOKEN_BLACKLIST: TokenBlacklist is valid when enable_token_blacklist is set
# MIN_IDENTITY_APPROVAL_TIP: minimum amount of tip to be paid to the approver per request
# UNIQUE_IDENTITY_KEYS: manages unique identity keys to be used as an address identifier

# set network property for unique identity keys
sekaid tx customgov proposal set-network-property UNIQUE_IDENTITY_KEYS moniker,email  --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# check all proposals
sekaid query customgov proposals
# check changed network properties
sekaid query customgov network-properties
