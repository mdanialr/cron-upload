package logger

import (
	"io"
	"log"
)

// NewFile return logger that write the message to a file.
func NewFile(file io.Writer) Writer {
	return &LogFile{
		file: file,
	}
}

// LogFile log writer that write to a file.
type LogFile struct {
	file           io.Writer
	infLog, errLog *log.Logger
}

func (l *LogFile) Init() {
	l.infLog = log.New(l.file, "[INF] ", log.Ldate|log.Ltime)
	l.errLog = log.New(l.file, "[ERR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogFile) WriteInf(msg ...any) {
	l.infLog.Println(msg...)
}

func (l *LogFile) WriteErr(msg ...any) {
	l.errLog.Println(msg...)
}
