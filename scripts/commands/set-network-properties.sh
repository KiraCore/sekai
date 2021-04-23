#!/bin/bash

# command with fee set
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --fees=100ukex --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid

# no error response
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set-network-properties\"}]}]}]"

# response when not enough permissions to change tx fee
# "failed to execute message; message index: 0: PermChangeTxFee: not enough permissions"

# command without fee set
# sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid

# response
# "fee out of range [1, 10000]: invalid request"