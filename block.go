package main

import (
	"strings"
	"time"
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
	if len(b.Command) > 0 {
		return []string{kCodeTick + " " + strings.Join(b.Command, " ")}
	}
	return []string{kCodeTick}
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

const (
	kQuotePrefix = ">        "
)

func OutputBlockFromResult(output string, err error, style DisplayStyle) OutputBlock {
	var b = OutputBlock{Lines: strings.Split(output, "\n")}
	if err != nil {
		m := strings.Split(err.Error(), "\n")
		b.Lines = append(append(b.Lines, "failed with:"), m...)
	}

	if style == QUOTE {
		var nl []string
		for _, l := range b.Lines {
			nl = append(nl, kQuotePrefix+l)
		}
		b.Lines = nl
	}

	return b
}

func (b *OutputBlock) addLine(line *Line) AddAction {
	switch line.Type {
	case TextLine:
		if !strings.HasPrefix(line.Content, kOutputPrefix) {
			b.Lines = append(b.Lines, line.Content)
			return Add
		} else if t, err := time.Parse(kOutputModifiedTimeLineFormat, line.Content); err == nil {
			b.Modified = t
			return Add
		} else {
			return Bad
		}
	case CodeTickAndLine:
		return Bad
	case CodeTickLine:
		return Bad
	case OutputStartLine:
		if len(b.Lines) == 0 {
			return Add
		}
		return Bad
	case OutputEndLine:
		return AddAndClose
	}
	return Bad
}

func (b OutputBlock) Header() []string {
	var lines []string
	lines = append(lines, kOutputBlockStart)
	if !b.Modified.IsZero() {
		lines = append(lines, b.Modified.Format(kOutputModifiedTimeLineFormat))
	}
	if b.Lines[0] != "" {
		lines = append(lines, "")
	}
	return lines
}

func (b OutputBlock) Content() []string {
	return b.Lines
}

func (b OutputBlock) Footer() []string {
	if b.Lines[len(b.Lines)-1] == "" {
		return []string{kOutputBlockEnd}
	}
	return []string{"", kOutputBlockEnd}
}

func (b OutputBlock) Empty() bool {
	return len(b.Lines) == 0
}

func (b OutputBlock) String() string {
	return strings.Join(append(append(b.Header(), b.Content()...), b.Footer()...), "\n")
}
