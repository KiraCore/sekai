#!/bin/bash

# command to set PermChangeTxFee permission
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid
# good response
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"whitelist-permissions\"}]}]}]"