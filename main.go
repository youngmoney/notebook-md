package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func readLines() []string {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return lines
}

func errorOutput(s string) *OutputBlock {
	return &OutputBlock{Lines: []string{s}, Modified: time.Now()}
}

func executeGroup(i *ExecutionGroup, config *Config) *ExecutionGroup {
	var o = ExecutionGroup{Code: i.Code}

	var name = config.Notebook.Default.CommandName
	var args []string
	if len(i.Code.Command) > 0 {
		name = i.Code.Command[0]
		args = i.Code.Command[1:]
	}
	if name == "" {
		o.Output = errorOutput("no command specified")
		return &o
	}
	m := Match(name, &config.Notebook.Commands)
	if m == nil {
		o.Output = errorOutput("command not found: " + name)
		return &o
	}
	out, err := ExecuteCommandCapture(m.Command, args, i.Code.CommandBody())
	n := OutputBlockFromResult(strings.Trim(out, "\n"), err)
	n.Modified = time.Now()
	o.Output = &n

	return &o
}

func commandExecute(config *Config) error {
	lines := readLines()
	groups, err := parse(&lines)
	if err != nil {
		return err
	}
	var out []Group
	for _, g := range groups {
		switch g.(type) {
		case *TextGroup:
			out = append(out, g)
		case *ExecutionGroup:
			out = append(out, executeGroup(g.(*ExecutionGroup), config))
		}
	}
	for _, g := range out {
		fmt.Println(g)
	}
	return nil
}

func commandExpand(config *Config) error {
	return nil
}

func main() {
	configFilename := flag.String("config", os.Getenv("NOTEBOOK_MD_CONFIG"), "config file (yaml), or set NOTEBOOK_MD_CONFIG")
	flag.Parse()

	config := ReadConfig(*configFilename)

	switch flag.Arg(0) {
	case "execute":
		fs := flag.NewFlagSet("execute", flag.ExitOnError)
		fs.Parse(flag.Args()[1:])
		ExitIfNonZero(commandExecute(&config))
	case "expand":
		fs := flag.NewFlagSet("expand", flag.ExitOnError)
		fs.Parse(flag.Args()[1:])
		ExitIfNonZero(commandExpand(&config))
	default:
		if flag.NArg() > 0 {
			fmt.Println("unknown command:", flag.Arg(0))
		}
		fmt.Println("supported commands: execute, expand")
		os.Exit(1)

	}
}
