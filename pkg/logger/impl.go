package logger

func Debug(msg string, options ...interface{}) {
	regLogger.Debugw(msg, options...)
}

func Info(msg string, options ...interface{}) {
	regLogger.Infow(msg, options...)
}

func Warn(msg string, options ...interface{}) {
	regLogger.Warnw(msg, options...)
}

func Error(msg string, options ...interface{}) {
	regLogger.Errorw(msg, options...)
}

func Panic(msg string, options ...interface{}) {
	regLogger.Panicw(msg, options...)
}

func Fatal(msg string, options ...interface{}) {
	regLogger.Fatalw(msg, options...)
}
