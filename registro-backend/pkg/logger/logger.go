package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init(level string) {
	Log.SetOutput(os.Stdout)

	// Set JSON formatter for production-ready logs
	Log.SetFormatter(&logrus.JSONFormatter{})

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	Log.SetLevel(parsedLevel)
}
