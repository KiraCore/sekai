#!/bin/bash

sekaid tx upgrade set-plan --resource-id=1 --resource-git=1 --resource-checkout=1 --resource-checksum=1 --min-halt-time=1 --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo=1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --log_level=debug --yes

# {"height":"175","txhash":"F942BD2DC74DB31334477FD4BCF8BED6A9CA173A691F9A7D1A0AB2885C72DD47","codespace":"","code":0,"data":"0A1E0A1870726F706F73652D736F6674776172652D7570677261646512020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"propose-software-upgrade\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"propose-software-upgrade"}]}]}],"info":"","gas_wanted":"0","gas_used":"15650","tx":null,"timestamp":""}

sekaid query customgov proposals --log_level=debug
# proposals:
# - content:
#     '@type': /kira.upgrade.ProposalSoftwareUpgrade
#     max_enrolment_duration: "1"
#     memo: "1"
#     min_halt_time: "1"
#     new_chain_id: "1"
#     old_chain_id: "1"
#     resources:
#       checkout: "1"
#       checksum: "1"
#       git: "1"
#       id: "1"
#     rollback_checksum: "1"
#   description: "1"
#   enactment_end_time: "2021-06-07T03:45:27.918603Z"
#   min_enactment_end_block_height: "296"
#   min_voting_end_block_height: "177"
#   proposal_id: "1"
#   result: VOTE_RESULT_UNKNOWN
#   submit_time: "2021-06-07T03:30:27.918603Z"
#   voting_end_time: "2021-06-07T03:40:27.918603Z"

sekaid tx customgov permission whitelist-permission --permission=29 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

sekaid query upgrade current-plan --log_level=debug