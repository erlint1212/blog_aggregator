# Blog Aggregator 

Boot.Dev guided project, [link to project page](https://www.boot.dev/courses/build-blog-aggregator-golang)

RSS feed aggragor in Go

Build a blog aggregator microservice in Go. Put your API, database, and web scraping skills to the test.

## Use case

* Add RSS feeds from across the internet to be collected
* Store the collected posts in a PostgreSQL database
* Follow and unfollow RSS feeds that other users have added
* View summaries of the aggregated posts in the terminal, with a link to the full post

## Learning goals

* Learn how to integrate a Go application with a PostgreSQL database
* Practice using your SQL skills to query and migrate a database (using sqlc and goose, two lightweight tools for typesafe SQL in Go)
* Learn how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database

## Make proper nix env later

```
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go get github.com/google/uuid
go get github.com/lib/pq
```
Important tools used

* PostgreSQL
* sqlc (generate go code from sql query files)

## TODO

Config file is not properly isolated, add a proper Mutex to it so read/write
operations don't get screwed in the future.
