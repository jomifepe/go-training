package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

var (
	Logger *logrus.Logger
)

// NewLogger creates and configures a new logrus Logger.
func NewLogger() *logrus.Logger {
	Logger = logrus.New()
	if viper.GetBool("LOG_FORMAT_JSON") {
		Logger.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: false,
		}
	} else {
		Logger.Formatter = &logrus.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp: false,
		}
	}

	level, err := logrus.ParseLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}
	Logger.Level = level
	return Logger
}