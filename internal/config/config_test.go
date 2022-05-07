package config

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeReader just fake type to satisfies io.Reader interfaces so it could
// trigger error buffer read from.
type fakeReader struct{}

func (_ fakeReader) Read(_ []byte) (_ int, _ error) {
	return 0, fmt.Errorf("this should trigger error in test")
}

func TestNewConfig(t *testing.T) {
	fakeConfigFile := `log: /tmp/cron-upload-log`
	buf := bytes.NewBufferString(fakeConfigFile)
	t.Run("Should pass when using valid value", func(t *testing.T) {
		mod, err := NewConfig(buf)
		require.NoError(t, err)

		assert.Equal(t, "/tmp/cron-upload-log", mod.LogDir)
	})

	fakeConfigFile = `provider: 101`
	buf = bytes.NewBufferString(fakeConfigFile)
	t.Run("Should error when using mismatch type in yaml unmarshalling", func(t *testing.T) {
		_, err := NewConfig(buf)
		require.Error(t, err)
	})

	t.Run("Injecting fake reader should be error in buffer read from", func(t *testing.T) {
		_, err := NewConfig(fakeReader{})
		require.Error(t, err)
	})
}

func TestModel_Sanitization(t *testing.T) {
	testCases := []struct {
		name    string
		sample  Model
		expect  Model
		wantErr bool
	}{
		{
			name:   "Should has default of '/tmp/' for log dir if not provided",
			sample: Model{Provider: Provider{Name: "cloud"}},
			expect: Model{LogDir: "/tmp/", Provider: Provider{Name: "cloud"}, RootFolder: "Cron-Backups"},
		},
		{
			name:   "Log dir should has prefix and trailing slash '/'",
			sample: Model{LogDir: "log/dir", Provider: Provider{Name: "cloud"}},
			expect: Model{LogDir: "/log/dir/", Provider: Provider{Name: "cloud"}, RootFolder: "Cron-Backups"},
		},
		{
			name:   "Should has default value of '1' for max worker if not provided",
			sample: Model{Provider: Provider{Name: "cloud"}},
			expect: Model{LogDir: "/tmp/", Provider: Provider{Name: "cloud"}, RootFolder: "Cron-Backups"},
		},
		{
			name:    "Should error if required `provider's fields` not provided",
			sample:  Model{},
			wantErr: true,
		},
		{
			name:    "Should error if `provider.name` is 'drive' but `auth` field is not provided",
			sample:  Model{Provider: Provider{Name: "drive"}},
			wantErr: true,
		},
		{
			name:   "Should has leading and trailing slash for `auth` field if `provider.name` is 'drive'",
			sample: Model{Provider: Provider{Name: "drive", Auth: "auth", Token: "dir"}},
			expect: Model{LogDir: "/tmp/", RootFolder: "Cron-Backups",
				Provider: Provider{Name: "drive", Auth: "auth", Token: "/dir/cron-upload-token.json"},
			},
		},
		{
			name:   "Should has default value of '/tmp/cron-upload-token.json' if `provider.name` is 'drive' and `token` field is not provided",
			sample: Model{Provider: Provider{Name: "drive", Auth: "auth"}},
			expect: Model{LogDir: "/tmp/", RootFolder: "Cron-Backups",
				Provider: Provider{Name: "drive", Auth: "auth", Token: "/tmp/cron-upload-token.json"},
			},
		},
		{
			name:   "Should has default value of 'Cron-Backups' if `root_folder` field is not provided",
			sample: Model{Provider: Provider{Name: "cloud"}},
			expect: Model{LogDir: "/tmp/", RootFolder: "Cron-Backups",
				Provider: Provider{Name: "cloud"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, tc.sample)
			case true:
				require.Error(t, err)
			}
		})
	}
}
