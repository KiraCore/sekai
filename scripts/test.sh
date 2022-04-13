#!/usr/bin/env bash
set -e
set -x
. /etc/profile

echo "INFO: Cleaning up system resources"
kill -9 $(lsof -t -i:46147) || echo "WARNING: Nothing running on port 9090, or failed to kill processes"

echo "INFO: Cleaning up system resources"
go test -mod=readonly $(go list ./... | grep -v '/simulation')