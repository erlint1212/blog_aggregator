
package main

import (
    "context"
    "net/http"
    "io"
    "encoding/xml"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        return &RSSFeed{}, err
    }

    req.Header.Set("User-Agent", "gator")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return &RSSFeed{}, err
    }
    defer resp.Body.Close()


    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return &RSSFeed{}, err
    }

    var RSS RSSFeed
    if err = xml.Unmarshal(data, &RSS); err != nil {
		return &RSSFeed{}, err
	}

    return &RSS, nil
}
