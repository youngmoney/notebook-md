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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return lines
}

func errorOutput(s string) *OutputBlock {
	return &OutputBlock{Lines: []string{"Error: " + s}, Modified: time.Now()}
}

func executeGroup(i *ExecutionGroup, config *Config) *ExecutionGroup {
	o := *i
	o.Output = nil

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
	n := OutputBlockFromResult(strings.Trim(out, "\n"), err, m.DisplayStyle)
	n.Modified = time.Now()
	o.Output = &n

	return &o
}

func commandExecute(config *Config, r MultiRange) error {
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
			if r.Overlaps(g.(*ExecutionGroup).LineRange()) {
				out = append(out, executeGroup(g.(*ExecutionGroup), config))
			} else {
				out = append(out, g)
			}
		}
	}
	for i, g := range out {
		if i == len(out)-1 {
			fmt.Println(strings.TrimSuffix(g.String(), "\n"))
		} else {
			fmt.Println(g)
		}
	}
	return nil
}

func expand(i *ExecutionGroup, config *Config) string {
	var name = config.Notebook.Default.CommandName
	if len(i.Code.Command) > 0 {
		name = i.Code.Command[0]
	}
	if name == "" {
		return i.Expand(Expand{})
	}
	m := Match(name, &config.Notebook.Commands)
	if m == nil {
		return i.Expand(Expand{})
	}
	return i.Expand(m.Expand)
}

func commandExpand(config *Config) error {
	lines := readLines()
	groups, err := parse(&lines)
	if err != nil {
		return err
	}
	for i, g := range groups {
		var out string
		switch g.(type) {
		case *TextGroup:
			out = g.String()
		case *ExecutionGroup:
			out = expand(g.(*ExecutionGroup), config)
		}
		if i == len(out)-1 {
			fmt.Println(strings.TrimSuffix(out, "\n"))
		} else {
			fmt.Println(out)
		}
	}
	return nil
}

func main() {
	configFilename := flag.String("config", os.Getenv("NOTEBOOK_MD_CONFIG"), "config file (yaml), or set NOTEBOOK_MD_CONFIG")
	flag.Parse()

	config := ReadConfig(*configFilename)

	switch flag.Arg(0) {
	case "execute":
		fs := flag.NewFlagSet("execute", flag.ExitOnError)
		lines := fs.String("line", "", "lines to execute <line>|<start line>-<end line>|-<end-line>|<start-line>- comma seperated")
		fs.Parse(flag.Args()[1:])
		lineRange, lerr := MultiRangeFromString(*lines)
		if lerr != nil {
			fmt.Fprintln(os.Stderr, lerr)
			os.Exit(1)
		}
		err := commandExecute(&config, lineRange)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "expand":
		fs := flag.NewFlagSet("expand", flag.ExitOnError)
		fs.Parse(flag.Args()[1:])
		ExitIfNonZero(commandExpand(&config))
	default:
		if flag.NArg() > 0 {
			fmt.Fprintln(os.Stderr, "unknown command:", flag.Arg(0))
		}
		fmt.Fprintln(os.Stderr, "supported commands: execute, expand")
		os.Exit(1)

	}
}
