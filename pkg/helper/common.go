package helper

import (
	"fmt"
	"time"
)

// primitiveType general type that contain all primitive/built in types.
type primitiveType interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string | bool
}

// Def return the default value of the given t if nil otherwise the real value
// instead.
func Def[T primitiveType](t *T) (newT T) {
	if t != nil {
		return *t
	}
	return newT
}

// Ptr return pointer of the given t.
func Ptr[T primitiveType](t T) *T {
	return &t
}

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
