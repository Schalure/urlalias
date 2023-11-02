package handlers

import (
	"fmt"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
)

type LoggerType int

const (
	LoggerTypeZap = LoggerType(iota)
	LoggerTypeSlog
)

type Loggerer interface {
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Close()
}

func NewLogger(loggerType LoggerType) (Loggerer, error) {

	switch loggerType {
	case LoggerTypeZap:
		return zaplogger.NewZapLogger("")
	default:
		return nil, fmt.Errorf("logger type is not supported")
	}
}
