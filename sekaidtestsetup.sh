#!/bin/bash

rm -rf $HOME/.sekaid/
rm -rf $HOME/.sekaicli/

cd $HOME

sekaid init --chain-id=testing testing
sekaicli keys add validator
sekaid add-genesis-account $(sekaicli keys show validator -a) 1000000000stake,1000000000validatortoken
sekaid gentx --name validator
sekaid collect-gentxs
sekaid start
