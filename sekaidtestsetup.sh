#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing
sekaid keys add validator --keyring-backend=test
sekaid add-genesis-account $(sekaid keys show validator -a --keyring-backend=test) 1000000000stake,1000000000validatortoken --keyring-backend=test
sekaid gentx-claim validator --keyring-backend=test --moniker="hello"
sekaid start
