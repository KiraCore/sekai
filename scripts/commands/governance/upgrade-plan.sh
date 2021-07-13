#!/bin/bash

sekaid tx upgrade set-plan --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=1 --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo=1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --log_level=debug --yes

# {"height":"175","txhash":"F942BD2DC74DB31334477FD4BCF8BED6A9CA173A691F9A7D1A0AB2885C72DD47","codespace":"","code":0,"data":"0A1E0A1870726F706F73652D736F6674776172652D7570677261646512020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"propose-software-upgrade\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"propose-software-upgrade"}]}]}],"info":"","gas_wanted":"0","gas_used":"15650","tx":null,"timestamp":""}


sekaid tx customgov permission whitelist-permission --permission=29 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

sekaid query upgrade current-plan --log_level=debug

sekaid query customgov proposals --log_level=debug