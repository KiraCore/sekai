# sekai
Kira Hub

## Set environment variables

```sh
sh env.sh
```

### Set permission via governance process

```sh
sekaid tx customgov proposal assign-permission $PermClaimValidator --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

sekaid query customgov proposals
sekaid query customgov proposal 1

sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
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

# Query signing infos per validator's consensus address
```sh
# query all signing infos as an array
sekaid query customslashing signing-infos
# response
info:
- address: kiravalcons166p6nw8gm4cq7xescmyzm8qsqf56za0x5q6ep9
  inactive_until: "1970-01-01T00:00:00Z"
  index_offset: "2"
  missed_blocks_counter: "0"
  start_height: "0"
  tombstoned: false
pagination:
  next_key: null
  total: "0"
```
```sh
# query signing info by validator
sekaid query customslashing signing-info $(sekaid tendermint show-validator)
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

# Commands for poor network management
```sh
# create proposal for setting poor network msgs
sekaid tx customgov proposal set-poor-network-msgs AAA,BBB --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=1000ukex --yes
# query for proposals
sekaid query customgov proposals
# set permission to vote on proposal
sekaid tx customgov permission whitelist-permission --permission=19 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# vote on the proposal
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# check votes
sekaid query customgov votes 1 
# wait until vote end time finish
sekaid query customgov proposals
# query poor network messages
sekaid query customgov poor-network-messages

# whitelist permission for modifying network properties
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=7 --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# test poor network messages after modifying min_validators section
sekaid tx customgov set-network-properties --from validator --min_validators="2" --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# set permission for upsert token rate
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermUpsertTokenRate --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# try running upser token rate which is not allowed on poor network
sekaid tx tokens upsert-rate --from validator --keyring-backend=test --denom="mykex" --rate="1.5" --fee_payments=true --chain-id=testing --fees=100ukex --home=$HOME/.sekaid  --yes
# try sending more than allowed amount via bank send
sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# try setting network property by governance to allow more amount sending
sekaid tx customgov proposal set-network-property POOR_NETWORK_MAX_BANK_SEND 100000000 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
# try sending after modification of poor network bank send param
sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000000ukex --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
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

# get ValAddress (kiravaloperxxx) from validator key
sekaid val-address $(sekaid keys show -a validator --keyring-backend=test)

# sekaid tx claim-validator-seat --from val2 --keyring-backend=test --home=$HOME/.sekaid --pubkey=kiravalconspub1zcjduepqdllep3v5wv04hmu987rv46ax7fml65j3dh5tf237ayn5p59jyamq04048n --validator-key=kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w --moniker="val2" --chain-id=testing --fees=100ukex --yes
# sekaid tx claim-validator-seat --from val2 --keyring-backend=test --home=$HOME/.sekaid --validator-key=kiravaloper1ewgq8gtsefakhal687t8hnsw5zl4y8eksup39w --moniker="val2" --chain-id=testing --fees=100ukex --yes
```

# Tx for set network property proposal

```
sekaid tx customgov proposal set-network-property MIN_TX_FEE 101 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
```

# Tx for Unjailing

Modify genesis json to have jailed validator for Unjail testing

```json
{
  "accounts": [
    {
      "@type": "/cosmos.auth.v1beta1.BaseAccount",
      "address": "kira126f48ukar7ntqqvk0qxgd3juu7r42mjmsddjrq",
      "pub_key": null,
      "account_number": "0",
      "sequence": "0"
    }
  ],
  "balances": [
    {
      "address": "kira126f48ukar7ntqqvk0qxgd3juu7r42mjmsddjrq",
      "coins": [
        {
          "denom": "stake",
          "amount": "1000000000"
        },
        {
          "denom": "ukex",
          "amount": "1000000000"
        },
        {
          "denom": "validatortoken",
          "amount": "1000000000"
        }
      ]
    }
  ],
  "validators": [
    {
      "moniker": "hello2",
      "website": "",
      "social": "social2",
      "identity": "",
      "commission": "1.000000000000000000",
      "val_key": "kiravaloper126f48ukar7ntqqvk0qxgd3juu7r42mjmrt33mv",
      "pub_key": {
        "@type": "/cosmos.crypto.ed25519.PubKey",
        "key": "tC8mzxDI3bzfZtToxU6ZpZIOw6nqQx87OZ1fD6FpD7E="
      },
      "status": "JAILED",
      "rank": "0",
      "streak": "0"
    }
  ],
  "network_actors": [
    {
      "address": "kira126f48ukar7ntqqvk0qxgd3juu7r42mjmsddjrq",
      "roles": ["1"],
      "status": "ACTIVE",
      "votes": [
        "VOTE_OPTION_YES",
        "VOTE_OPTION_ABSTAIN",
        "VOTE_OPTION_NO",
        "VOTE_OPTION_NO_WITH_VETO"
      ],
      "permissions": {
        "blacklist": [],
        "whitelist": []
      },
      "skin": "1"
    }
  ],
}
```

Add jailed validator key to kms.
```sh
  sekaid keys add jailed_validator --keyring-backend=test --home=$HOME/.sekaid --recover
  "dish rather zoo connect cross inhale security utility occur spell price cute one catalog coconut sort shuffle palm crop surface label foster slender inherit"
```

```sh
# make proposal to unjail validator from jailed_validator
sekaid tx customstaking proposal proposal-unjail-validator hash reference --from=jailed_validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

# vote on unjail validator proposal
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes

# proposal for jail max time - max to 1440min = 1d
sekaid tx customgov proposal set-network-property JAIL_MAX_TIME 1440 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes
```

# Proposal Tx for freeze / unfreeze tokens

```sh
# create a proposal to blacklist validatortoken
sekaid tx tokens propose-update-tokens-blackwhite --is_blacklist=true --is_add=true --tokens=validatortoken --tokens=kava --from validator --chain-id=testing --keyring-backend=test --fees=100ukex --home=$HOME/.sekaid --yes
# check proposal ID
sekaid query customgov proposals
# whitelist permission to vote on proposal
sekaid tx customgov permission whitelist-permission --from validator --keyring-backend=test --permission=$PermVoteTokensWhiteBlackChangeProposal --addr=$(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# vote on proposal
sekaid tx customgov proposal vote 1 1 --from validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes 
# get all votes
sekaid query customgov vote 1 $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)
```

# Query for current frozen / unfronzen tokens [these values are valid only when specific network property is enabled]

```sh
# query token blacklists and whitelists
sekaid query tokens token-black-whites
# response
data:
  blacklisted:
  - frozen
  whitelisted:
  - ukex

# query network properties
sekaid query customgov network-properties
# response
properties:
  enable_foreign_fee_payments: true
  enable_token_blacklist: false # useful for blacklist use or not
  enable_token_whitelist: false # useful for whitelist use or not
  inactive_rank_decrease_percent: "50"
  jail_max_time: "10"
  max_tx_fee: "1000000"
  min_tx_fee: "100"
  min_validators: "1"
  mischance_rank_decrease_amount: "10"
  poor_network_max_bank_send: "1000000"
  proposal_enactment_time: "300"
  proposal_end_time: "600"
  vote_quorum: "33"

# try sending frozen token
sekaid tx bank send validator $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid) 100000frozen --keyring-backend=test --chain-id=testing --fees=100ukex --home=$HOME/.sekaid --yes
# response
token is frozen: invalid request
```
