# sekai
Kira Hub

## Set permission environment variables

```sh
# permissions
export PermZero=0
export PermSetPermissions=1
export PermClaimValidator=2
export PermClaimCouncilor=3
export PermCreateSetPermissionsProposal=4
export PermVoteSetPermissionProposal=5
export PermUpsertTokenAlias=6
export PermChangeTxFee=7
export PermUpsertTokenRate=8

# transaction_type
export TypeMsgSend      = "send"
export TypeMsgMultiSend = "multisend"
export MsgTypeProposalSetNetworkProperty = "proposal-set-network-property"
export MsgTypeProposalAssignPermission   = "proposal-assign-permission"
export MsgTypeProposalUpsertDataRegistry = "proposal-upsert-data-registry"
export MsgTypeProposalUpsertTokenAlias   = "proposal-upsert-token-alias"
export MsgTypeVoteProposal               = "vote-proposal"
export MsgTypeWhitelistPermissions = "whitelist-permissions"
export MsgTypeBlacklistPermissions = "blacklist-permissions"
export MsgTypeClaimCouncilor       = "claim-councilor"
export MsgTypeSetNetworkProperties = "set-network-properties"
export MsgTypeSetExecutionFee      = "set-execution-fee"
export MsgTypeCreateRole = "create-role"
export MsgTypeAssignRole = "assign-role"
export MsgTypeRemoveRole = "remove-role"
export MsgTypeWhitelistRolePermission       = "whitelist-role-permission"
export MsgTypeBlacklistRolePermission       = "blacklist-role-permission"
export MsgTypeRemoveWhitelistRolePermission = "remove-whitelist-role-permission"
export MsgTypeRemoveBlacklistRolePermission = "remove-blacklist-role-permission"
export MsgTypeClaimValidator = "claim-validator"
export MsgTypeUpsertTokenAlias = "upsert-token-alias"
export MsgTypeUpsertTokenRate  = "upsert-token-rate"

export FuncIDMsgSend   = 1
export FuncIDMultiSend = 2
export FuncIDMsgProposalSetNetworkProperty = 3
export FuncIDMsgProposalAssignPermission   = 4
export FuncIDMsgProposalUpsertDataRegistry = 5
export FuncIDMsgVoteProposal               = 6
export FuncIDMsgWhitelistPermissions = 7
export FuncIDMsgBlacklistPermissions = 8
export FuncIDMsgClaimCouncilor       = 9
export FuncIDMsgSetNetworkProperties = 10
export FuncIDMsgSetExecutionFee      = 11
export FuncIDMsgCreateRole = 12
export FuncIDMsgAssignRole = 13
export FuncIDMsgRemoveRole = 14
export FuncIDMsgWhitelistRolePermission       = 15
export FuncIDMsgBlacklistRolePermission       = 16
export FuncIDMsgRemoveWhitelistRolePermission = 17
export FuncIDMsgRemoveBlacklistRolePermission = 18
export FuncIDMsgClaimValidator = 19
export FuncIDMsgUpsertTokenAlias = 20
export FuncIDMsgUpsertTokenRate  = 21
export FuncIDMsgProposalUpsertTokenAlias = 22
```

## Set ChangeTxFee permission
```sh
# command to set PermChangeTxFee permission
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid
# good response
"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"whitelist-permissions\"}]}]}]"
```

## Query permission of an address
```sh
# command
sekaid query customgov permissions $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# response
blacklist: []
whitelist:
- 4
- 3
```
## Set network properties
```sh

# command with fee set
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid

# no error response
"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set-network-properties\"}]}]}]"

# response when not enough permissions to change tx fee
"failed to execute message; message index: 0: PermChangeTxFee: not enough permissions"

# command without fee set
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --keyring-backend=test --chain-id=testing --home=$HOME/.sekaid

# response
"fee out of range [1, 10000]: invalid request"

```
## Query network properties
```sh
# command
sekaid query customgov network-properties

# response
properties:
  max_tx_fee: "10000"
  min_tx_fee: "1"
```

## Set Execution Fee
```sh
# command
sekaid tx customgov set-execution-fee --from validator --execution_name="B" --transaction_type="B" --execution_fee=10 --failure_fee=1 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=10ukex --home=$HOME/.sekaid

# response
"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set-execution-fee\"}]}]}]"
```

## Set execution fee validation test
```sh
# command for setting execution fee
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="set-network-properties" --execution_fee=10000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid

Here, the value should be looked at is `--execution_name="set-network-properties"`, `--execution_fee=10000` and `--failure_fee=1000`.

# check execution fee validation
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid
# response
"fee is less than failure fee 1000: invalid request"

Here, the value should be looked at is `"fee is less than failure fee 1000: invalid request"`.
In this case, issue is found on ante step and fee is not being paid at all.

# preparation for networks (v1) failure=1000, execution=10000
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="set-network-properties" --execution_fee=10000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# preparation for networks (v2) failure=1000, execution=500
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="set-network-properties" --execution_fee=500 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# init user1 with 100000ukex
sekaid keys add user1 --keyring-backend=test --home=$HOME/.sekaid
sekaid tx bank send validator $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) 100000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) --yes

# try changing set-network-properties with user1 that does not have ChangeTxFee permission
sekaid tx customgov set-network-properties --from user1 --min_tx_fee="2" --max_tx_fee="25000" --keyring-backend=test --chain-id=testing --fees=1000ukex --home=$HOME/.sekaid --yes
# this should fail and balance should be (previousBalance - failureFee)
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid)

# whitelist user1's permission for ChangeTxFee and try again
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov set-network-properties --from user1 --min_tx_fee="2" --max_tx_fee="25000" --keyring-backend=test --chain-id=testing --fees=1000ukex --home=$HOME/.sekaid --yes
# this should fail and balance should be (previousBalance - successFee)
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid)
```

## Query execution fee
```sh
sekaid query customgov execution-fee <msg_type>
# command
sekaid query customgov execution-fee "B"
# response
fee:
  default_parameters: "0"
  execution_fee: "10"
  failure_fee: "1"
  name: ABC
  timeout: "10"
  transaction_type: B

# genesis fee configuration test
sekaid query customgov execution-fee "A"
fee:
  default_parameters: "0"
  execution_fee: "10"
  failure_fee: "1"
  name: Claim Validator Seat
  timeout: "10"
  transaction_type: A
```

## Upsert token alias
```sh
# set PermUpsertTokenAlias permission to validator address
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenAlias --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# run upsert alias
sekaid tx tokens upsert-alias --from validator --keyring-backend=test --expiration=0 --enactment=0 --allowed_vote_types=0,1 --symbol="ETH" --name="Ethereum" --icon="myiconurl" --decimals=6 --denoms="finney" --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
```
# Query token alias
```sh
# command
sekaid query tokens alias KEX
# response
allowed_vote_types:
- "yes"
- "no"
decimals: 6
denoms:
- ukex
enactment: 0
expiration: 0
icon: myiconurl
name: Kira
status: undefined
symbol: KEX
```
```sh
# command
sekaid query tokens alias KE
# response
Error: KE symbol does not exist
```
```sh
# command
sekaid query tokens all-aliases --chain-id=testing --home=$HOME/.sekaid
# response
data:
- allowed_vote_types:
  - "yes"
  - "no"
  decimals: 6
  denoms:
  - ukex
  enactment: 0
  expiration: 0
  icon: myiconurl
  name: Kira
  status: undefined
  symbol: KEX

# command
sekaid query tokens aliases-by-denom ukex --chain-id=testing --home=$HOME/.sekaid
# response
data:
  ukex:
    allowed_vote_types:
    - "yes"
    - "no"
    decimals: 6
    denoms:
    - ukex
    enactment: 0
    expiration: 0
    icon: myiconurl
    name: Kira
    status: undefined
    symbol: KEX
```

## Upsert token rates
```sh
# set PermUpsertTokenRate permission to validator address
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# run upsert rate
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="mykex" --rate="1.5" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
```
# Query token rate
```sh
# command
sekaid query tokens rate mykex
# response
denom: mykex
fee_payments: true
rate: "1.500000"
```
```sh
# command
sekaid query tokens rate invalid_denom
# response
Error: invalid_denom denom does not exist
```
```sh
# command
sekaid query tokens all-rates --chain-id=testing --home=$HOME/.sekaid
# response
data:
- denom: ubtc
  fee_payments: true
  rate: "0.000010"
- denom: ukex
  fee_payments: true
  rate: "1.000000"
- denom: xeth
  fee_payments: true
  rate: "0.000100"

# command
sekaid query tokens rates-by-denom ukex --chain-id=testing --home=$HOME/.sekaid
# response
data:
  ukex:
    denom: ukex
    fee_payments: true
    rate: "1.000000"
```
# Fee payment in foreign currency
```sh
# register stake token as 1ukex=100stake
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="stake" --rate="0.01" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
sekaid query tokens rate stake
# try to spend stake token as fee
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="valstake" --rate="0.01" --fee_payments=true --chain-id=testing --fees=10000stake --home=$HOME/.sekaid  --yes
# smaller amount of fee in foreign currency
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="valstake" --rate="0.02" --fee_payments=true --chain-id=testing --fees=1000stake --home=$HOME/.sekaid  --yes
# try to spend unregistered token (validatortoken) as fee
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="valstake" --rate="0.03" --fee_payments=true --chain-id=testing --fees=1000validatortoken --home=$HOME/.sekaid  --yes
```

# Fee payment in foreign currency returning failure - execution fee in foreign currency
```sh
# register stake token as 1ukex=100stake
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="stake" --rate="0.01" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
sekaid query tokens rate stake

# set execution fee and failure fee for upsert-rate transaction
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermChangeTxFee --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# set execution_fee=1000 failure_fee=5000
sekaid tx customgov set-execution-fee --from validator --execution_name="upsert-token-alias" --transaction_type="upsert-token-alias" --execution_fee=1000 --failure_fee=5000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# set execution_fee=5000 failure_fee=1000
sekaid tx customgov set-execution-fee --from validator --execution_name="upsert-token-alias" --transaction_type="upsert-token-alias" --execution_fee=5000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

# check current balance
sekaid query bank balances $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# try upsert-token-alias failure in foreign currency
sekaid tx tokens upsert-alias --from validator --keyring-backend=test --expiration=0 --enactment=0 --allowed_vote_types=0,1 --symbol="ETH" --name="Ethereum" --icon="myiconurl" --decimals=6 --denoms="finney" --chain-id=testing --fees=500000stake --home=$HOME/.sekaid  --yes
# set permission for this execution
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenAlias --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=10000stake --home=$HOME/.sekaid --yes
# try upsert-token-alias success in foreign currency
sekaid tx tokens upsert-alias --from validator --keyring-backend=test --expiration=0 --enactment=0 --allowed_vote_types=0,1 --symbol="ETH" --name="Ethereum" --icon="myiconurl" --decimals=6 --denoms="finney" --chain-id=testing --fees=500000stake --home=$HOME/.sekaid  --yes
```

# Query validator account
```sh
# query validator account
sekaid query validator --addr  $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)
```

# Custom governance module commands

```sh
sekaid tx customgov councilor claim-seat --from validator --keyring-backend=test --home=$HOME/.sekaid

sekaid tx customgov permission blacklist-permission
sekaid tx customgov permission whitelist-permission

sekaid tx customgov proposal assign-permission
sekaid tx customgov proposal vote

sekaid tx customgov role blacklist-permission
sekaid tx customgov role create
sekaid tx customgov role remove
sekaid tx customgov role remove-blacklist-permission
sekaid tx customgov role remove-whitelist-permission
sekaid tx customgov role whitelist-permission

# querying for voters of a specific proposal
sekaid query customgov voters 1
# querying for votes of a specific proposal
sekaid query customgov votes 1
# querying for a vote of a specific propsal/voter pair
sekaid query customgov vote 1 $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)
```

# Commands for adding more validators

```sh
# sekaid keys add val2 --keyring-backend=test --home=$HOME/.sekaid
# sekaid tx bank send validator $(sekaid keys show -a val2 --keyring-backend=test --home=$HOME/.sekaid) 100000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermCreateSetPermissionsProposal --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermVoteSetPermissionProposal --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes

sekaid tx customgov proposal assign-permission $PermClaimValidator --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 

sekaid tx claim-validator-seat --from validator --keyring-backend=test --home=$HOME/.sekaid --validator-key=kiravaloper1ntk7n5y38en5dvnhvmruwagmkemq76x8s4pnwu --moniker="validator" --chain-id=testing --fees=100ukex --yes

# sekaid tx claim-validator-seat --from val2 --keyring-backend=test --home=$HOME/.sekaid --pubkey=kiravalconspub1zcjduepqdllep3v5wv04hmu987rv46ax7fml65j3dh5tf237ayn5p59jyamq04048n --validator-key=kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w --moniker="val2" --chain-id=testing --fees=100ukex --yes
# sekaid tx claim-validator-seat --from val2 --keyring-backend=test --home=$HOME/.sekaid --validator-key=kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w --moniker="val2" --chain-id=testing --fees=100ukex --yes
```

# Tx for set network property proposal

```
sekaid tx customgov proposal set-network-property MIN_TX_FEE 101 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
```
---
`dev` branch
