package log

import "github.com/johandry/terranova/logger"

// Discard returns a Terranova Logger Middleware to intercept every log entry
// from Terraform. The Discard middleware do not print any Terraform log.
func Discard() *logger.Middleware {
	l := logger.DiscardLog()
	return logger.NewMiddleware(l)
}
