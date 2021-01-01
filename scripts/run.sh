#!/usr/bin/env bash

cd $(dirname "$0")
set -e

go run ../cmd/short-url/main.go
