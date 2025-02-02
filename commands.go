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
    "html"
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
    if len(cmd.args) != 0 {
        return fmt.Errorf("There should be no arguments for agg command")
    }

    const feedURL = "https://www.wagslane.dev/index.xml"
    ctx := context.Background()

    RSSFeed, err := fetchFeed(ctx, feedURL)
    if err != nil {
        return err
    }

    RSSFeed.Channel.Title = html.UnescapeString(RSSFeed.Channel.Title)
    RSSFeed.Channel.Description = html.UnescapeString(RSSFeed.Channel.Description)
    for i := 0; i < len(RSSFeed.Channel.Item); i++ {
        RSSFeed.Channel.Item[i].Title = html.UnescapeString(RSSFeed.Channel.Item[i].Title)
        RSSFeed.Channel.Item[i].Description = html.UnescapeString(RSSFeed.Channel.Item[i].Description)
    }

    fmt.Printf("%+v\n", RSSFeed)

    return nil
}

func handlerAddFeed(s *state, cmd command) error {
    if len(cmd.args) != 2 {
        return fmt.Errorf("the addfeed handler expects a two argument, name and url.")
    }

    name := cmd.args[0]
    url := cmd.args[1]

    ctx := context.Background()
    current_user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
    if err != nil {
        return err
    }
    
    feed_params := database.CreateFeedParams{
        uuid.New(),
        time.Now(),
        time.Now(),
        name,
        url,
        current_user.ID,
    }

    _, err = s.db.CreateFeed(ctx, feed_params)
    if err != nil {
        return err
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
