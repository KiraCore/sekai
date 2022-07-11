#!/bin/bash

sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test)
sekaid query customstaking validator --addr=$(sekaid keys show -a validator --keyring-backend=test)
sekaid query multistaking pools
sekaid query multistaking undelegations
sekaid query multistaking outstanding-rewards $(sekaid keys show -a validator --keyring-backend=test)
sekaid query multistaking compound-info $(sekaid keys show -a validator --keyring-backend=test)

sekaid tx multistaking upsert-staking-pool kiravaloper1h3f7w7ekjnpfcyjhktq06n4rl8yued9a0ap468 --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking delegate kiravaloper1h3f7w7ekjnpfcyjhktq06n4rl8yued9a0ap468 1000000ukex --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking undelegate kiravaloper1h3f7w7ekjnpfcyjhktq06n4rl8yued9a0ap468 10000ukex --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking claim-undelegation 1 --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking claim-rewards --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking set-compound-info true "" --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes --broadcast-mode=block

sekaid keys add delegator1 --keyring-backend=test --home=$HOME/.sekaid
sekaid keys add delegator2 --keyring-backend=test --home=$HOME/.sekaid
sekaid tx bank send validator $(sekaid keys show -a delegator1 --keyring-backend=test --home=$HOME/.sekaid) 100ubtc,10000ukex --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx bank send validator $(sekaid keys show -a delegator2 --keyring-backend=test --home=$HOME/.sekaid) 1000000ukex --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking delegate kiravaloper1h3f7w7ekjnpfcyjhktq06n4rl8yued9a0ap468 10ubtc --from=delegator1 --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking delegate kiravaloper1h3f7w7ekjnpfcyjhktq06n4rl8yued9a0ap468 100ukex --from=delegator2 --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block

sekaid tx multistaking set-compound-info true "" --from=delegator1 --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes --broadcast-mode=block
sekaid tx bank send validator $(sekaid keys show -a delegator1 --keyring-backend=test --home=$HOME/.sekaid) 1000000v1/ukex --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking register-delegator --from=delegator1 --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block

sekaid query bank balances $(sekaid keys show -a delegator1 --keyring-backend=test)
sekaid query bank balances $(sekaid keys show -a delegator2 --keyring-backend=test)
sekaid query multistaking outstanding-rewards $(sekaid keys show -a delegator1 --keyring-backend=test)
sekaid query multistaking outstanding-rewards $(sekaid keys show -a delegator2 --keyring-backend=test)
sekaid query multistaking staking-pool-delegators kiravaloper1h3f7w7ekjnpfcyjhktq06n4rl8yued9a0ap468