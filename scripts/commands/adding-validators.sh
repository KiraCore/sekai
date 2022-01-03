#!/bin/bash

# sekaid keys add val2 --keyring-backend=test --home=$HOME/.sekaid
# sekaid tx bank send validator $(sekaid keys show -a val2 --keyring-backend=test --home=$HOME/.sekaid) 100000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermCreateSetPermissionsProposal --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermVoteSetPermissionProposal --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

sekaid tx customgov proposal assign-permission $PermClaimValidator  --title="title" --description="description" --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

sekaid tx claim-validator-seat --from validator --keyring-backend=test --home=$HOME/.sekaid --validator-key=kiravaloper1ntk7n5y38en5dvnhvmruwagmkemq76x8s4pnwu --moniker="validator" --chain-id=testing --fees=100ukex --yes

# get ValAddress (kiravaloperxxx) from validator key
sekaid val-address $(sekaid keys show -a validator --keyring-backend=test)

# sekaid tx claim-validator-seat --from val2 --keyring-backend=test --home=$HOME/.sekaid --pubkey=kiravalconspub1zcjduepqdllep3v5wv04hmu987rv46ax7fml65j3dh5tf237ayn5p59jyamq04048n --validator-key=kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w --moniker="val2" --chain-id=testing --fees=100ukex --yes
# sekaid tx claim-validator-seat --from val2 --keyring-backend=test --home=$HOME/.sekaid --validator-key=kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w --moniker="val2" --chain-id=testing --fees=100ukex --yes