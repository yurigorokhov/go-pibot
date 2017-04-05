#!/bin/sh

export GOOS=linux
export GOARCH=arm

echo "[INFO] Fetching dependencies"
go get -d -v

echo "[INFO] Building"
go build -v github.com/yurigorokhov/go-pibot

echo "[INFO] Creating Package"
mkdir -p build
tar -c go-pibot public scripts > build/package.tar && gzip -f build/package.tar

echo "[INFO] Uploading Package"
s3cmd put build/package.tar.gz s3://batbot/package/package.tar.gz


