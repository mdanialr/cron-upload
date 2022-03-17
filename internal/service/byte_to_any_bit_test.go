package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBytesToAnyBit(t *testing.T) {
	testCases := []struct {
		name    string
		sample  int64
		unit    string
		expect  string
		wantErr bool
	}{
		{
			name:   "Should return 6Gb when using the equivalent of the bytes",
			sample: 6442450944, unit: "Gb", expect: "6Gb",
		},
		{
			name:   "Should return 6Mb when using the equivalent of the bytes",
			sample: 6291456, unit: "Mb", expect: "6Mb",
		},
		{
			name:   "Should return 6Kb when using the equivalent of the bytes",
			sample: 6144, unit: "Kb", expect: "6Kb",
		},
		{
			name:   "Should error when using unit that does not supported",
			sample: 1024, unit: "Tb",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := BytesToAnyBit(tc.sample, tc.unit)

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
