#!/usr/bin/env bash

GOOS=darwin GOARCH=arm64 go build -o bin/notebook-md-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -o bin/notebook-md-linux-amd64 .
