#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing --home=$HOME/.sekaid
sekaid keys add validator --keyring-backend=test --home=$HOME/.sekaid
sekaid add-genesis-account $(sekaid keys show validator -a --keyring-backend=test --home=$HOME/.sekaid) 1000000000ukex,1000000000validatortoken,1000000000stake,10000000frozen  --home=$HOME/.sekaid
sekaid gentx-claim validator --keyring-backend=test --moniker="hello" --home=$HOME/.sekaid
sekaid start --home=$HOME/.sekaid

