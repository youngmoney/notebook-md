package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type Default struct {
	CommandName string `yaml:"command_name"`
}

type ExpandStyle int64

const (
	NONE ExpandStyle = iota + 1 // leave as is
	HIDE                        // remove when expanding
	LINE                        // remove new lines from code block
	ONCE
	HEREDOC // replace with bash heredoc
)

func (t *ExpandStyle) UnmarshalYAML(n *yaml.Node) error {
	v := strings.ToLower(n.Value)
	switch v {
	case "none":
		*t = NONE
	case "hide":
		*t = HIDE
	case "line":
		*t = LINE
	case "once":
		*t = ONCE
	case "heredoc":
		*t = HEREDOC
	default:
		return errors.New("unkown expand style" + v)
	}
	return nil
}

type Expand struct {
	CommandName string      `yaml:"command_name"`
	BlockName   string      `yaml:"block_name"`
	Style       ExpandStyle `yaml:"style"`
}

type DisplayStyle int64

const (
	QUOTE DisplayStyle = iota + 1
	RAW
)

func (t *DisplayStyle) UnmarshalYAML(n *yaml.Node) error {
	v := strings.ToLower(n.Value)
	switch v {
	case "raw":
		*t = RAW
	case "quote":
		*t = QUOTE
	default:
		return errors.New("unkown display style" + v)
	}
	return nil
}

type Command struct {
	Name         string       `yaml:"name"`
	Command      string       `yaml:"command"`
	DisplayStyle DisplayStyle `yaml:"display_style"`
	Expand       Expand       `yaml:"expand"`
}

type Notebook struct {
	Commands []Command `yaml:"commands"`
	Default  Default   `yaml:"default"`
}

type Config struct {
	Notebook Notebook `yaml:"notebook"`
}

func ReadConfig(filename string) Config {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to read config: ", filename)
		os.Exit(1)
	}

	config := Config{}
	if err := yaml.Unmarshal(raw, &config); err != nil {
		fmt.Fprintln(os.Stderr, "unable to parse config: ", filename)
		os.Exit(1)
	}

	return config
}
