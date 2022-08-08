#!/bin/bash

sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test)
sekaid query customstaking validator --addr=$(sekaid keys show -a validator --keyring-backend=test)
sekaid query multistaking pools

sekaid tx multistaking upsert-staking-pool kiravaloper1zkka0hucf8mt8wwhqfft7pc0cgv7gk0merzmrs --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx multistaking delegate kiravaloper1zkka0hucf8mt8wwhqfft7pc0cgv7gk0merzmrs 1000000ukex --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block

# proposal to create slash validator
sekaid tx customslashing proposal-slash-validator --offender=kiravaloper1zkka0hucf8mt8wwhqfft7pc0cgv7gk0merzmrs --staking-pool-id=1 --misbehaviour-time=1659927223 --misbehaviour-type="manual-slash" --jail-percentage=10 --colluders="" --refutation="" --title="Slash validator" --description="Slash valiator" --from=validator --chain-id=testing --keyring-backend=test  --fees=100ukex --yes --log_format=json --broadcast-mode=async --output=json --home=$HOME/.sekaid --broadcast-mode=block

# refute slash proposal
sekaid tx customslashing refute-slash-validator-proposal --refutation="refutation.com/1" --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block

# vote slash validator proposal
sekaid tx customgov proposal vote 1 1 --slash=20 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block

# query slash validator proposals
sekaid query customslashing slash-proposals 
