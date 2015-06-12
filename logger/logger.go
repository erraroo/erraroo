package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/jserr/jserr/config"
)

func init() {
	//logrus.SetLevel(logrus.DebugLevel)
	logrus.SetLevel(logrus.InfoLevel)
	if config.Env == "production" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	//logrusrus.AddHook(airbrake.NewHook("https://example.com", "xyz", "development"))
}

func Debug(args ...interface{}) {
	message := args[0]
	fields := makeFields(args...)
	logrus.WithFields(fields).Debug(message)
}

func Error(args ...interface{}) {
	message := args[0]
	fields := makeFields(args...)
	logrus.WithFields(fields).Error(message)
}

func Info(args ...interface{}) {
	message := args[0]
	fields := makeFields(args...)
	logrus.WithFields(fields).Info(message)
}

func Fatal(args ...interface{}) {
	message := args[0]
	fields := makeFields(args...)
	logrus.WithFields(fields).Fatal(message)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args)
}

func makeFields(args ...interface{}) logrus.Fields {
	fields := logrus.Fields{}
	length := len(args)
	if length < 3 && (length-1)%2 != 0 {
		return fields
	}

	for i := 1; i < len(args); i += 2 {
		fields[args[i].(string)] = args[i+1]
	}

	return fields
}
