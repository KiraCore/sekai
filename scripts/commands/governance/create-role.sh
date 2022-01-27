#!/bin/bash

sekaid tx customgov proposal role create newrole "NewRole Description" --title="title" --description="description" --whitelist="1,2" --blacklist="3,4" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

# check newly created role
sekaid query customgov role 3
