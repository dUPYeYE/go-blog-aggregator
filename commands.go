package main

import "fmt"

type Command struct {
	name string
	args []string
}

type Commands struct {
	cmds map[string]func(*State, Command) error
}

func (c *Commands) register(name string, handler func(*State, Command) error) {
	c.cmds[name] = handler
}

func (c *Commands) run(s *State, cmd Command) error {
	if handler, ok := c.cmds[cmd.name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("Command %s not found", cmd.name)
}
