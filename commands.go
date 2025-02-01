package main

import (
    "errors"
    "database/sql"
    "os"
    "github.com/erlint1212/blog_aggregator/internal/config"
    "github.com/erlint1212/blog_aggregator/internal/database"
    "fmt"
    "time"
    "context"
	"github.com/google/uuid"
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
        os.Exit(1)
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
