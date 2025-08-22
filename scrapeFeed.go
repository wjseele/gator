package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wjseele/gator/internal/database"
	"github.com/wjseele/gator/internal/rss"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	fetchParams := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:            nextFeed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), fetchParams)
	if err != nil {
		return err
	}

	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	for post := range feed.Channel.Item {
		fmt.Printf("From %s: %s\n", feed.Channel.Title, feed.Channel.Item[post].Title)
	}
	return nil
}
