#!/usr/bin/env bash

function execute() {
	cat tests/"$1".md | go run . --config tests/"$1".config.yaml execute
}

function expand() {
	cat tests/"$1".md | go run . --config tests/"$1".config.yaml expand
}

function compare() {
	diff <(cat "$1" | sed 's/<!-- notebook output modified .* -->/<!-- modified -->/') <(cat "$2" | sed 's/<!-- notebook output modified .* -->/<!-- modified -->/')
}

compare <(execute simple) <(cat tests/simple.executed.md)
compare <(expand simple) <(cat tests/simple.expanded.md)

compare <(execute expand) <(cat tests/expand.executed.md)
compare <(expand expand) <(cat tests/expand.expanded.md)
