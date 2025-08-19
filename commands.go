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
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Cleared all users and feeds from the database")
	return nil
}

func handlerListUsers(s *state, _ command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println("URL is required")
		os.Exit(1)
	}

	feed, err := rss.FetchFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) < 2 {
		fmt.Println("Name and url for feed are required")
		os.Exit(1)
	}

	userInfo, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	feedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    userInfo.ID,
	}
	_, err = s.db.AddFeed(context.Background(), feedParams)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("New feed for %s was added\n", cmd.arguments[0])
	fmt.Printf("ID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\nURL: %v\nUser ID: %v\n", feedParams.ID, feedParams.CreatedAt, feedParams.UpdatedAt, feedParams.Name, feedParams.Url, feedParams.UserID)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.name](s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
