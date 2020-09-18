#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing --home=$HOME/.sekaid
sekaid keys add validator --keyring-backend=test --home=$HOME/.sekaid
sekaid add-genesis-account $(sekaid keys show validator -a --home=$HOME/.sekaid --keyring-backend=test) 1000000000stake,1000000000validatortoken  --home=$HOME/.sekaid
sekaid gentx-claim validator --keyring-backend=test --moniker="hello" --home=$HOME/.sekaid
sekaid start --home=$HOME/.sekaid
