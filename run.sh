#!/bin/bash

set -o errexit
set -o pipefail

go build main.go

dapr run --app-id viewer --app-port 8083 --components-path ./components --config ./components/config.yaml --log-level debug -- go run main.go
