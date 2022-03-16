package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrettyJson(t *testing.T) {
	testCases := []struct {
		name    string
		sample  []byte
		expect  string
		wantErr bool
	}{
		{
			name:   "Should pass if using valid json string",
			sample: []byte(`{"name":"user","admin":"admin"}`),
			expect: `{
    "name": "user",
    "admin": "admin"
}`,
		},
		{
			name:    "Should error if using invalid json string",
			sample:  []byte(`{"name":"me`),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := PrettyJson(tc.sample)

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, out)
			case true:
				require.Error(t, err)
			}
		})
	}
}
