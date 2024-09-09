package log

import (
	"log"
	"os"
	"strings"
)

var (
	debuglogger   *log.Logger
	infologger    *log.Logger
	warninglogger *log.Logger
	errorlogger   *log.Logger
)

func Debugf(msg string, v ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}
	if debuglogger == nil {
		debuglogger = log.New(os.Stderr, "DEBUG ", log.Default().Flags())
	}
	debuglogger.Printf(msg, v...)
}

func Infof(msg string, v ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}
	if infologger == nil {
		infologger = log.New(os.Stderr, "INFO ", log.Default().Flags())
	}
	infologger.Printf(msg, v...)
}

func Warningf(msg string, v ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}
	if warninglogger == nil {
		warninglogger = log.New(os.Stderr, "WARNING ", log.Default().Flags())
	}
	warninglogger.Printf(msg, v...)
}

func Errorf(msg string, v ...any) {
	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}
	if errorlogger == nil {
		errorlogger = log.New(os.Stderr, "ERROR ", log.Default().Flags())
	}
	errorlogger.Printf(msg, v...)
}

func Fatal(v ...any) {
	log.Fatal(v...)
}
