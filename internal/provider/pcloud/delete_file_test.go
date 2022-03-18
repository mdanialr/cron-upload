package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeleteFileUrl(t *testing.T) {
	const TOKEN = "token"
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Should has only one return's path",
			sample: "1234567890",
			expect: "https://eapi.pcloud.com/deletefile?auth=token&fileid=1234567890",
		},
		{
			name:   "Should return url as expected",
			sample: "0987654321",
			expect: "https://eapi.pcloud.com/deletefile?auth=token&fileid=0987654321",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := GetDeleteFileUrl(TOKEN, tc.sample)
			assert.Equal(t, tc.expect, out)
		})
	}
}
