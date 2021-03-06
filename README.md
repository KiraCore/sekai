# sekai
Kira Hub

## Set environment variables

```sh
sh env.sh
```
# Adding more validators
[scripts/commands/adding-validators.sh](scripts/commands/adding-validators.sh)
## Set ChangeTxFee permission
[scripts/commands/set-permission.sh](scripts/commands/set-permission.sh)
## Set network properties
[scripts/commands/set-network-properties.sh](scripts/commands/set-network-properties.sh)
## Set Execution Fee
[scripts/commands/set-execution-fee.sh](scripts/commands/set-execution-fee.sh)
## Upsert token rates
[scripts/commands/upsert-token-rates.sh](scripts/commands/upsert-token-rates.sh)
## Upsert token alias
[scripts/commands/upsert-token-alias.sh](scripts/commands/upsert-token-alias.sh)
# Fee payment in foreign currency
[scripts/commands/foreign-fee-payments.sh](scripts/commands/foreign-fee-payments.sh)
# Fee payment in foreign currency returning failure - execution fee in foreign currency
[scripts/commands/foreign-fee-payments-failure-return.sh](scripts/commands/foreign-fee-payments-failure-return.sh)
## Query permission of an address
[scripts/commands/query-permission.sh](scripts/commands/query-permission.sh)
## Query network properties
[scripts/commands/query-network-properties.sh](scripts/commands/query-network-properties.sh)
## Query execution fee
[scripts/commands/query-execution-fee.sh](scripts/commands/query-execution-fee.sh)
# Query token alias
[scripts/commands/query-token-alias.sh](scripts/commands/query-token-alias.sh)
# Query token rate
[scripts/commands/query-token-rate.sh](scripts/commands/query-token-rate.sh)
# Query validator account
[scripts/commands/query-validator.sh](scripts/commands/query-validator.sh)
# Query for current frozen / unfronzen tokens
**Notes**: these values are valid only when specific network property is enabled
[scripts/commands/query-frozen-token.sh](scripts/commands/query-frozen-token.sh)
# Query poor network messages
[scripts/commands/query-poor-network-messages.sh](scripts/commands/query-poor-network-messages.sh)
# Query signing infos per validator's consensus address
[scripts/commands/query-signing-infos.sh](scripts/commands/query-signing-infos.sh)
# Common commands for governance process
[scripts/commands/governance/common.sh](scripts/commands/governance/common.sh)
### Set permission via governance process
[scripts/commands/governance/assign-permission.sh](scripts/commands/governance/assign-permission.sh)
## Upsert token alias via governance process
[scripts/commands/governance/upsert-token-alias.sh](scripts/commands/governance/upsert-token-alias.sh)
## Upsert token rates via governance process
[scripts/commands/governance/upsert-token-rates.sh](scripts/commands/governance/upsert-token-rates.sh)
# Commands for poor network management via governance process
[scripts/commands/governance/poor-network-messages.sh](scripts/commands/governance/poor-network-messages.sh)
# Freeze / unfreeze tokens via governance process
[scripts/commands/governance/token-freeze.sh](scripts/commands/governance/token-freeze.sh)
# Set network property proposal via governance process
[scripts/commands/governance/set-network-property.sh](scripts/commands/governance/set-network-property.sh)

# Set application upgrade proposal via governance process
[scripts/commands/governance/upgrade-plan.sh](scripts/commands/governance/upgrade-plan.sh)

Export the status of chain before halt (should kill the daemon process at the time of genesis export)
[scripts/commands/export-state.sh](scripts/commands/export-state.sh)

The script for creating new chain from exported state should be written or manual edition process is required.
`ChainId` should be modified in this process.

For now, upgrade process requires manual conversion from old genesis to new genesis.
At each time of upgrade, genesis upgrade command will be built and infra could run the command like `sekaid genesis-migrate`

Note: export is not exporting the upgrade plan and if all validators run with exported genesis with the previous binary, consensus failure won't happen.

TODO:@Asmodat If we need to export upgrade plan as well, export should also have a rollback flag in case we want to get genesis that allows to continue producing blocks on the old chain after chain was halted but upgrade failed This way no matter what happens chain remains operational.

# Unjail via governance process

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

[scripts/commands/governance/unjail-validator.sh](scripts/commands/governance/unjail-validator.sh)
