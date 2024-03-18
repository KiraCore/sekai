# Contributing

- [Contributing](#contributing)
  - [Overview](#overview)
  - [Kira Improvement Proposal (KIP)](#kira-improvement-proposal-kip)
    - [Testing](#testing)
    - [Pull Requests](#pull-requests)
    - [Requesting Reviews](#requesting-reviews)
    - [Changelog](#changelog)
  - [Dependencies](#dependencies)
  - [Protobuf](#protobuf)
  - [Branching Model and Release](#branching-model-and-release)

## Overview

Contributing to this repo can mean many things such as participating in
discussion or proposing code changes.
Following the processes outlined in this document will lead to the best
chance of getting changes merged into the codebase.

## Kira Improvement Proposal (KIP)

When proposing an architecture decision for Gaia, please start by opening an [issue](https://github.com/cosmos/gaia/issues/new/choose) or a [discussion](https://github.com/cosmos/gaia/discussions/new) with a summary of the proposal. Once the proposal has been discussed and there is rough alignment on a high-level approach to the design, you may either start development, or write an ADR.

If your architecture decision is a simple change, you may contribute directly without writing an ADR. However, if you are proposing a significant change, please include a corresponding ADR.

To create an ADR, follow the [template](./docs/architecture/adr-template.md) and [doc](./docs/architecture/README.md). If you would like to see examples of how these are written, please refer to the current [ADRs](https://github.com/cosmos/gaia/tree/main/docs/architecture).

### Testing

Tests can be executed by running `make test` at the top level of the Sekai repository.

### Pull Requests

Before submitting a pull request:

- synchronize your branch with the latest `master` branch and resolve any arising conflicts, `git fetch origin/master && git merge origin/master`
- run `make install`, `make test`, to ensure that all checks and tests pass.

### Requesting Reviews

In order to accommodate the review process, the author of the PR should be in contact with Sekai repo maintainers.

### Changelog

Changelog keeps the changes made as part of releases. The logs are kept on [CHANGELOG](./CHANGELOG.md).

## Dependencies

We use [Go Modules](https://github.com/golang/go/wiki/Modules) to manage
dependency versions.

The main branch of every Cosmos repository should just build with `go get`,
which means they should be kept up-to-date with their dependencies so we can
get away with telling people they can just `go get` our software.

When dependencies in Sekai's `go.mod` are changed, it is generally accepted practice
to delete `go.sum` and then run `go mod tidy`.

Since some dependencies are not under our control, a third party may break our
build, in which case we can fall back on `go mod tidy -v`.

## Protobuf

We use [Protocol Buffers](https://developers.google.com/protocol-buffers) along with [gogoproto](https://github.com/cosmos/gogoproto) to generate code for use in sekai.

For deterministic behavior around Protobuf tooling, everything is containerized using Docker. Make sure to have Docker installed on your machine, or head to [Docker's website](https://docs.docker.com/get-docker/) to install it.

To generate the protobuf stubs, you can run `make proto-gen`.

## Branching Model and Release

Sekai branches should be one of `feature/{feature-description}` or `bugfix/{bugfix-description}` to join CI/CD process.

Sekai follows [semantic versioning](https://semver.org).

To release a new version

- Set a new sekai version on [`types/constants.go`](types/constants.go)
- Add relevant information on [`RELEASE.md`](RELEASE.md)
- Push the code to the branch
- The bot automatically creates release branch `release/vx.x.x` as configured in [`types/constants.go`](types/constants.go) and raise the PR from working branch to release branch automatically
- Check CI/CD pass on the PR
- Get manual review
- Get the PR merged into release branch
- New PR into master is raised after release branch PR merge
- Check CI/CD pass on the PR
- Merge into master
