package logger

import (
	"os"

	"github.com/iineva/git-nfs/pkg/signal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(name string) *zap.SugaredLogger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging := zapcore.Lock(_stdout)
	consoleErrors := zapcore.Lock(_stderr)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// consoleEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	l := zap.New(zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	))

	logger := l.Sugar().Named(name)
	signal.AddTermCallback(func(s os.Signal, done func()) {
		logger.Infof("[logger] receive signal (%v), closing", s)
		logger.Sync()
		done()
	})
	return logger
}
