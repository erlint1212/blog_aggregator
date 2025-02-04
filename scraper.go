package main

import (
    "github.com/erlint1212/blog_aggregator/internal/database"
    "database/sql"
    "fmt"
    "time"
    "context"
    "html"
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


    RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)
    RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)
    for i := 0; i < len(RSSFeed.Channel.Item); i++ {
        RSSFeed.Channel.Item[i].Title = html.UnescapeString(RSSFeed.Channel.Item[i].Title)
        RSSFeed.Channel.Item[i].Description = html.UnescapeString(RSSFeed.Channel.Item[i].Description)
        fmt.Printf("%s\n", RSSFeed.Channel.Item[i].Title)
    }

    return nil
}
