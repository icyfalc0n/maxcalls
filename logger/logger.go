package logger

type Logger interface {
	Debugf(format string, args ...interface{})
}
