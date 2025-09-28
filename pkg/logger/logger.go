package logger

import (
	"github.com/LiteMove/light-stack/internal/shared/config"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Init 初始化日志
func Init() {
	Log = logrus.New()

	cfg := config.Get()

	// 设置日志级别
	switch cfg.Log.Level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	if cfg.Log.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 设置输出
	if cfg.Log.Output == "file" {
		file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			Log.SetOutput(file)
		} else {
			Log.SetOutput(os.Stdout)
			Log.Warn("Failed to log to file, using default stderr")
		}
	} else {
		Log.SetOutput(os.Stdout)
	}
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return Log
}

// Debug 调试日志
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// WithField 带字段日志
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields 带多个字段日志
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
