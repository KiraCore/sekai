#!/bin/bash

# create id.json file with below content
# {
#     "key1": "value1",
#     "key2": "value2"
# }

sekaid tx customgov create-identity-record --infos-file="id.json" --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid tx customgov create-identity-record --infos-json='{"key":"moniker","value":"My Moniker"}' --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid query customgov all-identity-records
sekaid query customgov identity-record 1
sekaid query customgov identity-record-by-addr $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

sekaid tx customgov edit-identity-record --record-id=1 --infos-file="id.json" --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid tx customgov edit-identity-record --record-id=1 --infos-json='{"key":"moniker","value":"My Moniker"}' --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid tx customgov request-identity-record-verify --record-ids=1 --verifier=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --tip=10ukex --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid query customgov all-identity-record-verify-requests
sekaid query customgov identity-record-verify-request 1
sekaid query customgov identity-record-verify-requests-by-approver $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)
sekaid query customgov identity-record-verify-requests-by-requester $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

sekaid tx customgov approve-identity-records 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid tx customgov request-identity-record-verify --record-ids=1 --verifier=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --tip=10ukex --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid query customgov all-identity-record-verify-requests

sekaid tx customgov cancel-identity-records-verify-request 2 --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid query customgov all-identity-record-verify-requests
