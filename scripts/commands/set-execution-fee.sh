#!/bin/bash

# command
sekaid tx customgov set-execution-fee --from validator --execution_name="B" --transaction_type="B" --execution_fee=10 --failure_fee=1 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=10ukex --home=$HOME/.sekaid

# response
# "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set-execution-fee\"}]}]}]"

# command for setting execution fee
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="set-network-properties" --execution_fee=10000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid

# Here, the value should be looked at is `--execution_name="set-network-properties"`, `--execution_fee=10000` and `--failure_fee=1000`.

# check execution fee validation
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid
# response
# "fee is less than failure fee 1000: invalid request"

# Here, the value should be looked at is `"fee is less than failure fee 1000: invalid request"`.
# In this case, issue is found on ante step and fee is not being paid at all.

# preparation for networks (v1) failure=1000, execution=10000
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="set-network-properties" --execution_fee=10000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# preparation for networks (v2) failure=1000, execution=500
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="set-network-properties" --execution_fee=500 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# init user1 with 100000ukex
sekaid keys add user1 --keyring-backend=test --home=$HOME/.sekaid
sekaid tx bank send validator $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) 100000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) --yes

# try changing set-network-properties with user1 that does not have ChangeTxFee permission
sekaid tx customgov set-network-properties --from user1 --min_tx_fee="2" --max_tx_fee="25000" --keyring-backend=test --chain-id=testing --fees=1000ukex --home=$HOME/.sekaid --yes
# this should fail and balance should be (previousBalance - failureFee)
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid)

# whitelist user1's permission for ChangeTxFee and try again
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov set-network-properties --from user1 --min_tx_fee="2" --max_tx_fee="25000" --keyring-backend=test --chain-id=testing --fees=1000ukex --home=$HOME/.sekaid --yes
# this should fail and balance should be (previousBalance - successFee)
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid)