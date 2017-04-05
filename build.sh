#!/bin/sh

export GOOS=linux
export GOARCH=arm

echo "[INFO] Fetching dependencies"
go get -d -v

echo "[INFO] Building"
go build -v github.com/yurigorokhov/go-pibot
