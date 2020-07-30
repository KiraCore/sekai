#!/bin/bash

rm -rf $HOME/.sekaid/
rm -rf $HOME/.sekaicli/

cd $HOME

sekaid init --chain-id=testing testing
sekaid keys add validator
sekaid add-genesis-account $(sekaid keys show validator -a) 1000000000stake,1000000000validatortoken
#sekaid gentx validator --chain-id=testing --offline
#sekaid collect-gentxs
#sekaid start
