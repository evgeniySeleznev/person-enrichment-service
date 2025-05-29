package logger

import (
	"io"
	"log"
	"os"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type logger struct {
	debug *log.Logger
	info  *log.Logger
	error *log.Logger
	fatal *log.Logger
}

func NewLogger(out io.Writer) Logger {
	if out == nil {
		out = os.Stdout
	}

	return &logger{
		debug: log.New(out, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		info:  log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatal: log.New(out, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.debug.Printf(msg+"\n", args...)
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.info.Printf(msg+"\n", args...)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.error.Printf(msg+"\n", args...)
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.fatal.Fatalf(msg+"\n", args...)
}
