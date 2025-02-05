package main

import (
    "github.com/erlint1212/blog_aggregator/internal/database"
    "database/sql"
    "fmt"
    "time"
    "context"
    "html"
	"github.com/google/uuid"
    "github.com/lib/pq"
)

func scrapeFeeds(s *state) error {

    ctx := context.Background()

    next_feed, err := s.db.GetNextFeedToFetch(ctx)
    if err != nil {
        return err
    }
    feed_params := database.MarkFeedFetchedParams{
        sql.NullTime{time.Now(), true},
        next_feed.ID,
    }

    err = s.db.MarkFeedFetched(ctx, feed_params)
    if err != nil {
        return err
    }

    RSSFeed, err := fetchFeed(ctx, next_feed.Url)
    if err != nil {
        return err
    }

    const ColNotUniqueErr = "23505"

    RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)
    RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)
    for i := 0; i < len(RSSFeed.Channel.Item); i++ {
        RSSFeed.Channel.Item[i].Title = html.UnescapeString(RSSFeed.Channel.Item[i].Title)
        RSSFeed.Channel.Item[i].Description = html.UnescapeString(RSSFeed.Channel.Item[i].Description)

        err_m := fmt.Errorf("Item %d failed for RSSFeed \"%s\":", i, RSSFeed.Channel.Title)

        published_at, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", RSSFeed.Channel.Item[i].PubDate)
        if err != nil {
            return fmt.Errorf("%s %s", err_m, err)
        }
    
        new_post := database.CreatePostParams{
            uuid.New(),
            time.Now(),
            time.Now(),
            RSSFeed.Channel.Item[i].Title,
            RSSFeed.Channel.Item[i].Link,
            sql.NullString{RSSFeed.Channel.Item[i].Description, true},
            published_at,
            next_feed.ID,
        }

        _, err = s.db.CreatePost(ctx, new_post)
        if err != nil {
            return_err := true
            // Check if postgres error, then check if the error is non uniqueness(23505)
            if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == ColNotUniqueErr {
                return_err = false
            }

            if return_err {
                return fmt.Errorf("%s %s", err_m, err)
            }
        }
    }


    return nil
}
