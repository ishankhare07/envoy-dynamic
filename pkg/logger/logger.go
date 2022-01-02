package logger

import (
	"log"
)

type Logger struct {
	Debug bool
}

// Debugf logs a formatted debugging message.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
		log.Printf(format+"\n", args...)
	}
}

// Infof logs a formatted informational message.
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.Debug {
		log.Printf(format+"\n", args...)
	}
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}
