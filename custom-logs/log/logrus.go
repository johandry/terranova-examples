package log

import (
	"github.com/johandry/terranova/logger"
	"github.com/sirupsen/logrus"
)

// Logrus returns a Terranova Logger Middleware to intercept every log entry
// from Terraform and print it using the sirupsen/logrus package
func Logrus(count int) *logger.Middleware {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetLevel(logrus.InfoLevel)

	return logger.NewMiddleware(l.WithFields(logrus.Fields{
		"platform": "AWS",
		"count":    count,
	}))
}
