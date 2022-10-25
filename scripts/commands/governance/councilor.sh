#!/bin/bash

# councilor management commands
sekaid tx customgov councilor claim-seat --moniker="seedvalidator" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov councilor pause --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov councilor unpause --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov councilor activate --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# councilor proposal commands
sekaid tx customgov proposal proposal-reset-whole-councilor-rank --title="Title" --description="Description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov proposal proposal-jail-councilor $(sekaid keys show -a councilor1 --keyring-backend=test --home=$HOME/.sekaid) --title="Title" --description="Description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov proposal vote 1 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# councilor queries
sekaid query customgov councilors
sekaid query customgov non-councilors
sekaid query customgov whitelisted-permission-addresses 1
sekaid query customgov blacklisted-permission-addresses 1
sekaid query customgov whitelisted-role-addresses 1
sekaid query customgov blacklisted-role-addresses 1

# check waiting councilor creation on permission assign
sekaid keys add councilor1 --keyring-backend=test
PermClaimCouncilor=3
sekaid tx customgov permission whitelist --from validator --keyring-backend=test --permission=$PermClaimCouncilor --addr=$(sekaid keys show -a councilor1 --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes --broadcast-mode=block

# check waiting councilor creation on role assign
sekaid keys add councilor2 --keyring-backend=test
RoleSudo=1
sekaid tx customgov role assign $RoleSudo --addr=$(sekaid keys show -a councilor2 --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes --broadcast-mode=block

# try claim councilor
sekaid tx bank send validator $(sekaid keys show -a councilor2 --keyring-backend=test --home=$HOME/.sekaid) 100000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx customgov councilor claim-seat --moniker="councilor2" --from=councilor2 --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# open a proposal
sekaid tx customgov proposal proposal-reset-whole-councilor-rank --title="Title" --description="Description" --from=councilor2 --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# try voting twice and see rank changes
sekaid tx customgov proposal vote 1 1 --from=councilor2 --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov proposal vote 1 1 --from=councilor2 --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# moniker update on councilor
sekaid tx customgov register-identity-records --infos-json='{"moniker":"My Moniker"}' --from=councilor1 --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes --broadcast-mode=block
sekaid tx customgov register-identity-records --infos-json='{"moniker":"My Moniker2"}' --from=councilor1 --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes --broadcast-mode=block