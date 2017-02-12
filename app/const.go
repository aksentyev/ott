package app

import (
    "time"
)

const (
    messagesRedisKey  = "messages"
    erroredRedisKey   = "errors"
    generatorRedisKey = "generator"

    generatorPingInterval = time.Duration(50) * time.Millisecond
    generatorPingTTL      = time.Duration(200) * time.Millisecond
)
