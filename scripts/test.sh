#!/usr/bin/env bash

function run() {
	cat tests/"$1".md | go run . --config tests/"$1".config.yaml execute
}

function compare() {
	diff <(cat "$1" | sed 's/<!-- notebook output modified .* -->/<!-- modified -->/') <(cat "$2" | sed 's/<!-- notebook output modified .* -->/<!-- modified -->/')
}

compare <(run simple) <(cat tests/simple.executed.md)
