package main

import _ "github.com/lib/pq" //Import for side effect, not because you need it

import (
    "github.com/erlint1212/blog_aggregator/internal/config"
    "github.com/erlint1212/blog_aggregator/internal/database"
    "os"
    "fmt"
    "log"
    "database/sql"
)


func check(err error) {
    if err != nil {
        log.Fatalf("error: %v\n", err)
    }
}

func main() {

    cfg, err := config.Read()
    check(err)

    db, err := sql.Open("postgres", cfg.DbUrl)
    check(err)

    dbQueries := database.New(db)

    cfg_state := state{dbQueries, &cfg}

    commands := commands{make(map[string]func(*state, command) error)}

    commands.register("login", handlerLogin)
    commands.register("register", handlerRegister)
    commands.register("reset", handlerReset)
    commands.register("users", handlerUsers)
    commands.register("agg", handlerAgg)

    args := os.Args[1:] //Without prog
    if len(args) < 1 {
        check(fmt.Errorf("Need more than 1 argument"))
    }
    cmd := command{args[0], args[1:]}

    err = commands.run(&cfg_state, cmd)
    check(err)
}
