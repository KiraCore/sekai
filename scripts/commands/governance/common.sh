#!/bin/bash

sekaid tx customgov councilor claim-seat --from validator --keyring-backend=test --home=$HOME/.sekaid

sekaid tx customgov permission blacklist-permission
sekaid tx customgov permission whitelist-permission

sekaid tx customgov proposal assign-permission
sekaid tx customgov proposal vote

sekaid tx customgov role blacklist-permission
sekaid tx customgov role create
sekaid tx customgov role remove
sekaid tx customgov role remove-blacklist-permission
sekaid tx customgov role remove-whitelist-permission
sekaid tx customgov role whitelist-permission

# querying for voters of a specific proposal
sekaid query customgov voters 1
# querying for votes of a specific proposal
sekaid query customgov votes 1
# querying for a vote of a specific propsal/voter pair
sekaid query customgov vote 1 $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# whitelist permission for claim validator
sekaid keys add lladmin --keyring-backend=test
sekaid tx bank send validator $(sekaid keys show -a lladmin --keyring-backend=test) 1000000ukex --keyring-backend=test --chain-id=testing --fees=200ukex --yes
sekaid tx customgov permission whitelist-permission --from=validator --keyring-backend=test --addr=$(sekaid keys show -a lladmin --keyring-backend=test) --permission=30 --chain-id=testing --fees=200ukex --yes
sekaid tx customgov permission whitelist-permission --from=lladmin --keyring-backend=test --addr=$(sekaid keys show -a lladmin --keyring-backend=test) --permission=2 --chain-id=testing --fees=200ukex --yes
