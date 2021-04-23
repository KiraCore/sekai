#!/bin/bash

# set PermUpsertTokenRate permission to validator address
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# run upsert rate
sekaid tx tokens upsert-rate --from=validator --keyring-backend=test --denom="mykex" --rate="1.5" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes