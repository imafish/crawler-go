package main

// Logger interface provides common methods for logging
type Logger interface {
	Tracef(format string, a ...interface{}) (n int, err error)
	Debugf(format string, a ...interface{}) (n int, err error)
	Infof(format string, a ...interface{}) (n int, err error)
	Warningf(format string, a ...interface{}) (n int, err error)
	Errorf(format string, a ...interface{}) (n int, err error)
	Fatalf(format string, a ...interface{}) (n int, err error)

	Trace(a ...interface{}) (n int, err error)
	Debug(a ...interface{}) (n int, err error)
	Info(a ...interface{}) (n int, err error)
	Warning(a ...interface{}) (n int, err error)
	Error(a ...interface{}) (n int, err error)
	Fatal(a ...interface{}) (n int, err error)

	SetPrefix(prefix string) error
}
