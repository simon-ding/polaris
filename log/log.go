package log

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger
var atom zap.AtomicLevel

const dataPath = "./data"

func init() {
	atom = zap.NewAtomicLevel()
	atom.SetLevel(zap.DebugLevel)
	filer, _, err := zap.Open(filepath.Join(dataPath, "polaris.log"))
	if err != nil {
		panic(err)
	}
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	
	logger := zap.New(zapcore.NewCore(consoleEncoder, zapcore.Lock(filer), atom), zap.AddCallerSkip(1))

	sugar = logger.Sugar()

}



func SetLogLevel(l string) {
	switch strings.TrimSpace(strings.ToLower(l)) {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "info":
		atom.SetLevel(zap.InfoLevel)
	case "warn", "warnning":
		atom.SetLevel(zap.WarnLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
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
