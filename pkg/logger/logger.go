package logger

type Logger interface {
	Info(args ...any)
	// Infof logs a message at level Info on the standard logger.
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Debug(args ...any)
	Debugf(format string, args ...any)
}
