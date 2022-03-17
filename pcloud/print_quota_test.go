package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQuotaUrl(t *testing.T) {
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Should has only one result's path",
			sample: "token",
			expect: "https://eapi.pcloud.com/userinfo?auth=token",
		},
		{
			name:   "Should return url value as expected",
			sample: "secret",
			expect: "https://eapi.pcloud.com/userinfo?auth=secret",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := GetQuotaUrl(tc.sample)
			assert.Equal(t, tc.expect, out)
		})
	}
}
