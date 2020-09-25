package logger

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	Info(message string)
	Debug(message string)
}

type NullLogger struct {
}

func (nl NullLogger) Log(v ...interface{}) {
	return
}

func (nl NullLogger) Logf(format string, v ...interface{}) {
	return
}

func (nl NullLogger) Info(message string) {
	return
}

func (nl NullLogger) Debug(message string) {
	return
}
