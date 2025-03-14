package main

type Group interface {
	add(block Block) AddAction
	String() string
	Empty() bool
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

func (g TextGroup) Empty() bool {
	return len(g.Blocks) == 0
}

type ExecutionGroup struct {
	Code              *CodeBlock
	Output            *OutputBlock
	ConsumedEmptyLine bool
	StartedWithOutput bool
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
			g.StartedWithOutput = true
			return AddAndClose
		}
		return Close
	case *TextBlock:
		if g.Code != nil && g.Output == nil && block.String() == "" {
			g.ConsumedEmptyLine = true
			return Add
		}
		return Close
	}
	return Bad
}

func (g ExecutionGroup) LineRange() Range {
	if g.Code == nil {
		return Range{Lower: 0, Upper: 0}
	}
	return g.Code.LineRange
}

func (g ExecutionGroup) String() string {
	var extra string
	if g.ConsumedEmptyLine && !g.StartedWithOutput {
		extra = "\n"
	}
	if g.Code == nil {
		return ""
	}
	if g.Output == nil {
		return g.Code.String() + extra
	}
	return g.Code.String() + "\n\n" + g.Output.String() + extra
}

func (g ExecutionGroup) Empty() bool {
	return false
}

func (g ExecutionGroup) Expand(e Expand) string {
	var extra string
	if g.ConsumedEmptyLine && !g.StartedWithOutput {
		extra = "\n"
	}
	if g.Code == nil {
		return ""
	}
	if g.Output == nil {
		return g.Code.Expand(e) + extra
	}
	return g.Code.Expand(e) + "\n\n" + g.Output.String() + extra
}
