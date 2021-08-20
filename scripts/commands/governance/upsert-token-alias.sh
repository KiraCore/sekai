#!/bin/bash

# upsert alias by governance
## proposal
sekaid tx tokens proposal-upsert-alias --from=validator --keyring-backend=test --title="upsert alias" --description="upsert alias proposal" --symbol="ETH" --name="Ethereum" --icon="myiconurl" --decimals=6 --denoms="finney" --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
## query
sekaid query proposals
## vote
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 