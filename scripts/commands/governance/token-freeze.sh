#!/bin/bash

# create a proposal to blacklist validatortoken
sekaid tx tokens proposal-update-tokens-blackwhite --is_blacklist=true --is_add=true --tokens=validatortoken --tokens=kava --from=validator --title="title" --description="description" --chain-id=testing --keyring-backend=test --fees=100ukex --home=$HOME/.sekaid --yes
# check proposal ID
sekaid query customgov proposals
# whitelist permission to vote on proposal
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermVoteTokensWhiteBlackChangeProposal --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# vote on proposal
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# get all votes
sekaid query customgov vote 1 $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)