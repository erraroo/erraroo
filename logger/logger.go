package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/erraroo/erraroo/config"
)

var logger = logrus.New()

func init() {
	logger.Level = logrus.DebugLevel

	if config.Env == "production" {
		logger.Level = logrus.InfoLevel
		logger.Formatter = &logrus.JSONFormatter{}
	}
}

func Debug(msg string, args ...interface{}) {
	logger.WithFields(makeFields(args...)).Debug(msg)
}

func Info(msg string, args ...interface{}) {
	logger.WithFields(makeFields(args...)).Info(msg)
}

func Error(msg string, args ...interface{}) {
	logger.WithFields(makeFields(args...)).Error(msg)
}

func Fatal(msg string, args ...interface{}) {
	logger.WithFields(makeFields(args...)).Fatal(msg)
}

func makeFields(args ...interface{}) logrus.Fields {
	fields := logrus.Fields{}

	for i := 0; i < len(args)-1; i += 2 {
		if name, ok := args[i].(string); ok {
			fields[name] = args[i+1]
		}
	}

	return fields
}
