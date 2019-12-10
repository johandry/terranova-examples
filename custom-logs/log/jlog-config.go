package log

import (
	"os"

	"github.com/johandry/log"
	"github.com/johandry/terranova/logger"
	"github.com/spf13/viper"
)

const configFilename = "config"

// JLogConfig returns a Terranova Logger Middleware to intercept every log entry
// from Terraform and print it in a custom way. JLogConfig use the johandry/log
// package which uses sirupsen/logrus to format the logs and spf13/viper to
// configure it from code, environment variables and config files.
func JLogConfig() *logger.Middleware {
	v := viper.New()

	// Set default parameters via code:
	v.SetDefault(log.OutputKey, os.Stderr)
	v.SetDefault(log.ForceColorsKey, true)
	v.SetDefault(log.DisableColorsKey, false)
	v.SetDefault(log.LevelKey, "debug")
	v.SetDefault(log.PrefixField, "Terranova")

	// Set parameters from environment variables:
	v.BindEnv(log.LevelKey)

	// Set parameters from the config file
	v.SetConfigName(configFilename)
	v.AddConfigPath(".")
	v.ReadInConfig()

	l := log.New(v)
	return logger.NewMiddleware(l)
}

// Set some environment variables in case they are not set by the user.
// This is just a trick to test this custom log, there is no need to do it in
// production
func init() {
	if os.Getenv(log.LevelKey) != "" {
		os.Setenv(log.LevelKey, "warn")
	}
}
