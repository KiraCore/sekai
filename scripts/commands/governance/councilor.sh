#!/bin/bash

# councilor management commands
sekaid tx customgov councilor claim-seat --moniker="seedvalidator" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov councilor pause --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov councilor unpause --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex
sekaid tx customgov councilor activate --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# councilor proposal commands
sekaid tx customgov proposal proposal-reset-whole-councilor-rank --title="Title" --description="Description" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing -y --broadcast-mode=block --fees=100ukex

# councilor queries
sekaid query customgov councilors
sekaid query customgov non-councilors
sekaid query customgov whitelisted-permission-addresses 1
sekaid query customgov blacklisted-permission-addresses 1
sekaid query customgov whitelisted-role-addresses 1
sekaid query customgov blacklisted-role-addresses 1
