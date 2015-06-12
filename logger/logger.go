package logger

import (
	"os"

	"github.com/erraroo/erraroo/config"
	"github.com/inconshreveable/log15"
)

var root log15.Logger

func init() {
	root = config.Logger()
}

func Debug(msg string, args ...interface{}) {
	root.Debug(msg, args...)
}

func Error(msg string, args ...interface{}) {
	root.Error(msg, args...)
}

func Info(msg string, args ...interface{}) {
	root.Info(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	root.Crit(msg, args...)
	os.Exit(1)
}
