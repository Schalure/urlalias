package zaplogger

import (
	"fmt"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

func NewZapLogger(logFileName string) (*ZapLogger, error) {

	if logFileName != "" {
		panic("Save to log file not implemented")
	}

	aliasLogger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		return nil, fmt.Errorf("cannot initialize zap: %s", err)
	}
	suggarLogger := aliasLogger.Sugar()

	return &ZapLogger{
		logger:        aliasLogger,
		sugaredLogger: suggarLogger,
	}, nil
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args)
}

func (l *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Infow(msg, keysAndValues)
}

func (l *ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Errorw(msg, keysAndValues)
}

func (l *ZapLogger) Close() {
	l.logger.Sync()
}

func (l *ZapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Fatalw(msg, keysAndValues)
}
