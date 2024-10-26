package main

import (
    "github.com/erlint1212/blog_aggregator/internal/config"
    "fmt"
)

type state struct {
    CfgPointer          *config.Config
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

    err := s.CfgPointer.SetUser(cmd.args[0])
    if err != nil {
        return err
    }

    fmt.Printf("User %s has been set\n", cmd.args[0])

    return nil
}
