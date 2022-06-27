#!/bin/bash

sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test)
sekaid query customstaking validator --addr=$(sekaid keys show -a validator --keyring-backend=test)
sekaid query multistaking pools
sekaid query multistaking undelegations
sekaid query multistaking outstanding-rewards $(sekaid keys show -a validator --keyring-backend=test)
sekaid query multistaking compound-info $(sekaid keys show -a validator --keyring-backend=test)

sekaid tx multistaking upsert-staking-pool kiravaloper1nwcljqs98zkr39pwenngquryaaueztv6ejtljt --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking delegate kiravaloper1nwcljqs98zkr39pwenngquryaaueztv6ejtljt 1000000ukex --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking undelegate kiravaloper1nwcljqs98zkr39pwenngquryaaueztv6ejtljt 10000ukex --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking claim-undelegation 1 --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking claim-rewards --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking set-compound-info true "" --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes --broadcast-mode=block