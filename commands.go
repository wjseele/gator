package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/wjseele/gator/internal/database"
	"github.com/wjseele/gator/internal/rss"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Username is required")
		os.Exit(1)
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s\n", cmd.arguments[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("Name for new user is required")
		os.Exit(1)
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}
	_, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	fmt.Printf("New user %s created\n", cmd.arguments[0])

	err = s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s\n", cmd.arguments[0])
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.ClearDB(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Cleared all users and feeds from the database")
	return nil
}

func handlerListUsers(s *state, _ command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return err
	}
	for _, name := range users {
		if name == s.cfg.CurrentUser {
			fmt.Printf("* %s (current)\n", name)
		} else {
			fmt.Printf("* %s\n", name)
		}
	}
	return nil
}

func handlerFetcher(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Specify time in duration/unit (1s, 2h, etc)")
	}

	time_between_reqs, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerListFeeds(s *state, _ command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Here are the feeds currently in the database:")
	for iter := range feeds {
		fmt.Printf("Name: %s, URL: %s, ", feeds[iter].Name, feeds[iter].Url)
		userName, err := s.db.GetUserByID(context.Background(), feeds[iter].UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Owner: %s\n", userName.Name)
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.name](s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
