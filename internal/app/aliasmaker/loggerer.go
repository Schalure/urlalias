package aliasmaker

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
