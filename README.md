# sekai
Kira Hub

## Set ChangeTxFee permission
```sh
# command to set changeTxFee permission
sekaid tx customgov set-whitelist-permissions --from validator --keyring-backend=test --permission=4 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid
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
sekaid tx customgov set-execution-fee --from validator --execution_name="ABC" --transaction_type="B" --execution_fee=10 --failure_fee=1 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=10ukex --home=$HOME/.sekaid

# response
"[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"set-execution-fee\"}]}]}]"
```

## Set execution fee validation test
```sh
# command for setting execution fee
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="B" --execution_fee=10000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid

Here, the value should be looked at is `--execution_name="set-network-properties"`, `--execution_fee=10000` and `--failure_fee=1000`.

# check execution fee validation
sekaid tx customgov set-network-properties --from validator --min_tx_fee="2" --max_tx_fee="20000" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid
# response
"fee is less than failure fee 1000: invalid request"

Here, the value should be looked at is `"fee is less than failure fee 1000: invalid request"`.
In this case, issue is found on ante step and fee is not being paid at all.

# preparation for networks (v1) failure=1000, execution=10000
sekaid tx customgov set-whitelist-permissions --from validator --keyring-backend=test --permission=4 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="B" --execution_fee=10000 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y

# preparation for networks (v2) failure=1000, execution=500
sekaid tx customgov set-whitelist-permissions --from validator --keyring-backend=test --permission=4 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y
sekaid tx customgov set-execution-fee --from validator --execution_name="set-network-properties" --transaction_type="B" --execution_fee=500 --failure_fee=1000 --timeout=10 default_parameters=0 --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y

# init user1 with 100000ukex
sekaid keys add user1 --keyring-backend=test --home=$HOME/.sekaid
sekaid tx bank send validator $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) 100000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) <<< y

# try changing set-network-properties with user1 that does not have ChangeTxFee permission
sekaid tx customgov set-network-properties --from user1 --min_tx_fee="2" --max_tx_fee="25000" --keyring-backend=test --chain-id=testing --fees=1000ukex --home=$HOME/.sekaid <<< y
# this should fail and balance should be (previousBalance - failureFee)
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid)

# whitelist user1's permission for ChangeTxFee and try again
sekaid tx customgov set-whitelist-permissions --from validator --keyring-backend=test --permission=4 --addr=$(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y
sekaid tx customgov set-network-properties --from user1 --min_tx_fee="2" --max_tx_fee="25000" --keyring-backend=test --chain-id=testing --fees=1000ukex --home=$HOME/.sekaid <<< y
# this should fail and balance should be (previousBalance - successFee)
sekaid query bank balances $(sekaid keys show -a user1 --keyring-backend=test --home=$HOME/.sekaid)
```

## Query execution fee
```sh
sekaid query customgov execution-fee <msg_type>
# command
sekaid query customgov execution-fee ABC
# response
fee:
  default_parameters: "0"
  execution_fee: "10"
  failure_fee: "1"
  name: ABC
  timeout: "10"
  transaction_type: B

# genesis fee configuration test
sekaid query customgov execution-fee "Claim Validator Seat"
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
# set PermUpsertTokenAlias(10) permission to validator address
sekaid tx customgov set-whitelist-permissions --from validator --keyring-backend=test --permission=10 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid <<< y
# run upsert alias
sekaid tx tokens upsert-alias --from validator --keyring-backend=test --expiration=0 --enactment=0 --allowed_vote_types=0,1 --symbol="ETH" --name="Ethereum" --icon="myiconurl" --decimals=6 --denoms="finney" --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  <<< y
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
---
`dev` branch
