package config

import (
	"fmt"
)

func (c *Commands) Run(s *State, cmd Command) error {
	com, ok := c.CommandMap[cmd.Name]
	if !ok {
		err := fmt.Errorf("This command (%v) doesn't exist.", cmd.Name)
		return err
	}
	if err := com(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, handlerFunc func(*State, Command) error) {
	c.CommandMap[name] = handlerFunc
}
