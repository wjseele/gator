package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wjseele/gator/internal/config"
	"os"
)

func main() {
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	var configState state
	configState.config = &dbConfig

	commands := commands{
		commands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("No commands given.")
		os.Exit(1)
	}
	command := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}
	err = commands.run(&configState, command)
	if err != nil {
		fmt.Println(err)
	}
}
