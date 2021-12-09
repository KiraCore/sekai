#!/bin/bash

sekaid tx customgov councilor claim-seat --from validator --keyring-backend=test --home=$HOME/.sekaid

sekaid tx customgov permission blacklist-permission
sekaid tx customgov permission remove-blacklisted-permission
sekaid tx customgov permission whitelist-permission
sekaid tx customgov permission remove-whitelisted-permission

# add / remove / query whitelisted permissions
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=7 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid query customgov permissions $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)
sekaid tx customgov permission remove-whitelisted-permission --from validator --keyring-backend=test --permission=7 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid query customgov permissions $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

sekaid tx customgov proposal assign-permission
sekaid tx customgov proposal vote

# role creation, role permission add / remove
sekaid tx customgov role create testRole "testRole Description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
sekaid tx customgov role whitelist-role-permission testRole 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
sekaid tx customgov role blacklist-role-permission testRole 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
sekaid tx customgov role remove-whitelisted-role-permission testRole 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
sekaid tx customgov role remove-blacklisted-role-permission testRole 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes

# query all roles
sekaid query customgov all-roles
# query roles for an address
sekaid query customgov roles $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# query a single role
sekaid query customgov role sudo
sekaid query customgov role 1

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
