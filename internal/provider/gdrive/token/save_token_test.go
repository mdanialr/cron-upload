package token

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestSaveToken(t *testing.T) {
	testCases := []struct {
		name       string
		sample     *oauth2.Token
		samplePath string
		wantErr    bool
	}{
		{
			name:       "Should pass because using correct and accessible file path",
			sample:     &oauth2.Token{RefreshToken: "refresh", AccessToken: "access", TokenType: "bearer"},
			samplePath: "/tmp/auth.json",
		},
		{
			name:       "Should error because using invalid and not enough permission file path",
			sample:     &oauth2.Token{RefreshToken: "refresh", AccessToken: "access", TokenType: "bearer"},
			samplePath: "/fake/path/auth.json",
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantErr {
			case false:
				require.NoError(t, SaveToken(tc.samplePath, tc.sample))
			case true:
				require.Error(t, SaveToken(tc.samplePath, tc.sample))
			}
		})
	}

	// cleanup and remove /tmp/auth.json
	t.Cleanup(func() {
		for _, tc := range testCases {
			os.Remove(tc.samplePath)
		}
	})
}
