// Package human is a convenience wrapper for logging human-readable messages of various
// levels to the terminal. Much of this code has been taken from
// https://github.com/helm/helm/blob/d87ce93e1e287ece84d940dbfe09b0de493d9953/pkg/kube/log.go
//
// Thank you Helm team!
package log

import (
	"io"
	"os"
)

// Stdout is the logging destination for normal messages.
var Stdout io.Writer = os.Stdout

// Stderr is the logging destination for error messages.
var Stderr io.Writer = os.Stderr

// IsDebugging toggles whether or not to enable debug output and behavior.
var IsDebugging = false

var defaultLogger = NewDefault(IsDebugging)

// SetDebug sets the debug flag for the default logger
func SetDebug(dbg bool) {
	defaultLogger = NewDefault(dbg)
}

// Msg passes through the formatter, but otherwise prints exactly as-is.
//
// No prettification.
func Msg(format string, v ...interface{}) {
	defaultLogger.Msg(format, v...)
}

// Check checks err, prints an error message using Err(), and returns err
// all if err != nil
// Otherwise, prints the success format string, formatted with successArgs
// using Msg, and returns nil.
func Check(err error, successFmt string, successArgs ...interface{}) error {
	return defaultLogger.Check(err, successFmt, successArgs...)
}

// Err prints an error message. It does not cause an exit.
func Err(format string, v ...interface{}) {
	defaultLogger.Err(format, v...)
}

// ErrRet does the same thing as Err(format, v...), except returns an
// error with the given format string and arguments
func ErrRet(format string, v ...interface{}) error {
	return defaultLogger.ErrRet(format, v...)
}

// Info prints a green-tinted message.
func Info(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

// Debug prints a cyan-tinted message if IsDebugging is true.
func Debug(msg string, v ...interface{}) {
	defaultLogger.Debug(msg, v...)
}

// Warn prints a yellow-tinted warning message.
func Warn(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}
