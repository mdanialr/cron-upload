package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUploadFileUrl(t *testing.T) {
	const TOKEN = "token"
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Should has only one return's path",
			sample: "1234567890",
			expect: "https://eapi.pcloud.com/uploadfile?auth=token&folderid=1234567890&nopartial=1&renameifexists=1",
		},
		{
			name:   "Should return url as expected",
			sample: "0987654321",
			expect: "https://eapi.pcloud.com/uploadfile?auth=token&folderid=0987654321&nopartial=1&renameifexists=1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := GetUploadFileUrl(TOKEN, tc.sample)
			assert.Equal(t, tc.expect, out)
		})
	}
}
