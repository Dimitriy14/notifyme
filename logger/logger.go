package logger

import (
	"github.com/Dimitriy14/notifyme/config"
	"github.com/sirupsen/logrus"

	"os"
)

var (
	Log Logger
)

// Logger - represents methods for logging
type Logger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

// Load loads logger
func Load() error {
	output := os.Stdout
	if config.Conf.LogFile != "" {
		logFile, err := os.Create(config.Conf.LogFile)
		if err != nil {
			return err
		}
		output = logFile
	}

	logLvl, err := logrus.ParseLevel(config.Conf.LogLevel)
	if err != nil {
		return err
	}

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(output)
	log.SetLevel(logLvl)

	Log = log

	return nil
}
