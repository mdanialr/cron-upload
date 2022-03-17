package pcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDigestUrl(t *testing.T) {
	t.Run("Should be only one path that always generate same result", func(t *testing.T) {
		expect := "https://eapi.pcloud.com/getdigest"
		assert.Equal(t, expect, GetDigestUrl())
	})
}
