package provider_test

import (
	"testing"

	"github.com/mdanialr/cron-upload/internal/provider"
	"github.com/stretchr/testify/assert"
)

var samplePayloads = []*provider.Payload{
	{
		Id:     "12",
		Name:   "Root",
		Parent: []string{},
	},
	{
		Id:     "2",
		Name:   "DB",
		Parent: []string{"12"},
	},
	{
		Id:     "3",
		Name:   "app",
		Parent: []string{"12"},
	},
	{
		Id:     "24",
		Name:   "sample",
		Parent: []string{"3"},
	},
}

func TestLookupRoute(t *testing.T) {
	testCases := []struct {
		name        string
		sampleRoute string
		expect      string
	}{
		{
			name:        "Given route 'Root' without parent should return 'Root' without any added route",
			sampleRoute: "Root",
			expect:      "Root",
		},
		{
			name:        "Given route that does not exist yet in the stack should return just like route without parent",
			sampleRoute: "Backup",
			expect:      "Backup",
		},
		{
			name:        "Given route 'DB' that has single parent 'Root' should return 'Root/DB'",
			sampleRoute: "DB",
			expect:      "Root/DB",
		},
		{
			name:        "Given route 'sample' that has nested parent should return 'Root/app/sample'",
			sampleRoute: "sample",
			expect:      "Root/app/sample",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := provider.LookupRoute(samplePayloads, tc.sampleRoute)
			assert.Equal(t, tc.expect, res)
		})
	}
}

func TestLookupRouteName(t *testing.T) {
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Given want to search route id '3' and exist in stack should return 'app'",
			sample: "3",
			expect: "app",
		},
		{
			name:   "Given want to search route id '5' and does not exist in stack should return empty string",
			sample: "5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := provider.LookupRouteName(samplePayloads, tc.sample)
			assert.Equal(t, tc.expect, res)
		})
	}
}
