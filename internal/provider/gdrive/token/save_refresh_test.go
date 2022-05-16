package token

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveRefreshToken(t *testing.T) {
	var sample = NewToken{
		TokenUrl: "https://example.com/oauth", RefreshToken: "refresh",
		ClientSecret: "secret", ClientID: "client",
	}

	testCases := []struct {
		name       string
		samplePath string
		wantErr    bool
	}{
		{
			name:       "Should pass when using correct and accessible file path",
			samplePath: "/tmp/test-credentials.json",
		},
		{
			name:       "Should fail when using invalid and or inaccessible file path",
			samplePath: "/fake/path/test-credentials.json",
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantErr {
			case false:
				require.NoError(t, SaveRefreshToken(tc.samplePath, sample))
			case true:
				require.Error(t, SaveRefreshToken(tc.samplePath, sample))
			}
		})
	}

	t.Cleanup(func() {
		for _, tc := range testCases {
			os.Remove(tc.samplePath)
		}
	})
}
