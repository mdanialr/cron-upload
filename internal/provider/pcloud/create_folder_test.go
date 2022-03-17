package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCreateFolderUrl(t *testing.T) {
	const TOKEN = "token"
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Should has only one result's path",
			sample: "/vps",
			expect: "https://eapi.pcloud.com/createfolderifnotexists?auth=token&path=/vps",
		},
		{
			name:   "Should return url as expected",
			sample: "/vps/backup",
			expect: "https://eapi.pcloud.com/createfolderifnotexists?auth=token&path=/vps/backup",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := GetCreateFolderUrl(TOKEN, tc.sample)
			assert.Equal(t, tc.expect, out)
		})
	}
}
