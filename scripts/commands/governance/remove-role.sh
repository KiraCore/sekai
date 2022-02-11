#!/bin/bash

# for complexity of the operation, this proposal is not implemented
sekaid tx customgov proposal role remove newrole --title="title" --description="description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

# check deleted role
sekaid query customgov role 3
sekaid query customgov role newrole
