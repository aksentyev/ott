package app

import (
    s "github.com/aksentyev/ott/state"
    "fmt"
    "time"
    "os"
)

type Generator struct{
    GenInterval   int
    MessagesLimit int
    DoneCh        chan bool
    counter       int
}

func (g *Generator) Run() {
    s.Logger.Info("generator started")

    go g.generatorHeartBeats()

    for {
        select {
        case <-g.DoneCh:
            s.Logger.Info("generator stopped")
            return
        default:
            g.push(g.generateMsg())
            g.checkNumMessagesGenerated()

            if g.GenInterval > 0 {
                time.Sleep(time.Duration(g.GenInterval) * time.Millisecond)
            }
        }
    }
}

func (g *Generator) markSelfAsGenerator() (ok bool) {
    // set generator key if not exists with 1 second ttl to give generator enough time to start
    err := s.Redis.SetNX(generatorRedisKey, true, time.Second).Err()

    if err != nil {
        s.Logger.Info("generator already exists")
        return false
    }

    return true
}

func (g *Generator) generatorHeartBeats() {
    for {
        ok, err := s.Redis.SetXX(generatorRedisKey, true, generatorPingTTL).Result()
        if err != nil || !ok {
            g.DoneCh<- true
            return
        }

        time.Sleep(generatorPingInterval)
    }
}

func (g *Generator) checkNumMessagesGenerated() {
    if g.counter == g.MessagesLimit {
        s.Logger.Infof("%v messages were generated", g.counter)
        os.Exit(0)
    }
    if g.counter % 100000 == 0 {
        s.Logger.Infof("%v messages were generated", g.counter)
    }
}

func (g *Generator) generateMsg() string {
    msg := fmt.Sprintf("{msg: %v}", time.Now().UnixNano())
    g.counter++
    return msg
}

func (g *Generator) push(msg string) error {
    return s.Redis.RPush(messagesRedisKey, msg).Err()
}
