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

// InitLogger initializes logger instance
// This method MUST be called before calling any other method.
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
func Fatal(format string, args ...any) {
	checkInitState()
	gLogger.Fatalf(format, args...)
}

// Error ...
func Error(format string, args ...any) {
	checkInitState()
	gLogger.Errorf(format, args...)
}

// Warn ...
func Warn(format string, args ...any) {
	checkInitState()
	gLogger.Warningf(format, args...)
}

// Info ...
func Info(format string, args ...any) {
	checkInitState()
	gLogger.Infof(format, args...)
}

// Debug ...
func Debug(format string, args ...any) {
	checkInitState()
	gLogger.Debugf(format, args...)
}

func checkInitState() {
	if gLogger == nil {
		panic("logger not initialized")
	}
}
