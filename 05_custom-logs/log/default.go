package log

import "github.com/johandry/terranova/logger"

// Default returns a Terranova Logger Middleware to intercept every log entry
// from Terraform and print it. The Default middleware just print the log level
// entries Info, Warn and Error from Terraform in the format:
// LEVEL [date] message
func Default() *logger.Middleware {
	return logger.NewMiddleware()
}
