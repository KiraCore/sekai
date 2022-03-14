#!/usr/bin/env bash
set -e
set -x
. /etc/profile

go test -mod=readonly $(go list ./... | grep -v '/simulation')