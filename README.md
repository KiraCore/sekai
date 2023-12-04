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

## Features

### Core modules for consensus

- staking
- slashing
- evidence
- distributor

### Basic modules

- spending
- tokens
- ubi

### Liquid staking

- multistaking

### Derivatives

- basket
- collectives

### Governance

- gov

### Layer2

- layer2

### Fees

- feeprocessing

### Utilities & Upgrade

- custody
- recovery
- genutil
- upgrade

## Contributing

Check out [contributing.md](./CONTRIBUTING.md) for our guidelines & policies for how we develop the Kira chain. Thank you to all those who have contributed!
