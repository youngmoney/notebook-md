package main

const (
	kCodeTick                     = "```"
	kOutputBlockStart             = "<!-- notebook output start -->"
	kOutputBlockEnd               = "<!-- notebook output end -->"
	kOutputPrefix                 = "<!-- notebook output "
	kOutputModifiedPrefix         = "<!-- notebook output modified "
	kOutputModifiedSuffix         = " -->"
	kOutputModifiedTimeFormat     = "2006-01-02T15:04:05"
	kOutputModifiedTimeLineFormat = kOutputModifiedPrefix + kOutputModifiedTimeFormat + kOutputModifiedSuffix
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
