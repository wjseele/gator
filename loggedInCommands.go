package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/wjseele/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		userID, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return err
		}
		return handler(s, cmd, userID)
	}
}

func handlerAddFeed(s *state, cmd command, userID database.User) error {
	if len(cmd.arguments) < 2 {
		fmt.Println("Name and url for feed are required")
		os.Exit(1)
	}

	feedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    userID.ID,
	}
	_, err := s.db.AddFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID.ID,
		FeedID:    feedParams.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("New feed for %s was added\n", cmd.arguments[0])
	fmt.Printf("ID: %v\nCreated at: %v\nUpdated at: %v\nName: %v\nURL: %v\nUser ID: %v\n", feedParams.ID, feedParams.CreatedAt, feedParams.UpdatedAt, feedParams.Name, feedParams.Url, feedParams.UserID)
	return nil
}

func handlerFollow(s *state, cmd command, userID database.User) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("The URL is required.")
		os.Exit(1)
	}
	feedID, err := s.db.GetFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID.ID,
		FeedID:    feedID,
	}
	newFollowedFeed, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("User %s is now following feed %s\n", newFollowedFeed.UserName, newFollowedFeed.FeedName)
	return nil
}

func handlerUnfollow(s *state, cmd command, userID database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Please provide an URL")
	}
	feedID, err := s.db.GetFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}
	deleteParams := database.DeleteFeedParams{
		UserID: userID.ID,
		FeedID: feedID,
	}
	err = s.db.DeleteFeed(context.Background(), deleteParams)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully unfollowed %s", cmd.arguments[0])
	return nil
}

func handlerFollowing(s *state, _ command, userID database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), userID.ID)
	if err != nil {
		return err
	}
	for i := range followedFeeds {
		fmt.Printf("Feed %s is followed by %s\n", followedFeeds[i].FeedName, followedFeeds[i].UserName)
	}
	return nil
}

func handlerBrowse(s *state, cmd command, userID database.User) error {
	limit := 2
	if len(cmd.arguments) != 0 {
		i, err := strconv.Atoi(cmd.arguments[0])
		if err != nil {
			return err
		}
		limit = i
	}

	getPostParams := database.GetPostsForUserParams{
		UserID: userID.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		return err
	}

	for i := range posts {
		fmt.Printf("Title: %s\n", posts[i].Title)
		fmt.Printf("URL: %s\n", posts[i].Url)
		fmt.Printf("Content: %s\n", posts[i].Description)
	}

	return nil
}
