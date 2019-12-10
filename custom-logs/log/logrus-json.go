package log

import (
	"github.com/johandry/terranova/logger"
	"github.com/sirupsen/logrus"
)

// LogrusJSON returns a Terranova Logger Middleware to intercept every log entry
// from Terraform and print it using the sirupsen/logrus package and a JSON format
func LogrusJSON(count int) *logger.Middleware {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.InfoLevel)

	return logger.NewMiddleware(l.WithFields(logrus.Fields{
		"platform": "AWS",
		"count":    count,
	}))
}
