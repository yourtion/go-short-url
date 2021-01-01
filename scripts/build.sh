#!/usr/bin/env bash

export GO111MODULE=on

cd $(dirname "$0")
set -e

cd ..

# 构建时候自动更新构建时间
updateBuildDate() {
  file="internal/base/define/service.go"
  sed -in-place -e s/build-[1,2][0-9]*/build-`date +%Y%m%d%H%M`/g ${file}
}

build() {
    rm -rf release/${1}
    mkdir -p release/${1}
    GOOS=${1} GOARCH=amd64 go build -v -o release/${1}/short-url cmd/short-url/main.go
    cp configs/config_template.toml release/${1}/config.toml
}

updateBuildDate

build darwin
build linux
#build windows
