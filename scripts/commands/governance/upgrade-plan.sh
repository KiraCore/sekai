#!/bin/bash

sekaid tx upgrade propose-upgrade-plan --name="upgrade1" --instate-upgrade=true --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=1626540651 --height=15  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade1 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --log_level=debug --yes

# {"height":"175","txhash":"F942BD2DC74DB31334477FD4BCF8BED6A9CA173A691F9A7D1A0AB2885C72DD47","codespace":"","code":0,"data":"0A1E0A1870726F706F73652D736F6674776172652D7570677261646512020801","raw_log":"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"propose-software-upgrade\"}]}]}]","logs":[{"msg_index":0,"log":"","events":[{"type":"message","attributes":[{"key":"action","value":"propose-software-upgrade"}]}]}],"info":"","gas_wanted":"0","gas_used":"15650","tx":null,"timestamp":""}


sekaid tx customgov permission whitelist-permission --permission=29 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

sekaid query upgrade show-plan --log_level=debug

sekaid query customgov proposals --log_level=debug

# After chain halt detect, get new binary and install
# Binary's app.go should contain instate upgrade handler like below
#     app.upgradeKeeper.SetUpgradeHandler(
#         "upgrade1", func(ctx sdk.Context, plan upgradetypes.Plan) {
#         })

# Once it has been setup, start binary again
sekaid start --home=$HOME/.sekaid

# propose second upgrade plan
sekaid tx upgrade propose-upgrade-plan --name="upgrade2" --instate-upgrade=true --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=1626540651 --height=40  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade2 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --log_level=debug --yes

sekaid tx upgrade propose-upgrade-plan --name="upgrade3" --instate-upgrade=true --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=1626540651 --height=70  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade3 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --log_level=debug --yes

sekaid tx upgrade propose-upgrade-plan --name="upgrade4" --instate-upgrade=true --resources="[{\"id\":\"infra\",\"git\":\"https://aaa/bbb.com\"}]" --min-upgrade-time=1626540651 --height=90  --old-chain-id=1 --new-chain-id=1 --rollback-memo=1 --max-enrollment-duration=1 --upgrade-memo="upgrade4 test" --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --log_level=debug --yes
