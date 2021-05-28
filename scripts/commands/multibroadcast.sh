#!/bin/bash

# block_time modification
# sed -i '' 's/timeout_commit = "5s"/timeout_commit = "20s"/g' $HOME/.sekaid/config/config.toml

sekaid tx customgov proposal set-poor-network-msgs AAA --description="" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=2000000ukex --generate-only > AAA_proposal.json
sekaid tx customgov proposal set-poor-network-msgs BBB --description="" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=3000000ukex --generate-only > BBB_proposal.json
sekaid tx customgov proposal set-poor-network-msgs CCC --description="" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=4000000ukex --generate-only > CCC_proposal.json
sekaid tx customgov proposal set-poor-network-msgs DDD --description="" --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000000ukex --generate-only > DDD_proposal.json

sekaid tx sign AAA_proposal.json --log_level=debug --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_AAA.json
sekaid tx sign BBB_proposal.json --log_level=debug --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_BBB.json
sekaid tx sign CCC_proposal.json --log_level=debug --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_CCC.json
sekaid tx sign DDD_proposal.json --log_level=debug --chain-id=testing --keyring-backend=test --from=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --home=$HOME/.sekaid > signed_DDD.json

sekaid tx broadcast signed_DDD.json --log_level=debug --broadcast-mode=async
sekaid tx broadcast signed_BBB.json --log_level=debug --broadcast-mode=async
sekaid tx broadcast signed_CCC.json --log_level=debug --broadcast-mode=async
sekaid tx broadcast signed_AAA.json --log_level=debug --broadcast-mode=async

sleep 20

sekaid query customgov proposals --log_level=debug