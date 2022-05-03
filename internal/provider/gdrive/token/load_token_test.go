package token

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadToken(t *testing.T) {
	const tokenPath = "/tmp/token-test.json"

	// prepare sample file auth.json in /tmp
	b := []byte(`{"key":"value"}`)
	require.NoError(t, os.WriteFile(tokenPath, b, 0644))

	t.Run("Should be success because loading token file that exist and has enough permission", func(t *testing.T) {
		_, err := LoadToken(tokenPath)
		require.NoError(t, err)
	})

	t.Run("Should error because loading token file that does not exist", func(t *testing.T) {
		_, err := LoadToken("/fake/path/auth.json")
		require.Error(t, err)
	})

	// cleanup and remove /tmp/auth.json file
	t.Cleanup(func() {
		os.Remove(tokenPath)
	})
}
