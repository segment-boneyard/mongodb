package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/segmentio/kit/schema"
)

// defaultLogger defines the internal default logging
var defaultLogger = logrus.New()

type M logrus.Fields

// signalListener listens on SIGHUP, SIGINT, SIGTERM, SIGQUIT
// opens a 5 minute window to debug level
func signalListener(serviceSchema schema.Service) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	for range c {
		defaultLogger.Level = logrus.DebugLevel
		go func() {
			time.Sleep(1 * time.Minute)
			defaultLogger.Level = logrus.InfoLevel
		}()
	}
}

func Init(serviceSchema schema.Service) error {
	go signalListener(serviceSchema)
	defaultLogger.Level = logrus.InfoLevel

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		lvl, err := logrus.ParseLevel(level)
		if err == nil {
			defaultLogger.Level = lvl
		}
	}

	if logrus.IsTerminal() {
		defaultLogger.Formatter = &logrus.TextFormatter{}
	} else {
		defaultLogger.Formatter = &logrus.JSONFormatter{}
	}

	return nil
}

func SetWriter(out io.Writer) {
	defaultLogger.Out = out
}

// Debugf formats message according to format specifier
// and writes to log with level = Debug.
func Debugf(format string, params ...interface{}) {
	defaultLogger.Debugf(format, params...)
}

// Infof formats message according to format specifier
// and writes to log with level = Info.
func Infof(format string, params ...interface{}) {
	defaultLogger.Infof(format, params...)
}

// Warnf formats message according to format specifier
// and writes to log with level = Warn.
func Warnf(format string, params ...interface{}) error {
	defaultLogger.Warnf(format, params...)
	return fmt.Errorf(format, params...)
}

// Errorf formats message according to format specifier
// and writes to log with level = Error.
func Errorf(format string, params ...interface{}) error {
	defaultLogger.Errorf(format, params...)
	return fmt.Errorf(format, params...)
}

// Criticalf formats message according to format specifier
// and writes to log with level = Critical.
func Fatalf(format string, params ...interface{}) error {
	defaultLogger.Fatalf(format, params...)
	return fmt.Errorf(format, params...)
}

func Panicf(format string, params ...interface{}) {
	defaultLogger.Panicf(format, params...)
}

// Debug formats message using the default formats for its operands
// and writes to log with level = Debug
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Info formats message using the default formats for its operands
// and writes to log with level = Info
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Warn formats message using the default formats for its operands
// and writes to log with level = Warn
func Warn(v ...interface{}) error {
	defaultLogger.Warn(v...)
	return fmt.Errorf(strings.Repeat("%v ", len(v)), v...)
}

// Error formats message using the default formats for its operands
// and writes to log with level = Error
func Error(v ...interface{}) error {
	defaultLogger.Error(v...)
	return fmt.Errorf(strings.Repeat("%v ", len(v)), v...)
}

// Critical formats message using the default formats for its operands
// and writes to log with level = Critical
func Fatal(v ...interface{}) error {
	defaultLogger.Fatal(v...)
	return fmt.Errorf(strings.Repeat("%v ", len(v)), v...)
}

func Panic(v ...interface{}) {
	defaultLogger.Panic(v...)
}

// Inspect returns a string that inspects the passed data structures and dumps
// the contents in a human readable manner using `go-spew`
//
// Example:
// 	(main.Foo) {
//  	unexportedField: (*main.Bar)(0xf84002e210)({
//  		flag: (main.Flag) flagTwo,
//   		data: (uintptr) <nil>
//  	}),
//  	ExportedField: (map[interface {}]interface {}) {
//   		(string) "one": (bool) true
//  	}
// 	}
//
// Use of this method should be for debugging purposes only and not
// shipped in staging or production environments
func Inspect(v ...interface{}) string {
	return spew.Sdump(v...)
}

func toM(v interface{}) M {
	var x M
	res, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(res, &x); err != nil {
		return nil
	}
	return x
}

func With(v interface{}) *logrus.Entry {
	if value, ok := v.(M); ok {
		return defaultLogger.WithFields(logrus.Fields(value))
	} else {
		return defaultLogger.WithFields(logrus.Fields(toM(v)))
	}
}
