package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFolder_Sanitization(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Folder
		expect  Folder
		wantErr bool
	}{
		{
			name:    "Should be failed when not providing required `path` field",
			sample:  Folder{Name: "test/dir"},
			wantErr: true,
		},
		{
			name:   "Should has same value as `path` field if `name` field is not provided",
			sample: Folder{Path: "/full/path/upload/folders"},
			expect: Folder{Path: "/full/path/upload/folders", Name: "full/path/upload/folders", Retain: uint(0)},
		},
		{
			name:   "Should pass when providing required field",
			sample: Folder{Name: "/test/dir/", Path: "full/path/upload/folders"},
			expect: Folder{Name: "test/dir", Path: "full/path/upload/folders", Retain: uint(0)},
		},
		{
			name:   "Default value of `retain_days` field should be '0' if not provided",
			sample: Folder{Name: "/app", Path: "/path/folders"},
			expect: Folder{Name: "app", Path: "/path/folders", Retain: uint(0)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, tc.sample)
			}
		})
	}
}

func TestUpload_Sanitization(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Upload
		expect  Upload
		wantErr bool
	}{
		{
			name: "Should error if not providing required `path` field",
			sample: Upload{
				{Folders: Folder{}},
			},
			wantErr: true,
		},
		{
			name: "Should pass if providing minimum required `path` field",
			sample: Upload{
				{Folders: Folder{Path: "/full/path/folders"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}
