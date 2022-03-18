package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClearTrashUrl(t *testing.T) {
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Should has only one return's path",
			sample: "token",
			expect: "https://eapi.pcloud.com/trash_clear?auth=token&folderid=0",
		},
		{
			name:   "Should return url value as expected",
			sample: "secret",
			expect: "https://eapi.pcloud.com/trash_clear?auth=secret&folderid=0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := GetClearTrashUrl(tc.sample)
			assert.Equal(t, tc.expect, out)
		})
	}
}
