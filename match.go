package main

import ()

func Match(name string, commands *[]Command) *Command {
	for _, m := range *commands {
		if m.Name == name {
			return &m
		}
	}
	return nil
}
