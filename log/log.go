package log

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger
var atom zap.AtomicLevel

const dataPath = "./data"

func InitLogger(toFile bool) {
	atom = zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)

	w := zapcore.Lock(os.Stdout)
	if toFile {
		w = zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dataPath, "logs", "polaris.log"),
			MaxSize:    50, // megabytes
			MaxBackups: 3,
			MaxAge:     30, // days
			Compress:   true,
		})
	
	}

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	logger := zap.New(zapcore.NewCore(consoleEncoder, w, atom), zap.AddCallerSkip(1),zap.AddCaller())

	sugar = logger.Sugar()

}

func SetLogLevel(l string) {
	switch strings.TrimSpace(strings.ToLower(l)) {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
		Debug("set log level to debug")
	case "info":
		atom.SetLevel(zap.InfoLevel)
		Info("set log level to info")
	case "warn", "warning":
		atom.SetLevel(zap.WarnLevel)
		Warn("set log level to warning")
	case "error":
		atom.SetLevel(zap.ErrorLevel)
		Error("set log level to error")
	}
}

func Logger() *zap.SugaredLogger {
	return sugar
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}
