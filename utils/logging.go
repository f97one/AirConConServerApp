package utils

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var logger *logrus.Logger

const (
	logFilename string = "airconcon.log"
)

func GetLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		conf := Load("")

		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006/01/02 15:04:05.000 MST",
		})
		logger.SetReportCaller(true)

		if conf.ProductionMode {
			logger.SetOutput(&lumberjack.Logger{
				Filename:  filepath.Join("./logs", logFilename),
				MaxAge:    7,
				LocalTime: true,
				Compress:  true,
			})
			logger.SetLevel(logrus.WarnLevel)
		} else {
			logger.SetOutput(os.Stdout)
			logger.SetLevel(logrus.TraceLevel)
		}
	}

	return logger
}
