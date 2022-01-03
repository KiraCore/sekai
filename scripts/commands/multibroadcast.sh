#!/bin/bash

# block_time modification
# sed -i '' 's/timeout_commit = "5s"/timeout_commit = "20s"/g' $HOME/.sekaid/config/config.toml

sekaid tx customgov proposal set-poor-network-msgs AAA --title="title" --description="description" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=2000000ukex --generate-only > AAA_proposal.json
sekaid tx customgov proposal set-poor-network-msgs BBB --title="title" --description="description" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=3000000ukex --generate-only > BBB_proposal.json
sekaid tx customgov proposal set-poor-network-msgs CCC --title="title" --description="description" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=4000000ukex --generate-only > CCC_proposal.json
sekaid tx customgov proposal set-poor-network-msgs DDD --title="title" --description="description" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000000ukex --generate-only > DDD_proposal.json

sekaid tx sign AAA_proposal.json --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_AAA.json
sekaid tx sign BBB_proposal.json --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_BBB.json
sekaid tx sign CCC_proposal.json --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_CCC.json
sekaid tx sign DDD_proposal.json --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_DDD.json

sekaid tx broadcast signed_DDD.json --broadcast-mode=async
sekaid tx broadcast signed_BBB.json --broadcast-mode=async
sekaid tx broadcast signed_CCC.json --broadcast-mode=async
sekaid tx broadcast signed_AAA.json --broadcast-mode=async

sleep 20

sekaid query customgov proposals