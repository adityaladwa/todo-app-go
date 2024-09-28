package logger

import (
    "strings"

    "github.com/sirupsen/logrus"
)

func NewLogger(level string) *logrus.Logger {
    log := logrus.New()
    lvl, err := logrus.ParseLevel(strings.ToLower(level))
    if err != nil {
        log.SetLevel(logrus.InfoLevel)
    } else {
        log.SetLevel(lvl)
    }
    log.SetFormatter(&logrus.JSONFormatter{})
    return log
}
