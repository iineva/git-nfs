package logger

import (
	"os"

	"github.com/iineva/git-nfs/pkg/signal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_debug      = true
	_baseLogger *zap.SugaredLogger
)

func Debug(debug bool) {
	_debug = debug
}

func initBaseLogger() *zap.SugaredLogger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging := zapcore.Lock(_stdout)
	consoleErrors := zapcore.Lock(_stderr)
	var consoleEncoder zapcore.Encoder
	if _debug {
		consoleEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		consoleEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	return zap.New(zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)).Sugar()
}

func New(name string) *zap.SugaredLogger {
	if _baseLogger == nil {
		_baseLogger = initBaseLogger()
	}
	logger := _baseLogger.Named(name)
	signal.AddTermCallback(func(s os.Signal, done func()) {
		logger.Infof("receive signal (%v), closing", s)
		logger.Sync()
		done()
	})
	return logger
}
