#!/bin/bash

# set PermUpsertTokenAlias permission to validator address
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenAlias --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# run upsert alias
sekaid tx tokens upsert-alias --from validator --keyring-backend=test --expiration=0 --enactment=0 --allowed_vote_types=0,1 --symbol="ETH" --name="Ethereum" --icon="myiconurl" --decimals=6 --denoms="finney" --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes