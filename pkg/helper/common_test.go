package helper_test

import (
	"testing"
	"time"

	h "github.com/mdanialr/cron-upload/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestToWib(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func() time.Time
		expect string
	}{
		{
			name: "Given date 2023-01-15 08:16:25 UTC should return 2023-01-15 15:16:25 WIB",
			setup: func() time.Time {
				return time.Date(2023, 01, 15, 8, 16, 25, 0, time.UTC)
			},
			expect: "2023-01-15 15:16:25 WIB",
		},
		{
			name: "Given date 2023-01-15 08:16:25 UTC+9 should return 2023-01-15 06:16:25 WIB",
			setup: func() time.Time {
				jpn, _ := time.LoadLocation("Asia/Tokyo")
				return time.Date(2023, 01, 15, 8, 16, 25, 0, jpn)
			},
			expect: "2023-01-15 06:16:25 WIB",
		},
		{
			name: "Given date 2023-01-15 00:16:25 UTC+8 should return 2023-01-14 23:16:25 WIB",
			setup: func() time.Time {
				sgp, _ := time.LoadLocation("Asia/Singapore")
				return time.Date(2023, 01, 15, 0, 16, 25, 0, sgp)
			},
			expect: "2023-01-14 23:16:25 WIB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatTime := "2006-01-02 15:04:05 MST"
			assert.Equal(t, tc.expect, h.ToWib(tc.setup()).Format(formatTime))
		})
	}
}
