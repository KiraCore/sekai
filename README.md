# sekai
Kira Hub

## Set environment variables

```sh
sh env.sh
```
# Get version info
[scripts/commands/version.sh](scripts/commands/version.sh)
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

Note: state export command is not exporting the upgrade plan and if all validators run with exported genesis with the previous binary, consensus failure won't happen.

# Identity registrar
[scripts/commands/identity-registrar.sh](scripts/commands/identity-registrar.sh)

# Unjail via governance process

Modify genesis json to have jailed validator for Unjail testing
Add jailed validator key to kms.
```sh
  sekaid keys add jailed_validator --keyring-backend=test --home=$HOME/.sekaid --recover
  "dish rather zoo connect cross inhale security utility occur spell price cute one catalog coconut sort shuffle palm crop surface label foster slender inherit"
```

[scripts/commands/governance/unjail-validator.sh](scripts/commands/governance/unjail-validator.sh)

# New genesis file generation process from exported version

In order to manually generate new genesis file when the hard fork is activated, following steps should be taken:

1. Export current genesis, e.g: sekaid export --home=<path>
2. Change chain-id to new_chain_id as indicated by the upgrade plan
3. Replace current upgrade plan in the app_state.upgrade with next plan and set next plan to null

Using a command it can be done in this way.
1. sekaid export > exported-genesis.json
2. sekaid new-genesis-from-exported exported-genesis.json new-genesis.json
