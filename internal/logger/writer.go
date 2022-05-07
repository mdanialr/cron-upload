package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/mdanialr/cron-upload/internal/config"
)

var (
	// InfL info level logger instance that would be used throughout all this app.
	InfL *log.Logger
	// ErrL error level logger instance that would be used throughout all this app.
	ErrL *log.Logger
)

// InitLogger init and setup log file to write internal logger for this app.
func InitLogger(conf *config.Model) error {
	fl, err := os.OpenFile(conf.LogDir+"cron-upload-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		return fmt.Errorf("failed to open|create log file: %v\n", err)
	}

	InfL = log.New(fl, "[INFO] ", log.Ldate|log.Ltime)
	ErrL = log.New(fl, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
