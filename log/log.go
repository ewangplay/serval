package log

import (
	"io"
	"strings"

	l "github.com/apsdehal/go-logger"
)

// LoggerConfig defines the config for Logger
type LoggerConfig struct {
	Module   string
	Color    int
	LogLevel string
	Writer   io.Writer
}

var gLogger *l.Logger

// InitLogger init logger instance
func InitLogger(conf *LoggerConfig) (err error) {
	// If logger has been initialized, return directly
	if gLogger != nil {
		return nil
	}

	gLogger, err = l.New(
		conf.Module,
		conf.Color,
		conf.Writer,
		parseLogLevel(conf.LogLevel))
	if err != nil {
		return err
	}

	return nil
}

func parseLogLevel(levelstr string) l.LogLevel {
	lvl := l.InfoLevel
	switch strings.ToLower(levelstr) {
	case "debug":
		lvl = l.DebugLevel
	case "info":
		lvl = l.InfoLevel
	case "warn":
		lvl = l.WarningLevel
	case "error":
		lvl = l.ErrorLevel
	case "fatal":
		lvl = l.CriticalLevel
	}
	return lvl
}

// Fatal ...
func Fatal(format string, args ...interface{}) {
	if gLogger == nil {
		panic("logger not initialized")
	}
	gLogger.Fatalf(format, args...)
}

// Error ...
func Error(format string, args ...interface{}) {
	if gLogger == nil {
		panic("logger not initialized")
	}
	gLogger.Errorf(format, args...)
}

// Warn ...
func Warn(format string, args ...interface{}) {
	if gLogger == nil {
		panic("logger not initialized")
	}
	gLogger.Warningf(format, args...)
}

// Info ...
func Info(format string, args ...interface{}) {
	if gLogger == nil {
		panic("logger not initialized")
	}
	gLogger.Infof(format, args...)
}

// Debug ...
func Debug(format string, args ...interface{}) {
	if gLogger == nil {
		panic("logger not initialized")
	}
	gLogger.Debugf(format, args...)
}
