package main

import (
    "errors"
    "database/sql"
    "github.com/erlint1212/blog_aggregator/internal/config"
    "github.com/erlint1212/blog_aggregator/internal/database"
    "fmt"
    "time"
    "context"
	"github.com/google/uuid"
    "strconv"
)

type state struct {
    db      *database.Queries
    cfg     *config.Config
}

type command struct {
    name        string
    args        []string
}

type commands struct {
    handlers    map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
    f, ok := c.handlers[cmd.name]
    if !ok {
        return fmt.Errorf("Command dosen't exist")
    }

    err := f(s, cmd)
    if err != nil {
        return err
    }

    return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, cmd command) error {
        ctx := context.Background()
        current_user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
        if err != nil {
            return err
        }

        err = handler(s, cmd, current_user)
        if err != nil {
            return err
        }

        return nil
    }
    
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("the login handler expects a single argument, the username.")
    }

    username := cmd.args[0] 

    ctx := context.Background()
    _, err := s.db.GetUser(ctx, username)
    if err != nil {
        return err
    }

    err = s.cfg.SetUser(username)
    if err != nil {
        return err
    }

    fmt.Printf("User %s has been set\n", username)

    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("the login handler expects a single argument, the username.")
    }

    username := cmd.args[0]
    
    
    ctx := context.Background()
    user_params := database.CreateUserParams{
        uuid.New(),
        time.Now(),
        time.Now(),
        username,
    }

    user, err := s.db.GetUser(ctx, username)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return err
    }
    if user.Name == username {
        return fmt.Errorf("The user is already in the database")
    }

    _, err = s.db.CreateUser(ctx, user_params)
    if err != nil {
        return err
    }

    err = s.cfg.SetUser(username)
    if err != nil {
        return err
    }

    fmt.Printf("User %s has been created\n", username)
    fmt.Printf("%+v\n", user_params)

    return nil
}

func handlerReset(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return fmt.Errorf("There should be no arguments for reset command")
    }

    ctx := context.Background()
    err := s.db.DeleteAllUsers(ctx)
    if err != nil {
        return err
    }

    fmt.Printf("All rows in 'users' table have been deleted\n")
    return nil
}

func handlerUsers(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return fmt.Errorf("There should be no arguments for agg command")
    }

    ctx := context.Background()
    users, err := s.db.GetAllUsers(ctx)
    if err != nil {
        return err
    }

    current_loggedIn_user := s.cfg.CurrentUserName 

    for i := 0; i < len(users); i++ {
        special_status := ""
        if users[i].Name == current_loggedIn_user {
            special_status = " (current)"
        }
        fmt.Printf("* %s%s\n", users[i].Name, special_status)
    }

    return nil
}

func handlerAgg(s *state, cmd command) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("There should be 1 argument for agg command, time_between_reqs")
    }

    time_between_reqs, err := time.ParseDuration(cmd.args[0])
    if err != nil {
        return err
    }

    fmt.Printf("Collecting feeds every %s\n", time_between_reqs)

    ticker := time.NewTicker(time_between_reqs)
    for ; ; <-ticker.C {
        err = scrapeFeeds(s)
        if err != nil {
            return err
        }
    }

    return nil
}

func handlerAddFeed(s *state, cmd command, current_user database.User) error {
    if len(cmd.args) != 2 {
        return fmt.Errorf("the addfeed handler expects a two argument, name and url.")
    }

    name := cmd.args[0]
    url := cmd.args[1]

    ctx := context.Background()
    
    feed_params := database.CreateFeedParams{
        uuid.New(),
        time.Now(),
        time.Now(),
        name,
        url,
        current_user.ID,
    }

    feed, err := s.db.CreateFeed(ctx, feed_params)
    if err != nil {
        return err
    }
    
    cmd_follow_url := cmd
    cmd_follow_url.args = []string{feed.Url}
    err = middlewareLoggedIn(handlerFollow)(s, cmd_follow_url)
    if err != nil {
        return fmt.Errorf("Feed created but unable to follow feed: %s", err)
    }

    fmt.Printf("Feed \"%s\" has been created\n", name)
    fmt.Printf("%+v\n", feed_params)

    return nil
}

func handlerFeeds(s *state, cmd command) error {
    if len(cmd.args) != 0 {
        return fmt.Errorf("There should be no arguments for feeds command")
    }

    ctx := context.Background()
    feeds, err := s.db.GetAllFeeds(ctx)
    if err != nil {
        return err
    }

    for i := 0; i < len(feeds); i++ {
        creator_user, err := s.db.GetUserByID(ctx, feeds[i].UserID)
        if err != nil {
            return err
        }
        fmt.Printf("Name: %s URL: %s Username: %s\n", feeds[i].Name, feeds[i].Url, creator_user.Name)
    }

    return nil
}

func handlerFollow(s *state, cmd command, current_user database.User) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("There should be one arguments for follow command, url")
    }

    url := cmd.args[0]
    
    ctx := context.Background()

    feed, err := s.db.GetFeedByURL(ctx, url)
    if err != nil {
        return err
    }
    
    feedFollow_params := database.CreateFeedFollowParams{
        uuid.New(),
        time.Now(),
        time.Now(),
        feed.ID,
        current_user.ID,
    }

    feed_follow, err := s.db.CreateFeedFollow(ctx, feedFollow_params)
    if err != nil {
        return err
    }

    fmt.Printf("Feed_follow has been created\n")
    fmt.Printf("User: \"%s\" Title: \"%s\" \n", feed_follow.UserName, feed_follow.FeedName)

    return nil
}

func handlerFollowing(s *state, cmd command, current_user database.User) error {
    if len(cmd.args) != 0 {
        return fmt.Errorf("There should be no arguments for following command")
    }

    username := s.cfg.CurrentUserName

    ctx := context.Background()

    feed_follows, err := s.db.GetFeedFollowsForUser(ctx, username)
    if err != nil {
        return err
    }

    fmt.Printf("Feeds user \"%s\" is currently following\n", username)
    for i := 0; i < len(feed_follows); i++ {
        fmt.Printf("Feed name: \"%s\"\n", feed_follows[i].FeedName)
    }

    return nil
}

func handlerUnfollow(s *state, cmd command, current_user database.User) error {
    if len(cmd.args) != 1 {
        return fmt.Errorf("There should be 1 argument for unfollowing command, feed_url")
    }

    url := cmd.args[0]

    ctx := context.Background()

    params := database.DeleteFeedFollowByUrlAndNameParams{
        url,
        current_user.Name,
    }

    err := s.db.DeleteFeedFollowByUrlAndName(ctx, params)
    if err != nil {
        return err
    }

    fmt.Printf("User \"%s\" unfollowed feed with url \"%s\"", current_user.Name, url)

    return nil
}

func handlerBrowse(s *state, cmd command, current_user database.User) error {
    if len(cmd.args) < 0 && len(cmd.args) > 1 {
        return fmt.Errorf("There should be 1 or 0 arguments for command, limit")
    }


    limit := int32(2)
    if len(cmd.args) == 1 {
        i, err := strconv.ParseInt(cmd.args[0], 10, 32) 
        if err != nil {
            return err
        }
        limit = int32(i)
    }

    ctx := context.Background()
    query_params := database.GetPostsForUserParams{
        s.cfg.CurrentUserName,
        limit,
    }


    posts, err := s.db.GetPostsForUser(ctx, query_params)
    if err != nil {
        return err
    }
    
    fmt.Printf("%s posts first %d posts\n", s.cfg.CurrentUserName, limit)

    if len(posts) == 0 {
        fmt.Printf("No posts found!\n")
    }

    for i := 0; i < len(posts); i++ {
        fmt.Printf("Published_at: \"%s\" Title: \"%s\" Description: \"%s\"\n", posts[i].PublishedAt, posts[i].Title, posts[i].Description)
    }

    return nil
}
