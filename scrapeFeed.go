package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

	for i := range feed.Channel.Item {
		parsedTime, err := time.Parse(time.RFC1123Z, feed.Channel.Item[i].PubDate)
		if err != nil {
			return err
		}
		post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       feed.Channel.Item[i].Title,
			Url:         feed.Channel.Item[i].Link,
			Description: feed.Channel.Item[i].Description,
			PublishedAt: sql.NullTime{Time: parsedTime},
			FeedID:      nextFeed.ID,
		}
		err = s.db.CreatePost(context.Background(), post)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					continue
				}
			} else {
				return err
			}
		}
	}

	fmt.Println("Added new posts to the database")

	return nil
}
