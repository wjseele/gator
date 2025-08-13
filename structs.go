package main

import (
	"github.com/wjseele/gator/internal/config"
	"github.com/wjseele/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands map[string]func(*state, command) error
}
