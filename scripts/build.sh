#!/usr/bin/env bash

set -e
set -u
set -o pipefail
set -x

go tool templ generate
go build -o ./tmp/main cmd/server/main.go
