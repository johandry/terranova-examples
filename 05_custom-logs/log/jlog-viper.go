package log

import (
	"os"

	"github.com/johandry/log"
	"github.com/johandry/terranova/logger"
	"github.com/spf13/viper"
)

// JLogViper returns a Terranova Logger Middleware to intercept every log entry
// from Terraform and print it in a custom way. JLogViper use the johandry/log
// package which uses sirupsen/logrus to format the logs and spf13/viper to
// configure it
func JLogViper() *logger.Middleware {
	v := viper.New()
	v.Set(log.OutputKey, os.Stderr)
	v.Set(log.ForceColorsKey, true)
	v.Set(log.DisableColorsKey, false)
	v.Set(log.LevelKey, "debug")
	v.Set(log.PrefixField, "Terranova")

	l := log.New(v)
	return logger.NewMiddleware(l)
}
