#!/bin/bash

# create id.json file with below content
# {
#     "key1": "value1",
#     "key2": "value2"
# }

sekaid tx customgov register-identity-records --infos-file="id.json" --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid tx customgov register-identity-records --infos-json='{"moniker":"My Moniker","social":"My Social"}' --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid query customgov all-identity-records
sekaid query customgov identity-record 1
sekaid query customgov identity-records-by-addr $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# query by specific keys
sekaid query customgov identity-records-by-addr $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --keys="moniker"

sekaid tx customgov delete-identity-records --keys="moniker" --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid keys add test --keyring-backend=test --home=$HOME/.sekaid
sekaid tx customgov request-identity-record-verify --record-ids=1 --verifier=$(sekaid keys show -a test --keyring-backend=test --home=$HOME/.sekaid) --tip=200ukex --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid tx customgov request-identity-record-verify --record-ids=2 --verifier=$(sekaid keys show -a test --keyring-backend=test --home=$HOME/.sekaid) --tip=200ukex --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid query customgov all-identity-record-verify-requests
sekaid query customgov identity-record-verify-request 1

sekaid query customgov identity-record-verify-requests-by-requester $(sekaid keys show -a test --keyring-backend=test --home=$HOME/.sekaid)
sekaid query customgov identity-record-verify-requests-by-requester $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

sekaid query customgov identity-record-verify-requests-by-approver $(sekaid keys show -a test --keyring-backend=test --home=$HOME/.sekaid)
sekaid query customgov identity-record-verify-requests-by-approver $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

sekaid tx customgov handle-identity-records-verify-request 1 --from=validator --approve=true --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid tx customgov handle-identity-records-verify-request 2 --from=validator --approve=false --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes

sekaid tx customgov request-identity-record-verify --record-ids=1 --verifier=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --tip=200ukex --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid query customgov all-identity-record-verify-requests

sekaid tx customgov cancel-identity-records-verify-request 2 --from=validator --keyring-backend=test --home=$HOME/.sekaid --fees=100ukex --chain-id=testing --yes
sekaid query customgov all-identity-record-verify-requests
