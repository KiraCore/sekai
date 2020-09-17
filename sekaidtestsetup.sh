#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing --home=/Users/jgimeno/.sekaid
sekaid keys add validator --keyring-backend=test --home=/Users/jgimeno/.sekaid
sekaid add-genesis-account $(sekaid keys show validator -a --home=/Users/jgimeno/.sekaid --keyring-backend=test) 1000000000stake,1000000000validatortoken  --home=/Users/jgimeno/.sekaid
sekaid gentx-claim validator --keyring-backend=test --moniker="hello" --home=/Users/jgimeno/.sekaid
sekaid start --home=/Users/jgimeno/.sekaid
