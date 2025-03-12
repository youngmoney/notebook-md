package main

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
