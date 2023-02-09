# sekai

KIRA Relay Chain

## Quick setup from Github

```bash
# dont forget to specify branch name or hash
cd $HOME && rm -fvr ./sekai && SEKAI_BRANCH="<branch-name-or-checkout-hash>" && \
 git clone https://github.com/KiraCore/sekai.git -b $SEKAI_BRANCH && \
 cd ./sekai && chmod -R 777 ./scripts && make proto-gen && \
 make install && echo "SUCCESS installed sekaid $(sekaid version)" || echo "FAILED"
```

## Signatures

All files in KIRA repositories are always signed with [cosign](https://github.com/sigstore/cosign/releases), you should NEVER install anything on your machine unless you verified integrity of the files!

Cosign requires simple initial setup of the signer keys described more precisely [here](https://dev.to/n3wt0n/sign-your-container-images-with-cosign-github-actions-and-github-container-registry-3mni)

```bash
# install cosign
COSIGN_VERSION="v1.7.2" && \
if [[ "$(uname -m)" == *"ar"* ]] ; then ARCH="arm64"; else ARCH="amd64" ; fi && echo $ARCH && \
PLATFORM=$(uname) && FILE=$(echo "cosign-${PLATFORM}-${ARCH}" | tr '[:upper:]' '[:lower:]') && \
 wget https://github.com/sigstore/cosign/releases/download/${COSIGN_VERSION}/$FILE && chmod +x -v ./$FILE && \
 mv -fv ./$FILE /usr/local/bin/cosign && cosign version

# save KIRA public cosign key
KEYS_DIR="/usr/keys" && KIRA_COSIGN_PUB="${KEYS_DIR}/kira-cosign.pub" && \
mkdir -p $KEYS_DIR  && cat > ./cosign.pub << EOL
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/IrzBQYeMwvKa44/DF/HB7XDpnE+
f+mU9F/Qbfq25bBWV2+NlYMJv3KvKHNtu3Jknt6yizZjUV4b8WGfKBzFYw==
-----END PUBLIC KEY-----
EOL

# download desired files and the corresponding .sig file from: https://github.com/KiraCore/tools/releases

# verify signature of downloaded files
cosign verify-blob --key=$KIRA_COSIGN_PUB--signature=./<file>.sig ./<file>
```

## Install bash-utils

Bash-utils is a KIRA tool that helps to simplify shell scripts and various bash commands that you might need to run

```
# Install bash-utils.sh KIRA tool to make downloads faster and easier
TOOLS_VERSION="v0.0.12.4" && mkdir -p /usr/keys && FILE_NAME="bash-utils.sh" && \
 if [ -z "$KIRA_COSIGN_PUB" ] ; then KIRA_COSIGN_PUB=/usr/keys/kira-cosign.pub ; fi && \
 echo -e "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/IrzBQYeMwvKa44/DF/HB7XDpnE+\nf+mU9F/Qbfq25bBWV2+NlYMJv3KvKHNtu3Jknt6yizZjUV4b8WGfKBzFYw==\n-----END PUBLIC KEY-----" > $KIRA_COSIGN_PUB && \
 wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}" -O ./$FILE_NAME && \
 wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}.sig" -O ./${FILE_NAME}.sig && \
 cosign verify-blob --key="$KIRA_COSIGN_PUB" --signature=./${FILE_NAME}.sig ./$FILE_NAME && \
 chmod -v 555 ./$FILE_NAME && ./$FILE_NAME bashUtilsSetup "/var/kiraglob" && . /etc/profile && \
 echoInfo "Installed bash-utils $(bash-utils bashUtilsVersion)"
```

## Quick setup from Releases

```bash
# TBD
```

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
2. sekaid new-genesis-from-exported exported-genesis.json new-genesis.json --json-minimize=true
