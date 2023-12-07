package aliaslogger

import (
	"fmt"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
)

type LoggerType int

const (
	LoggerTypeZap = iota
	LoggerTypeSlog
)

func (l LoggerType) String() string {
	return [...]string{"LoggerTypeZap", "LoggerTypeSlog"}[l]
}

// --------------------------------------------------
//
//	Choose logger for service
func NewLogger(loggerType LoggerType) (aliasmaker.Loggerer, error) {
	switch loggerType {
	case LoggerTypeZap:
		return zaplogger.NewZapLogger("")
	default:
		return nil, fmt.Errorf("logger type is not supported: %s", loggerType.String())
	}
}
