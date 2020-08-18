#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing
sekaid keys add validator --keyring-backend=test
sekaid add-genesis-account $(sekaid keys show validator -a --keyring-backend=test) 1000000000stake,1000000000validatortoken
sekaid gentx-claim validator --keyring-backend=test
sekaid start
