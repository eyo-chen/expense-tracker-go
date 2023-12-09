package logger

// Logger is the interface that wraps the basic logging methods
type Logger interface {
	Debugw(msg string, options ...interface{})
	Infow(msg string, options ...interface{})
	Warnw(msg string, options ...interface{})
	Errorw(msg string, options ...interface{})
	Panicw(msg string, options ...interface{})
	Fatalw(msg string, options ...interface{})
}
