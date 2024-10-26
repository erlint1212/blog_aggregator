package main

import (
    "fmt"
    "github.com/erlint1212/blog_aggregator/internal/config"
)

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    cfg, err := config.Read()
    check(err)

    err = cfg.SetUser("Erling")
    check(err)

    cfg, err = config.Read()
    check(err)

    fmt.Println(cfg)
}
