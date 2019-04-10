package log

import (
	"fmt"
	"io"
	"os"

	"github.com/labstack/gommon/color"
)

// Logger is the interface for logging
type Logger interface {
	Msg(string, ...interface{})
	Check(error, string, ...interface{}) error
	Err(string, ...interface{})
	ErrRet(string, ...interface{}) error
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Warn(string, ...interface{})
}
type logger struct {
	Stdout  io.Writer
	Stderr  io.Writer
	IsDebug bool
}

// NewDefault sets up a logger with defaults
//
// dbg is whether to turn on debug logging or not
func NewDefault(dbg bool) Logger {
	return &logger{
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		IsDebug: dbg,
	}
}

func (l *logger) Msg(format string, v ...interface{}) {
	fmt.Fprintf(l.Stdout, appendNewLine(format), v...)
}

func (l *logger) Check(err error, successFmt string, successArgs ...interface{}) error {
	if err != nil {
		l.Err(err.Error())
		return err
	}
	Msg(successFmt, successArgs...)
	return nil
}

func (l *logger) Err(format string, v ...interface{}) {
	fmt.Fprint(l.Stderr, color.Red("[ERROR] "))
	fmt.Fprintf(l.Stderr, appendNewLine(format), v...)
}

func (l *logger) ErrRet(format string, v ...interface{}) error {
	l.Err(format, v...)
	return fmt.Errorf(format, v...)
}

func (l *logger) Info(format string, v ...interface{}) {
	fmt.Fprint(l.Stderr, "---> ")
	fmt.Fprintf(l.Stderr, appendNewLine(format), v...)
}

func (l *logger) Debug(msg string, v ...interface{}) {
	if l.IsDebug {
		fmt.Fprint(l.Stderr, color.Cyan("[DEBUG] "))
		l.Msg(msg, v...)
	}
}

// Warn prints a yellow-tinted warning message.
func (l *logger) Warn(format string, v ...interface{}) {
	fmt.Fprint(l.Stderr, color.Yellow("[WARN] "))
	Msg(format, v...)
}

func appendNewLine(format string) string {
	return format + "\n"
}
