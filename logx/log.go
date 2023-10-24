package logx

import (
	"go.uber.org/zap"
)

type Level int

const (
	LvlFatal Level = iota
	LvlPanic
	LvlDPanic
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
)

const (
	FormatJson    = "json"
	FormatConsole = "console"
	FormatText    = "text"
)

var _globalLogger *zapLogger

// Debug logs a message at DebugLevel.
func Debug(msg string, fields ...zap.Field) {
	_globalLogger.Logger().Debug(msg, fields...)
}

// Info logs a message at InfoLevel.
func Info(msg string, fields ...zap.Field) {
	_globalLogger.Logger().Info(msg, fields...)
}

// Warn logs a message at WarnLevel.
func Warn(msg string, fields ...zap.Field) {
	_globalLogger.Logger().Warn(msg, fields...)
}

// Error logs a message at ErrorLevel.
func Error(msg string, fields ...zap.Field) {
	_globalLogger.Logger().Error(msg, fields...)
}

// Panic logs a message at PanicLevel.
func Panic(msg string, fields ...zap.Field) {
	_globalLogger.Logger().Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel.
func Fatal(msg string, fields ...zap.Field) {
	_globalLogger.Logger().Fatal(msg, fields...)
}

func init() {
	config := NewDefaultConfig()
	lg, _ := InitZapLogger(config)
	_globalLogger = lg
}
