#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing --home=$HOME/.sekaid
sekaid keys add validator --keyring-backend=test --home=$HOME/.sekaid
sekaid add-genesis-account $(sekaid keys show validator -a --keyring-backend=test --home=$HOME/.sekaid) 1000000000ukex,1000000000validatortoken,1000000000stake,10000000frozen,10000000samolean  --home=$HOME/.sekaid
sekaid gentx-claim validator --keyring-backend=test --moniker="hello" --home=$HOME/.sekaid

cat $HOME/.sekaid/config/genesis.json | jq '.app_state["customgov"]["network_properties"]["proposal_end_time"]="30"' > $HOME/.sekaid/config/tmp_genesis.json && mv $HOME/.sekaid/config/tmp_genesis.json $HOME/.sekaid/config/genesis.json
cat $HOME/.sekaid/config/genesis.json | jq '.app_state["customgov"]["network_properties"]["proposal_enactment_time"]="10"' > $HOME/.sekaid/config/tmp_genesis.json && mv $HOME/.sekaid/config/tmp_genesis.json $HOME/.sekaid/config/genesis.json

sekaid start --home=$HOME/.sekaid
