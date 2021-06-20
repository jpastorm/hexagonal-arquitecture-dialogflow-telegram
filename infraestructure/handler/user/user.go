package user

type Logger interface {
	Warningf(format string, arg ...interface{})
	Errorf(format string, arg ...interface{})
	Tracef(format string, arg ...interface{})
	Infof(format string, arg ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}