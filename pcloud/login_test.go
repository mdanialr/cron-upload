package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_GenerateTokenUrl(t *testing.T) {
	testCases := []struct {
		name   string
		sample User
		digest string
		expect string
	}{
		{
			name:   "1# Should pass",
			sample: User{Username: "user", Password: "pass"},
			digest: "digest",
			expect: "https://eapi.pcloud.com/userinfo?digest=digest&getauth=1&passworddigest=55f3b3ce9973588105765465d3a8b45613d426a3&username=user",
		},
		{
			name:   "2# Should pass",
			sample: User{Username: "admin", Password: "admin"},
			digest: "digest",
			expect: "https://eapi.pcloud.com/userinfo?digest=digest&getauth=1&passworddigest=37fb256ed24ce6677dc53fdd9aabe412fe35e248&username=admin",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.sample.GenerateTokenUrl(tc.digest)
			assert.Equal(t, tc.expect, out)
		})
	}
}
