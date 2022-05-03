package logger

import (
	"testing"

	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	testCases := []struct {
		name       string
		sampleConf config.Model
		wantErr    bool
	}{
		{
			name:       "Should be pass when using valid log dir",
			sampleConf: config.Model{Provider: config.Provider{Name: "drive"}},
		},
		{
			name:       "Should be failed when using fake log dir",
			sampleConf: config.Model{Provider: config.Provider{Name: "drive"}, LogDir: "/fake/log/dir"},
			wantErr:    true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sampleConf.Sanitization()

			switch tt.wantErr {
			case false:
				require.NoError(t, err)
				require.NoError(t, InitLogger(&tt.sampleConf))
			case true:
				require.Error(t, InitLogger(&tt.sampleConf))
			}
		})
	}
}
