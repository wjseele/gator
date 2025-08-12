package main

import (
	"fmt"
	"os"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Username is required")
		os.Exit(1)
	}
	err := s.config.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s\n", cmd.arguments[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.name](s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
