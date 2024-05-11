package logger

import "github.com/sirupsen/logrus"

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		TimestampFormat:        "15:04:05",
		DisableLevelTruncation: true,
	})

	logrus.SetLevel(logrus.DebugLevel)
}
