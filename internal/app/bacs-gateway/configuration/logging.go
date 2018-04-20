package configuration

import (
	"bufio"
	"os"
	"strings"
	"github.com/sirupsen/logrus"
)

// LoggingConfig specifies all the parameters needed for logging
type LoggingConfig struct {
	Level string
	File  string
}

const (
	DEFAULT_LOG_LEVEL = "INFO"
)

// ConfigureLogging will take the logging configuration and also adds
// a few default parameters
func ConfigureLogging(config *LoggingConfig) (*logrus.Entry, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// use a file if you want
	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return nil, errOpen
		}
		logrus.SetOutput(bufio.NewWriter(f))
	}

	logrus.ParseLevel(strings.ToUpper(GetOrDefault("logging.level", DEFAULT_LOG_LEVEL)))

	// always use the fulltimestamp
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
	})

	return logrus.StandardLogger().WithField("hostname", hostname), nil
}
