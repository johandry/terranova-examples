package log

import (
	"log"
	"os"

	"github.com/johandry/terranova/logger"
)

// Custom returns a Terranova Logger Middleware to intercept every log entry
// from Terraform and print it in a custom way
func Custom() *logger.Middleware {
	l := NewMyLog()
	return logger.NewMiddleware(l)
}

// MyLog is an implementation of a Logger that do nothing. It's simmilar to
// Discard.
type MyLog struct {
	log *log.Logger
}

// NewMyLog create a Terranova Logger that do nothing.
func NewMyLog() logger.Logger {
	l := log.New(os.Stderr, "", log.LstdFlags)
	return &MyLog{
		log: l,
	}
}

// Printf implements a standard Printf function of Logger interface
func (l *MyLog) Printf(format string, args ...interface{}) {
	l.output("     ", format, args...)
}

// Debugf implements a standard Debugf function of Logger interface
func (l *MyLog) Debugf(format string, args ...interface{}) {}

// Infof implements a standard Infof function of Logger interface
func (l *MyLog) Infof(format string, args ...interface{}) {
	l.output("INFO ", format, args...)
}

// Warnf implements a standard Warnf function of Logger interface
func (l *MyLog) Warnf(format string, args ...interface{}) {
	l.output("WARN ", format, args...)
}

// Errorf implements a standard Errorf function of Logger interface
func (l *MyLog) Errorf(format string, args ...interface{}) {
	l.output("ERROR", format, args...)
}

func (l *MyLog) output(levelStr string, format string, args ...interface{}) {
	l.log.SetPrefix(levelStr + " [ ")
	l.log.Printf("] "+format, args...)
}
