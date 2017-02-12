package state

import (
    log "github.com/Sirupsen/logrus"
    "fmt"
)

type AppLogger struct {
    *log.Logger
}

var Logger *AppLogger

func InitLogger(level string) *AppLogger {
    Logger = &AppLogger{log.New()}
    Logger.enableTimestamp("2006-01-02T15:04:05.1")
    Logger.setLevel(level)
    return Logger
}

func (l *AppLogger) enableTimestamp(format string) {
    formatter := new(log.TextFormatter)
    formatter.TimestampFormat = format
    formatter.FullTimestamp = true
    l.Formatter = formatter
}

func (l *AppLogger) setLevel(level string) {
    lvl, err := log.ParseLevel(level)
    Logger.Level = lvl
    if err != nil {
        log.Fatalf(fmt.Sprintf("Unable to parse '%v' as a logging level", level))
    }
}
