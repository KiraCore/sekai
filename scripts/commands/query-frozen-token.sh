#!/bin/bash

# query token blacklists and whitelists
sekaid query tokens token-black-whites
# response
# data:
#   blacklisted:
#   - frozen
#   whitelisted:
#   - ukex

# query network properties
sekaid query customgov network-properties
# response
# properties:
#   enable_foreign_fee_payments: true
#   enable_token_blacklist: false # useful for blacklist use or not
#   enable_token_whitelist: false # useful for whitelist use or not
#   inactive_rank_decrease_percent: "50"
#   jail_max_time: "10"
#   max_tx_fee: "1000000"
#   min_tx_fee: "100"
#   min_validators: "1"
#   mischance_rank_decrease_amount: "10"
#   poor_network_max_bank_send: "1000000"
#   proposal_enactment_time: "300"
#   proposal_end_time: "600"
#   vote_quorum: "33"

# try sending frozen token
sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000frozen --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# response
# token is frozen: invalid request