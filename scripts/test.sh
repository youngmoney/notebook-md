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

function line() {
	compare <(cat tests/line.md | go run . --config tests/line.config.yaml execute --line="$1") <(cat tests/line.executed."$1".md)
}

line 4
line 7
line -6
line 10-
line -
line 4,14
line 1,9-13
line 4,-
