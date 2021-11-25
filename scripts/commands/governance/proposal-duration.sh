#!/bin/bash

# TODO: add dynamic duration proposal set examples here


# create proposal for setting poor network msgs
sekaid tx customgov proposal set-proposal-duration-proposal UpsertDataRegistryProposal 300 --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes

# query for proposals
sekaid query customgov proposals

# vote on the proposal
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# check votes
sekaid query customgov votes 1 
# wait until vote end time finish
sekaid query customgov proposals

# query proposal duration
sekaid query customgov proposal-duration UpsertDataRegistryProposal

# query all proposal durations
sekaid query customgov all-proposal-durations

# batch operation
sekaid tx customgov proposal set-batch-proposal-durations-proposal UpsertDataRegistry,SetNetworkProperty 300,300 --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes