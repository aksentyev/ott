package app

import (
    s "github.com/aksentyev/ott/state"
)

type Supervisor struct {
    Generator *Generator
    Handler   *Handler
    doneChan chan bool
}

type runner interface {
    Run()
}

func NewSupervisor(genInterval, genLimit int) *Supervisor {
    doneChan := make(chan bool)

    return &Supervisor{
        Generator: &Generator{
            MessagesLimit: genLimit,
            GenInterval: genInterval,
        },
        Handler: &Handler{
            DoneCh: doneChan,
        },
        doneChan: doneChan,
    }
}

func (sv *Supervisor) Run() {
    var service runner = sv.Handler

    roleChanged := true
    for {
        if !isGeneratorUp(){
            s.Logger.Info("Generator is not running. Trying to become generator")
            if sv.Generator.markSelfAsGenerator() {

                go func() { sv.doneChan<- true }()
                close(sv.doneChan)

                sv.doneChan = make(chan bool)
                service = sv.Generator
                roleChanged = true
            }
        }
        if roleChanged {
            go service.Run()
            roleChanged = false
        }
    }
}

func isGeneratorUp() bool {
    return s.Redis.Exists(generatorRedisKey).Val()
}
