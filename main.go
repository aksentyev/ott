package main

import (
    "flag"
    "github.com/aksentyev/ott/state"
    "github.com/aksentyev/ott/app"
)

const defaultRedisURL = "redis://127.0.0.1:6379"

var (
    getErrors       bool
    genMsgInterval  int
    genMsgLimit     int
    redisURL        string
    loglevel        string
)

func init() {
    flag.BoolVar(&getErrors, "getErrors", false, "get errored messages from Redis and exit")
    flag.IntVar(&genMsgInterval, "interval", 500, "messages generation interval")
    flag.IntVar(&genMsgLimit, "limit", 0, "messages count limit (default 0)")
    flag.StringVar(&redisURL, "redis", defaultRedisURL, "redis url")
    flag.StringVar(&loglevel, "loglevel", "info", "logging level")
}

func setup() {
    flag.Parse()
    state.InitLogger(loglevel)
    state.InitRedisClient(redisURL)
}

func main() {
    setup()
    if getErrors {
        app.PrintErroredMessages()
        return
    }

    state.Logger.Debug("Starting app")
    app.NewSupervisor(genMsgInterval, genMsgLimit).Run()
}
