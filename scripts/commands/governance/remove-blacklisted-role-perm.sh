#!/bin/bash

PermVoteUpsertDataRegistryProposal=11

sekaid tx customgov proposal role remove-blacklisted-permission newrole $PermVoteUpsertDataRegistryProposal  --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

# check permissions
sekaid query customgov role newrole
