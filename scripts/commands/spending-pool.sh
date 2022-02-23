#!/bin/bash

sekaid tx spending create-spending-pool --name="ValidatorRewardsPool" --claim-start=$(($(date -u +%s))) --claim-end=0 --token="ukex" --rate=0.1 --vote-quorum="33" --vote-period="60" --vote-enactment="30" --owner-roles="" --owner-accounts=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --beneficiary-roles="1" --beneficiary-accounts="" --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes  --broadcast-mode=block 

sekaid tx spending deposit-spending-pool --name="ValidatorRewardsPool" --amount=1000000ukex --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 

sekaid tx spending register-spending-pool-beneficiary --name="ValidatorRewardsPool" --beneficiary-roles="1" --beneficiary-accounts="" --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 
 
sekaid tx spending claim-spending-pool --name="ValidatorRewardsPool" --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 

sekaid query spending pool-by-name ValidatorRewardsPool --home=$HOME/.sekaid
sekaid query spending pool-names

sekaid query customgov roles $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# proposals
sekaid tx spending proposal-update-spending-pool --title="title" --description="description" --name="ValidatorRewardsPool" --claim-start=$(($(date -u +%s))) --claim-end=0 --token="ukex" --rate=0.5 --vote-quorum="33" --vote-period="60" --vote-enactment="30" --owner-roles="" --owner-accounts=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --beneficiary-roles="1" --beneficiary-accounts="" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx spending proposal-spending-pool-withdraw --title="title" --description="description" --name="ValidatorRewardsPool" --beneficiary-accounts=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --amount=210000ukex --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
sekaid tx spending proposal-spending-pool-distribution --title="title" --description="description" --name="ValidatorRewardsPool" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 2 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes  --broadcast-mode=block 
