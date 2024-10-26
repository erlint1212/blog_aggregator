package main

import (
    "github.com/erlint1212/blog_aggregator/internal/config"
    "os"
    "fmt"
)

func check(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
}

func main() {
    cfg, err := config.Read()
    check(err)

    cfg_state := state{&cfg}

    commands := commands{make(map[string]func(*state, command) error)}

    commands.register("login", handlerLogin)

    args := os.Args[1:] //Without prog
    if len(args) < 1 {
        check(fmt.Errorf("Need more than 1 argument"))
    }
    cmd := command{args[0], args[1:]}

    err = commands.run(&cfg_state, cmd)
    check(err)
}
