package log

import "github.com/johandry/terranova/logger"

// Terraform returns a Terranova Logger Middleware to intercept every log entry
// from Terraform. In this case, the Terraform middleware do not modify any
// Terraform log entry and prints them as they are. It basicaly do nothing.
func Terraform() *logger.Middleware {
	return nil
}
