#!/usr/bin/env bash
set -e
set -x
. /etc/profile

PACKAGES=$(go list ./... | grep -v '/simulation')

go test -mod=readonly $(PACKAGES)