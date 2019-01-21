package mlog

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type ILogger interface {
	SetLevel(logrus.Level)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

type Level int

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logrus Level: %q", lvl)
}

// LogConfig log config
type LogConfig struct {
	Dir        string `toml:"dir,omitempty"`
	Buffer     int    `toml:"buffer,omitempty"`
	SubRelPath string `toml:"sub_rel_path,omitempty"`
	LogLevel   string `toml:"log_level,omitempty"`
	ServerIP   string `toml:"server_ip,omitempty"`
}

// DefaultLogConfig default log config
func DefaultLogConfig() LogConfig {
	return LogConfig{
		LogLevel:   "info",
		Buffer:     100,
		Dir:        "/data/logs/mimi",
		SubRelPath: "2006-01/2006-01-02/2006-01-02-15.log",
	}
}

// NewLogger new logger
func NewLogger(logConfig LogConfig) *Logger {
	logLevel, err := ParseLevel(logConfig.LogLevel)
	if err != nil {
		logLevel = DebugLevel
	}

	ilogger := logrus.New()
	ilogger.Out = os.Stdout
	ilogger.SetLevel(logrus.Level(logLevel))

	mlog := NewMlog(&logConfig)
	return &Logger{
		ILogger: ilogger,
		Mlog:    mlog,
	}
}

// InitStatDir init stat dir
func InitStatDir(logTypes []string) {
	for _, logType := range logTypes {
		if err := GetMlog().EnsureStatDir(logType); err != nil {
			panic(err)
		}
	}
}

type Logger struct {
	ILogger
	Mlog *Mlog
}

// Flush flush
func (l *Logger) Flush() {
	l.Mlog.Flush()
}

// Close close
func (l *Logger) Close() {
	l.Mlog.Close()
}

var defaultLogger = NewLogger(DefaultLogConfig())

func SetLogger(logger *Logger) {
	defaultLogger.Close() // 关掉原来的fd防止上层重复调用引起携程泄露
	defaultLogger = logger
}

func GetLogger() *Logger {
	return defaultLogger
}

func GetMlog() *Mlog {
	return defaultLogger.Mlog
}

func Flush() {
	defaultLogger.Flush()
}

func Close() {
	defaultLogger.Close()
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	defaultLogger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Print(args ...interface{}) {
	defaultLogger.Print(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Warning(args ...interface{}) {
	defaultLogger.Warning(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Debugln(args ...interface{}) {
	defaultLogger.Debugln(args...)
}

func Infoln(args ...interface{}) {
	defaultLogger.Infoln(args...)
}

func Println(args ...interface{}) {
	defaultLogger.Println(args...)
}

func Warnln(args ...interface{}) {
	defaultLogger.Warnln(args...)
}

func Warningln(args ...interface{}) {
	defaultLogger.Warningln(args...)
}

func Errorln(args ...interface{}) {
	defaultLogger.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	defaultLogger.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	defaultLogger.Panicln(args...)
}
