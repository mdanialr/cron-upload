package scan_test

import (
	"testing"

	"github.com/mdanialr/cron-upload/pkg/scan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiles(t *testing.T) {
	testCases := []struct {
		name      string
		sampleDir string
		wantErr   bool
	}{
		{
			name:      "Should error when scanning fake directory",
			sampleDir: "path/to/fake/dir",
			wantErr:   true,
		},
		{
			name:      "Should pass when scanning real directory",
			sampleDir: "testdata",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := scan.FilesAsc(tc.sampleDir)

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}

	t.Run("Scanning testdata directory should has exactly three files", func(t *testing.T) {
		out, err := scan.FilesAsc("testdata")
		require.NoError(t, err)
		assert.Equal(t, 3, len(out))
	})

	t.Run("Scanning testdata directory should ignore directory and indexing only the files", func(t *testing.T) {
		out, err := scan.FilesAsc("testdata")
		require.NoError(t, err)
		assert.Equal(t, 3, len(out))
	})
}
