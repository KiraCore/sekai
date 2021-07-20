#!/bin/bash

# create proposal for setting poor network msgs
sekaid tx customgov proposal set-poor-network-msgs AAA,BBB --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
# query for proposals
sekaid query customgov proposals
# set permission to vote on proposal
sekaid tx customgov permission whitelist-permission --permission=19 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
# vote on the proposal
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# check votes
sekaid query customgov votes 1 
# wait until vote end time finish
sekaid query customgov proposals
# query poor network messages
sekaid query customgov poor-network-messages

# whitelist permission for modifying network properties
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=7 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# test poor network messages after modifying min_validators section
sekaid tx customgov set-network-properties --from validator --min_validators="2" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# set permission for upsert token rate
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# try running upsert token rate which is not allowed on poor network
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="mykex" --rate="1.5" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
# try sending more than allowed amount via bank send
sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# try setting network property by governance to allow more amount sending
sekaid tx customgov proposal set-network-property POOR_NETWORK_MAX_BANK_SEND 100000000  --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
# try sending after modification of poor network bank send param
sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes