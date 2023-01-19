package logger

// Writer interface to write log message.
type Writer interface {
	// Init do some necessary setup before calling WriteInf & WriteErr
	Init()
	// WriteInf write message for info log
	WriteInf(...any)
	// WriteErr write message for error log
	WriteErr(...any)
}
