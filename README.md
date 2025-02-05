# Blog Aggregator - Gator


Boot.Dev guided project, [link to project page](https://www.boot.dev/courses/build-blog-aggregator-golang)

RSS feed aggragor in Go

Build a blog aggregator microservice in Go. Put your API, database, and web scraping skills to the test.

## Requirements
Neeed to have the following installed to run this code:

* Go (If not compiled, need to compile youself)
* PostgreSQL (or else the SQL stuff will fail)

Can install by writing this in console `go install github.com/erlint1212/blog_aggregator`

## Setup

Have to create own config file, doe this by either running the `reset_cfg.sh` file or manually typing `echo '{"db_url": "postgres://postgres:postgres@localhost:5432/gator"}' >| ~/.gatorconfig.json
` into console.

## Available commands

All commands start with `run . gator` or `gator` if installed with `go install`.

Some commands:

* `run . gator register <your name>` registers user by username, write your name in `<your name>`
* `run . gator reset` deletes all users, all their feeds and all their posts
* `run . gator users` prints all users into the console, with the currently logged in user having `(current)` behind their name.
* `run . gator agg <time interval>` initiates an infinite loop that gets all the posts from all the currently register feeds and stores them under feed creators name, updates the oldest first.
* `run . gator addfeed <feed title> <feed url>` add a name for the feed and the url to a RSS API to store the RSS for later use.
* `run . gator feeds` gets all the feeds title and username of the user that added it.
* `run . gator follow <feed URL>` currently logged in user will follow the feed if it is stored.
* `run . gator following` prints all the feeds that the currently logged in user is following.
* `run . gator unfollow <feed URL>` currently logged in user will unfollow the feed if it is stored and followed.
* `run . gator browse <amount of posts>` prints out the top `<amount of posts>` posts in followed feed, <amount of posts> is an optional argument, the default is 2.


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
