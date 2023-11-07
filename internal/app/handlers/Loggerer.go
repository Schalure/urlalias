package handlers

import (
	"fmt"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
)

type LoggerType int

const (
	LoggerTypeZap = iota
	LoggerTypeSlog
)

func (l LoggerType) String() string {
	return [...]string{"LoggerTypeZap", "LoggerTypeSlog"}[l]
}

type Loggerer interface {
	Info(args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
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
