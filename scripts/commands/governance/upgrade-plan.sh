#!/bin/bash

# PermCreateSoftwareUpgradeProposal PermValue = 28
# PermVoteSoftwareUpgradeProposal PermValue = 29
sekaid tx customgov permission whitelist-permission --permission=28 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
sekaid tx customgov permission whitelist-permission --permission=29 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

# The upgrade time should be set to future time if not it cause internal error
sekaid tx upgrade proposal-set-plan --name="upgrade1" --instate-upgrade=true --skip-handler=false --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=$(($(date -u +%s) + 200))  --old-chain-id=testing --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade1 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

# timestamp related commands
# echo $(date -u +%s)
# echo $(($(date -u +%s) + 200))

# {"height":"175","txhash":"F942BD2DC74DB31334477FD4BCF8BED6A9CA173A691F9A7D1A0AB2885C72DD47","codespace":"","code":0,"data":"0A1E0A1870726F706F73652D736F6674776172652D7570677261646512020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"propose-software-upgrade\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"propose-software-upgrade"}]}]}],"info":"","gas_wanted":"0","gas_used":"15650","tx":null,"timestamp":""}

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

sekaid query upgrade next-plan

sekaid query customgov proposals

# After chain halt detect, get new binary and install
# Binary's app.go should contain instate upgrade handler like below
#     app.upgradeKeeper.SetUpgradeHandler(
#         "upgrade1", func(ctx sdk.Context, plan upgradetypes.Plan) {
#         })

# Once it has been setup, start binary again
sekaid start --home=$HOME/.sekaid

# propose second upgrade plan
sekaid tx upgrade proposal-set-plan --name="upgrade2" --instate-upgrade=true --skip-handler=false --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=$(($(date -u +%s) + 200))  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade2 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx upgrade proposal-set-plan --name="upgrade3" --instate-upgrade=true --skip-handler=false --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=$(($(date -u +%s) + 200))  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade3 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx upgrade proposal-set-plan --name="upgrade4" --instate-upgrade=true --skip-handler=true --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=$(($(date -u +%s) + 200))  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade4 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx upgrade proposal-set-plan --name="upgrade5" --instate-upgrade=false --skip-handler=false --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=$(($(date -u +%s) + 200))  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade5 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx upgrade proposal-cancel-plan --name="cancel-upgrade4" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
sekaid tx customgov proposal vote 2 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
sekaid query customgov proposals

# upgrade to new json for hard-fork case
sekaid export > exported-genesis.json
sekaid new-genesis-from-exported exported-genesis.json new-genesis.json
