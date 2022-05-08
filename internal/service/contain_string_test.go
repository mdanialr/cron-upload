package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		name   string
		sample []string
		check  string
		expect bool
	}{
		{
			name:   "Should true if `sample` contains string `check`",
			sample: []string{"one.zip", "two.txt", "three.mp4"},
			check:  "two.txt",
			expect: true,
		},
		{
			name:   "Should false if `sample` does not contains string `check`",
			sample: []string{"one.zip", "two.txt", "three.mp4"},
			check:  "one.txt",
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := Contains(tc.sample, tc.check)
			assert.Equal(t, tc.expect, out)
		})
	}
}
