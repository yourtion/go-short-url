#!/usr/bin/env bash

cd $(dirname "$0")
set -e

go run short-url/src
