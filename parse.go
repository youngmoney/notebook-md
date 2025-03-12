package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func parseLineType(l string) LineType {
	if l == kCodeTick {
		return CodeTickLine
	}
	if strings.HasPrefix(l, kCodeTick) {
		return CodeTickAndLine
	}
	if l == kOutputBlockStart {
		return OutputStartLine
	}
	if l == kOutputBlockEnd {
		return OutputEndLine
	}
	return TextLine
}

type AddAction int64

const (
	Add AddAction = iota + 1
	Close
	AddAndClose
	Bad
)

func newBlock(t LineType) Block {
	switch t {
	case TextLine:
		return &TextBlock{}
	case CodeTickAndLine:
		return &CodeBlock{}
	case CodeTickLine:
		return &CodeBlock{}
	case OutputStartLine:
		return &OutputBlock{}
	case OutputEndLine:
		return nil
	}
	return nil
}

func parseBlocks(ls *[]Line) ([]Block, error) {
	var blocks []Block
	blocks = append(blocks, &TextBlock{})

	for _, l := range *ls {
		action := blocks[len(blocks)-1].addLine(&l)
		switch action {
		case Add:
			continue
		case AddAndClose:
			blocks = append(blocks, &TextBlock{})
			continue
		case Close:
		case Bad:
			return nil, errors.New(fmt.Sprint("line ", l.Number, " bad line type ", l.Type, " while creating block ", reflect.TypeOf(blocks[len(blocks)-1])))
		}

		// Close action must add this line to the next block
		nb := newBlock(l.Type)
		if nb == nil {
			return nil, errors.New(fmt.Sprint("line ", l.Number, " bad block start type ", l.Type))
		}
		blocks = append(blocks, nb)

		newAction := blocks[len(blocks)-1].addLine(&l)
		if newAction != Add {
			return nil, errors.New(fmt.Sprint("line ", l.Number, " bad block first line ", l.Type, " while creating block ", reflect.TypeOf(blocks[len(blocks)-1])))
		}
	}

	var keep []Block
	for _, b := range blocks {
		if b.Empty() {
			continue
		}
		keep = append(keep, b)
	}

	return keep, nil

}

func parseGroups(blocks []Block) ([]Group, error) {
	var groups []Group
	groups = append(groups, &TextGroup{})
	for _, b := range blocks {
		action := groups[len(groups)-1].add(b)
		switch action {
		case Add:
			continue
		case AddAndClose:
			groups = append(groups, &TextGroup{})
			continue
		case Close:
		case Bad:
			return nil, errors.New("bad block type") // + l.Type + " while creating block " + c.Type)
		}

		switch b.(type) {
		case *TextBlock:
			groups = append(groups, &TextGroup{})
		case *CodeBlock:
			groups = append(groups, &ExecutionGroup{})
		case *OutputBlock:
			groups = append(groups, &TextGroup{})
		}

		newAction := groups[len(groups)-1].add(b)
		if newAction != Add {
			return nil, errors.New("bad group first block")
		}
	}

	var keep []Group
	for _, g := range groups {
		if g.String() == "" {
			continue
		}
		keep = append(keep, g)
	}

	return keep, nil
}

func parse(lines *[]string) ([]Group, error) {
	var ls []Line
	for n, l := range *lines {
		nl := Line{Type: parseLineType(l), Content: l, Number: n + 1}
		ls = append(ls, nl)
	}

	blocks, err := parseBlocks(&ls)
	if err != nil {
		return nil, err
	}

	return parseGroups(blocks)
}
