#!/bin/bash

sekaid query customgov network-properties
sekaid query ubi ubi-record-by-name ValidatorQuaterlyIncome --home=$HOME/.sekaid
sekaid query ubi ubi-records

sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# set permissions
PermCreateUpsertUBIProposal=53
PermVoteUpsertUBIProposal=54
PermCreateRemoveUBIProposal=55
PermVoteRemoveUBIProposal=56
sekaid tx customgov proposal account whitelist-permission $PermCreateUpsertUBIProposal  --title="title" --description="description" --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx customgov proposal account whitelist-permission $PermVoteUpsertUBIProposal  --title="title" --description="description" --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx customgov proposal account whitelist-permission $PermCreateRemoveUBIProposal  --title="title" --description="description" --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx customgov proposal account whitelist-permission $PermVoteRemoveUBIProposal  --title="title" --description="description" --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 

# proposals
sekaid tx customgov proposal set-network-property UBI_HARDCAP 20000000  --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx ubi proposal-upsert-ubi --title="title" --description="description" --name="ValidatorQuaterlyIncome" --distr-start=0 --distr-end=0 --amount=2000000 --period=7776000 --pool-name="ValidatorBasicRewardsPool" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx ubi proposal-remove-ubi --title="title" --description="description" --name="ValidatorQuaterlyIncome" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes  --broadcast-mode=block 
sekaid tx customgov proposal vote 2 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes  --broadcast-mode=block 
sekaid tx customgov proposal vote 3 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes  --broadcast-mode=block 
