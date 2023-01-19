package provider_test

import (
	"testing"

	"github.com/mdanialr/cron-upload/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestValidateSupportedClouds(t *testing.T) {
	testCases := []struct {
		name          string
		sample        string
		wantErr       bool
		containErrMsg string
	}{
		{
			name:   "Add support for Google Drive as 'drive'",
			sample: "drive",
		},
		{
			name:          "Give error for unsupported cloud provider and contain error message 'is not supported'",
			sample:        "s3",
			wantErr:       true,
			containErrMsg: "is not supported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := provider.ValidateSupportedClouds(tc.sample)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.containErrMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}
