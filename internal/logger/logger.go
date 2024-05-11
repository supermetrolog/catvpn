package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/pkg/logger"
)

func NewLogger(f logrus.Formatter, l logrus.Level) logger.Logger {
	log := logrus.New()
	log.SetFormatter(f)
	log.SetLevel(l)

	return log
}
