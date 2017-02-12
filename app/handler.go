package app

import (
    s "github.com/aksentyev/ott/state"
    "math/rand"
    "errors"
)

type Handler struct{
    DoneCh  <-chan bool
    counter int
}

func (h *Handler) Run() {
    s.Logger.Info("handler started")

    for {
        select {
        case <-h.DoneCh:
            s.Logger.Info("handler stopped")
            return
        default:
            msg, ok := h.getMessage()
            if ok {
                s.Logger.Debug("msg handle started")
                h.handleMessage(msg)
                h.checkNumMessagesHandled()
                s.Logger.Debug("msg handle finished")
            } else {
                s.Logger.Debug("no messages in the queue")
            }
        }
    }
}

func (h *Handler) getMessage() (string, bool) {
    res, err := s.Redis.BLPop(generatorPingTTL, messagesRedisKey).Result()
    if err != nil {
        return "", false
    }
    return res[1], true
}

func (h *Handler) handleMessage(msg string) {
    if err := h.checkError(msg);  err != nil {
        s.Logger.Debugf("msg %v was errored. saving to errors set.")
        h.saveErrored(msg)
    }
    h.counter++
}

func (h *Handler) checkError(msg string) error {
    if chance := rand.Float32(); chance >= 0.95 {
        s.Logger.Debugf("an error occured while processing message: %v",msg)
        return errors.New("some error")
    }
    return nil
}

func (h *Handler) saveErrored(msg string) {
    s.Redis.RPush(erroredRedisKey, msg)
}

func (h *Handler) checkNumMessagesHandled() {
    if h.counter % 100000 == 0 {
        s.Logger.Infof("%v messages were handled", h.counter)
    }
}
