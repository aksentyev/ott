package app

import (
    s "github.com/aksentyev/ott/state"
    "fmt"
)

func PrintErroredMessages() {
    list := getErroredMessages()
    err := deleteMessagesWithErrorsList()
    if err != nil {
        s.Logger.Errorf("an error occured while removing errors list: %v", err)
    }

    if len(list) == 0 {
        s.Logger.Info("no messages with errors found")
    }
    for _, msg := range list {
        fmt.Printf("message with error: %v\n", msg)
    }
}

func getErroredMessages() []string {
    res, err := s.Redis.LRange("errors", 0, -1).Result()
    if err != nil {
        panic(err)
    }
    return res
}

func deleteMessagesWithErrorsList() error {
    return s.Redis.Del("errors").Err()
}
