#!/bin/bash

# register stake token as 1ukex=100stake
PermUpsertTokenRate=8
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="stake" --rate="0.01" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
sekaid query tokens rate stake
# try to spend stake token as fee
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="valstake" --rate="0.01" --fee_payments=true --chain-id=testing --fees=10000stake --home=$HOME/.sekaid  --yes
# smaller amount of fee in foreign currency
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="valstake" --rate="0.02" --fee_payments=true --chain-id=testing --fees=1000stake --home=$HOME/.sekaid  --yes
# try to spend unregistered token (validatortoken) as fee
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="valstake" --rate="0.03" --fee_payments=true --chain-id=testing --fees=1000validatortoken --home=$HOME/.sekaid  --yes