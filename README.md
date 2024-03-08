# SEKAI

SEKAI is KIRA Network's base layer (L1) blockchain application sometimes referred to as “backend”. The role of SEKAI is to be a source of shared security as well as a governance and settlement layer for all KIRA RollApps (L2). KIRA Blockchain preserves information such as user account balances, governance permissions, and RollApp state roots as well as other essential data for coordinating both L1 and L2 operations.

## Documentation

For the most up to date documentation please visit [docs.kira.network](https://docs.kira.network/)

## Quick Installation Guide

### Required Tools

KIRA requires the installation of two tools, [Cosign](https://docs.kira.network/docs/cosign) and [Bash-utils](https://docs.kira.network/docs/bash-utils), in order to secure the network and simplify the execution of various tasks.
All files in KIRA repositories are always signed with cosign, you should NEVER install anything on your machine unless you verified integrity of the files!

### Installation

Login as admin & load environment variables.

sudo -s

Set desired SEKAI release version and binaries repo as env variables within `/etc/profile` (with `bash-utils` or manually). Sourcing `/etc/profile` is necessary.
Check latest SEKAI release's version [here](https://github.com/KiraCore/sekai/releases).

```bash
setGlobEnv SEKAI_VERSION "v0.3.39" && \
setGlobEnv SEKAI_REPO "$HOME/sekai" && \
setGlobEnv NETWORK_NAME "test" && \
setGlobEnv SEKAID_HOME "~/.sekaid-$NETWORK_NAME" && \
loadGlobEnvs
```

Clone repository and install

```bash
rm -rf $SEKAI_REPO && rm -fr $GOBIN/sekaid && mkdir $SEKAI_REPO && cd $SEKAI_REPO && \
git clone https://github.com/KiraCore/sekai.git -b $SEKAI_VERSION $SEKAI_REPO && \
chmod -R 777 ./scripts && make install && \
echo "SUCCESS installed sekaid $(sekaid version)" || echo "FAILED"
```

Verify successful installation

```bash
sekaid version --long
```

## SEKAI Modules

### Consensus

- [staking](https://github.com/KiraCore/sekai/tree/master/x/staking)
- [slashing](https://github.com/KiraCore/sekai/tree/master/x/slashing)
- [evidence](https://github.com/KiraCore/sekai/tree/master/x/evidence)
- [distributor](https://github.com/KiraCore/sekai/tree/master/x/distributor)

### Basic modules

- [spending](https://github.com/KiraCore/sekai/tree/master/x/spending)
- [tokens](https://github.com/KiraCore/sekai/tree/master/x/tokens)
- [ubi](https://github.com/KiraCore/sekai/tree/master/x/ubi)

### Liquid Staking

- [multistaking](https://github.com/KiraCore/sekai/tree/master/x/multistaking)

### Derivatives

- [basket](https://github.com/KiraCore/sekai/tree/master/x/basket)
- [collectives](https://github.com/KiraCore/sekai/tree/master/x/collectives)

### Governance

- [gov](https://github.com/KiraCore/sekai/tree/master/x/gov)

### Layer2

- [layer2](https://github.com/KiraCore/sekai/tree/master/x/layer2)

### Fees

- [feeprocessing](https://github.com/KiraCore/sekai/tree/master/x/feeprocessing)

### Utilities & Upgrade

- [custody](https://github.com/KiraCore/sekai/tree/master/x/custody)
- [recovery](https://github.com/KiraCore/sekai/tree/master/x/recovery)
- [genutil](https://github.com/KiraCore/sekai/tree/master/x/genutil)
- [upgrade](https://github.com/KiraCore/sekai/tree/master/x/upgrade)

## Contributing

Check out [contributing.md](./CONTRIBUTING.md) for our guidelines & policies for how we develop the Kira chain. Thank you to all those who have contributed!
