#!/usr/bin/env bash

function run() {
	cat tests/"$1".md | go run . --config tests/"$1".config.yaml execute
}

run simple
