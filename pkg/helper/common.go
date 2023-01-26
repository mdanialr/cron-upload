package helper

import (
	"fmt"
	"time"
)

// LogStart append [BEGN] and the prefix to the given message and return it.
// Also append newline at the end.
func LogStart(prefix, message string) string {
	return fmt.Sprintf("[BEGN] | %s | %s", prefix, message)
}

// LogDone append [DONE] and the prefix to the given message and return it.
// Also append newline at the end.
func LogDone(prefix, message string) string {
	return fmt.Sprintf("[DONE] | %s | %s", prefix, message)
}

// ToWib convert the given time to Waktu Indonesia Barat which is time zone
// UTC+7.
func ToWib(t time.Time) time.Time {
	wib, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(wib)
}
