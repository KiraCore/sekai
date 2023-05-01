#!/bin/bash

# queries
sekaid query layer2 all-dapps
sekaid query layer2 execution-registrar l2dex
sekaid query layer2 transfer-dapps
sekaid query layer2 global-tokens

# transactions
sekaid tx layer2 create-dapp-proposal --dapp-name="l2dex" --denom="ul2d" --dapp-description="layer2 dex" \
  --website="website" --logo="logo" --social="social" --docs="docs" \
  --controller-roles="1" --controller-accounts="" --vote-quorum=30 --vote-period=86400 --vote-enactment=1000 \
  --bond="1000000ukex" \
  --issurance-config='{"premint":"10000","postmint":"10000","time":"1680044405"}' \
  --lp-pool-config='{"ratio": "1.0", "drip": 86400}' \
  --executors-min=1 --executors-max=3 --verifiers-min=1 \
  --binary-info='{"name":"layer2dex","hash":"0cc0","source":"github.com","reference":"","type":"exec"}' \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 

sekaid tx layer2 bond-dapp-proposal --dapp-name="l2dex" --bond="1000000ukex" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 

sekaid tx layer2 reclaim-dapp-proposal --dapp-name="l2dex" --bond="1000000ukex" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 

INTERX=$(sekaid keys show -a validator --keyring-backend=test)
sekaid tx layer2 join-dapp-verifier-with-bond "l2dex" $INTERX \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 exit-dapp "l2dex" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 execute-dapp-tx "l2dex" "l2dex.com" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

LEADER=$(sekaid keys show -a validator --keyring-backend=test)
sekaid tx layer2 denounce-leader "l2dex" $LEADER "bad actor" "v1" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 transition-dapp "l2dex" "08080818" "v1" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 approve-dapp-transition "l2dex" "v1" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 reject-dapp-transition "l2dex" "v1" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 proposal-join-dapp "l2dex" true true $INTERX --title="title" --description="description" \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid query customgov proposals
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes  --broadcast-mode=block 

sekaid tx layer2 proposal-upsert-dapp --title="title" --description="description" \
  --dapp-name="l2dex" --denom="ul2d" --dapp-description="layer2 dex" \
  --website="website" --logo="logo" --social="social" --docs="docs" \
  --controller-roles="1" --controller-accounts="" --vote-quorum=30 --vote-period=86400 --vote-enactment=1000 \
  --bond="1000000ukex" \
  --issurance-config='{"premint":"10000","postmint":"10000","time":"1680044405"}' \
  --lp-pool-config='{"ratio": "1.0", "drip": 86400}' \
  --binary-info='{"name":"layer2dex","hash":"0cc0","source":"github.com","reference":"","type":"exec"}' \
  --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 
sekaid query customgov proposals
sekaid tx customgov proposal vote 2 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes  --broadcast-mode=block 

sekaid tx layer2 redeem-dapp-pool l2dex 10lp/ul2d 0.1 --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block 
sekaid tx layer2 swap-dapp-pool l2dex 10ukex 0.1 --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 convert-dapp-pool l2dex l2agg 10lp/ul2d 0.1 --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 pause-dapp l2dex --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 unpause-dapp l2dex --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 reactivate-dapp l2dex --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 transfer-dapp '{"amounts":[{"bridge_token_index":"1","amount":"100000"}],"source_dapp":"1","source_account":"1","dest_dapp":"2","dest_beneficiary":"2","xam":""}' --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block

sekaid tx layer2 mint-create-ft-tx l2gov "L2 Governance Token" "L2GOV" "" "" "" "" 6 10000000 0 1000 $(sekaid keys show -a validator --keyring-backend=test) --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 mint-create-nft-tx l2gov "L2 Governance Token" "L2GOV" "" "" "" "" 6 10000000 0 1000 $(sekaid keys show -a validator --keyring-backend=test) "" "" --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 mint-issue-tx ku/l2gov 1000000 $(sekaid keys show -a validator --keyring-backend=test) --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
sekaid tx layer2 burn-issue-tx ku/l2gov 900000 --from=validator --chain-id=testing --fees=100ukex --keyring-backend=test --home=$HOME/.sekaid --yes --broadcast-mode=block
