#!/bin/bash

# command to generate recovery secret
sekaid tx recovery generate-recovery-secret 10a0fbe01030000122300000000000
NONCE=00
PROOF=29eeb09666c4792c314e631063932621f573cbb07af7274657d1314e1892eb93
CHALLENGE=220c7e47a53ef4c2161035308d4fdc52f619e54f8a657b208ba3708139fdc03d

# commmand to create an account
sekaid keys add recovery --keyring-backend=test
RECOVERY=$(sekaid keys show -a recovery --keyring-backend=test)

# commmand to create an account
sekaid keys add controller --keyring-backend=test
CONTROLLER=$(sekaid keys show -a controller --keyring-backend=test)

# commands to query
sekaid query recovery recovery-record $(sekaid keys show -a validator --keyring-backend=test)

# command to register recovery secret
sekaid tx recovery register-recovery-secret $CHALLENGE $NONCE $PROOF $CONTROLLER --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block

# command to rotate
sekaid tx recovery rotate-recovery-address $RECOVERY $PROOF $(sekaid keys show -a validator --keyring-backend=test) --from=validator --chain-id=testing --keyring-backend=test --fees=1000ukex -y --home=$HOME/.sekaid --yes --broadcast-mode=block

# command to check validators after rotation
sekaid query customstaking validators --moniker="" --addr="" --val-addr="" --moniker="" --status="" --pubkey="" --proposer=""

# upsert staking pool before rotation
sekaid tx multistaking upsert-staking-pool kiravaloper18ka9xpvwh75sgldgke69jmxsnkhjm0wa3ns9xa --commission=0.5 --from=validator --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block

# recovery token queries
sekaid query recovery rr-holder-rewards $(sekaid keys show -a validator --keyring-backend=test)
sekaid query recovery rr-holders rr/hello
sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test)

# recovery token txs
sekaid keys add recovery --keyring-backend=test
sekaid keys add recovery2 --keyring-backend=test
RECOVERY=$(sekaid keys show -a recovery --keyring-backend=test)
RECOVERY2=$(sekaid keys show -a recovery2 --keyring-backend=test)
VALIDATOR=$(sekaid keys show -a validator --keyring-backend=test)

sekaid tx recovery issue-recovery-tokens --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block
sekaid tx recovery burn-recovery-tokens 10000000rr/hello --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block
sekaid tx recovery register-rrtoken-holder --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block
sekaid tx recovery claim-rrtoken-rewards --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block
sekaid tx recovery rotate-validator-by-half-rr-holder $VALIDATOR $RECOVERY --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block

# upsert-staking-pool after recovery
sekaid tx bank send validator $RECOVERY 10000000000000rr/hello,100000ukex --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block

sekaid tx multistaking upsert-staking-pool kiravaloper1zwyxw66aw0fd3xv9e9ewyk58c0qdcn5synp3dt --commission=0.5 --from=recovery --keyring-backend=test --fees=100ukex --chain-id=testing -y --broadcast-mode=block
sekaid tx recovery rotate-validator-by-half-rr-holder $RECOVERY $VALIDATOR --from=recovery --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block
sekaid tx recovery rotate-validator-by-half-rr-holder $RECOVERY $RECOVERY2 --from=recovery --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block
