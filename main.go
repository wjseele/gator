package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/wjseele/gator/internal/config"
	"github.com/wjseele/gator/internal/database"
)

func main() {
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dbConfig.DB_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	configState := state{
		db:  dbQueries,
		cfg: &dbConfig,
	}

	commands := commands{
		commands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerFetcher)
	commands.register("addfeed", handlerAddFeed)
	commands.register("feeds", handlerListFeeds)

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
