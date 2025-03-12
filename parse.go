package main

import (
	"errors"
	"strings"
	"time"
)

const (
	kCodeTick         = "```"
	kOutputBlockStart = "<!-- notebook output start -->"
	kOutputBlockEnd   = "<!-- notebook output end -->"
	kOutputPrefix     = "<!-- notebook output "
)

type LineType string

const (
	TextLine        LineType = "TextLine"
	CodeTickLine    LineType = "CodeTickLine"
	CodeTickAndLine LineType = "CodeTickAndLine"
	OutputStartLine LineType = "OutputStartLine"
	OutputEndLine   LineType = "OutputEndLine"
)

type Line struct {
	Number  int
	Type    LineType
	Content string
}

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

type Block interface {
	// Return true if the block is closed
	addLine(line *Line) AddAction
	Header() []string
	Content() []string
	Footer() []string
	Empty() bool
	String() string
}

type TextBlock struct {
	Lines []string
}

func (b *TextBlock) addLine(line *Line) AddAction {
	switch line.Type {
	case TextLine:
		b.Lines = append(b.Lines, line.Content)
		return Add
	case CodeTickAndLine:
		return Close
	case CodeTickLine:
		return Close
	case OutputStartLine:
		return Close
	case OutputEndLine:
		return Bad
	}
	return Bad
}

func (b TextBlock) Header() []string {
	return []string{}
}

func (b TextBlock) Content() []string {
	return b.Lines
}

func (b TextBlock) Footer() []string {
	return []string{}
}

func (b TextBlock) String() string {
	return strings.Join(append(append(b.Header(), b.Content()...), b.Footer()...), "\n")
}

func (b TextBlock) Empty() bool {
	return len(b.Lines) == 0
}

type CodeBlock struct {
	Lines   []string
	Command []string
}

func (b CodeBlock) CommandBody() string {
	return strings.Join(b.Lines, "\n")
}

func (b *CodeBlock) addLine(line *Line) AddAction {
	if len(b.Lines) == 0 && strings.HasPrefix(line.Content, kCodeTick) {
		// TODO: consider a shlex
		c := strings.TrimPrefix(strings.TrimPrefix(line.Content, kCodeTick), " ")
		b.Command = strings.Split(c, " ")
		return Add
	}
	switch line.Type {
	case TextLine:
		b.Lines = append(b.Lines, line.Content)
		return Add
	case CodeTickAndLine:
		return Bad
	case CodeTickLine:
		return AddAndClose
	case OutputStartLine:
		return Bad
	case OutputEndLine:
		return Bad
	}
	return Bad
}

func (b CodeBlock) Header() []string {
	return []string{kCodeTick + strings.Join(b.Command, " ")}
}

func (b CodeBlock) Content() []string {
	return b.Lines
}

func (b CodeBlock) Footer() []string {
	return []string{kCodeTick}
}

func (b CodeBlock) String() string {
	return strings.Join(append(append(b.Header(), b.Content()...), b.Footer()...), "\n")
}

func (b CodeBlock) Empty() bool {
	return false
}

type OutputBlock struct {
	Lines    []string
	Modified time.Time
}

func OutputBlockFromResult(output string, err error) OutputBlock {
	var b = OutputBlock{Lines: strings.Split(output, "\n")}
	if err != nil {
		m := strings.Split(err.Error(), "\n")
		b.Lines = append(append(b.Lines, "failed with:"), m...)
	}
	return b
}

func (b *OutputBlock) addLine(line *Line) AddAction {
	switch line.Type {
	case TextLine:
		if !strings.HasPrefix(line.Content, kOutputPrefix) {
			b.Lines = append(b.Lines, line.Content)
		}
		return Add
	case CodeTickAndLine:
		return Bad
	case CodeTickLine:
		return Bad
	case OutputStartLine:
		return Bad
	case OutputEndLine:
		return Close
	}
	return Bad
}

func (b OutputBlock) Header() []string {
	if b.Modified.IsZero() {
		return []string{kOutputBlockStart}
	}
	return []string{kOutputBlockStart, "<!-- notebook output modified " + b.Modified.Format("2006-01-02T15:04:05") + " -->"}

}

func (b OutputBlock) Content() []string {
	return b.Lines
}

func (b OutputBlock) Footer() []string {
	return []string{kOutputBlockEnd}
}

func (b OutputBlock) Empty() bool {
	return len(b.Lines) == 0
}

func (b OutputBlock) String() string {
	return strings.Join(append(append(b.Header(), b.Content()...), b.Footer()...), "\n")
}

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
			return nil, errors.New("bad line type ") // + l.Type + " while creating block " + c.Type)
		}

		// Close action must add this line to the next block
		nb := newBlock(l.Type)
		if nb == nil {
			return nil, errors.New("bad block start type ") // + l.Type + " while creating block " + c.Type)
		}
		blocks = append(blocks, nb)

		newAction := blocks[len(blocks)-1].addLine(&l)
		if newAction != Add {
			return nil, errors.New("bad block first line")
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

type Group interface {
	add(block Block) AddAction
	String() string
}

type TextGroup struct {
	Blocks []Block
}

func (g *TextGroup) add(block Block) AddAction {
	switch block.(type) {
	case *TextBlock:
		g.Blocks = append(g.Blocks, block)
		return Add
	case *CodeBlock:
		return Close
	case *OutputBlock:
		g.Blocks = append(g.Blocks, block)
		return Add
	}
	return Bad
}

func (g TextGroup) String() string {
	var s string
	for _, b := range g.Blocks {
		if s == "" {
			s = b.String()
		} else {
			s = s + "\n" + b.String()
		}
	}
	return s
}

type ExecutionGroup struct {
	Code   *CodeBlock
	Output *OutputBlock
}

func (g *ExecutionGroup) add(block Block) AddAction {
	switch block.(type) {
	case *CodeBlock:
		if g.Code == nil {
			g.Code = block.(*CodeBlock)
			return Add
		}
		return Close
	case *OutputBlock:
		if g.Code != nil && g.Output == nil {
			g.Output = block.(*OutputBlock)
			return AddAndClose
		}
		return Close
	case *TextBlock:
		if g.Code != nil && g.Output == nil && block.String() == "" {
			return Add
		}
		return Close
	}
	return Bad
}

func (g ExecutionGroup) String() string {
	if g.Code == nil {
		return ""
	}
	if g.Output == nil {
		return g.Code.String()
	}
	return g.Code.String() + "\n\n" + g.Output.String()
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
		nl := Line{Type: parseLineType(l), Content: l, Number: n}
		ls = append(ls, nl)
	}

	blocks, err := parseBlocks(&ls)
	if err != nil {
		return nil, err
	}

	return parseGroups(blocks)
}
