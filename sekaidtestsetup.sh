#!/bin/bash

rm -rf $HOME/.sekaid/

cd $HOME

sekaid init --chain-id=testing testing --home=/Users/jgimeno/.sekaid
sekaid keys add validator
sekaid add-genesis-account $(sekaid keys show validator -a) 1000000000stake,1000000000validatortoken
#sekaid gentx validator --chain-id=testing --offline
#sekaid collect-gentxs

sekaid gentx-claim validator
sekaid start --home=/Users/jgimeno/.sekaid
