package config_test

import (
	"testing"

	"github.com/mdanialr/cron-upload/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	// 1st case for `log`
	testCasesLog := []struct {
		name   string
		setup  func() *viper.Viper
		expect string
	}{
		{
			name: "Given log is empty should has /tmp as the default value",
			setup: func() *viper.Viper {
				return viper.New()
			},
			expect: "/tmp",
		},
		{
			name: "Given log is /my/log/path should has /my/log/path as the value",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("log", "/my/log/path")
				return v
			},
			expect: "/my/log/path",
		},
	}
	for _, tc := range testCasesLog {
		t.Run(tc.name, func(t *testing.T) {
			conf := tc.setup()
			config.Sanitize(conf)
			assert.Equal(t, tc.expect, conf.GetString("log"))
		})
	}

	// 2nd case for `retain`
	testCasesRetain := []struct {
		name   string
		setup  func() *viper.Viper
		expect uint
	}{
		{
			name: "Given retain is empty should has 8640 minutes as the default value",
			setup: func() *viper.Viper {
				return viper.New()
			},
			expect: 8640,
		},
		{
			name: "Given retain is 60 minutes should has 60 minutes as the value",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("retain", 60)
				return v
			},
			expect: 60,
		},
	}
	for _, tc := range testCasesRetain {
		t.Run(tc.name, func(t *testing.T) {
			conf := tc.setup()
			config.Sanitize(conf)
			assert.Equal(t, tc.expect, conf.GetUint("retain"))
		})
	}

	// 3rd case for `worker`
	testCasesWorker := []struct {
		name   string
		setup  func() *viper.Viper
		expect uint
	}{
		{
			name: "Given worker is empty should has 2 workers as the default value",
			setup: func() *viper.Viper {
				return viper.New()
			},
			expect: 2,
		},
		{
			name: "Given worker is 8 workers should has 8 workers as the value",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("worker", 8)
				return v
			},
			expect: 8,
		},
	}
	for _, tc := range testCasesWorker {
		t.Run(tc.name, func(t *testing.T) {
			conf := tc.setup()
			config.Sanitize(conf)
			assert.Equal(t, tc.expect, conf.GetUint("worker"))
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name          string
		setup         func() *viper.Viper
		wantErr       bool
		containErrMsg string
	}{
		{
			name: "Given root is empty should throw error and contain message `is required`",
			setup: func() *viper.Viper {
				return viper.New()
			},
			wantErr:       true,
			containErrMsg: "is required",
		},
		{
			name: "Given root is 'Backup' and worker is '0' should throw error and contain message" +
				" `has minimum value of 1`",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("root", "Backup")
				return v
			},
			wantErr:       true,
			containErrMsg: "has minimum value of 1",
		},
		{
			name: "Given root is 'Backup', worker is '5' and provider.name is empty should throw error and" +
				" contain message `is required`",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("root", "Backup")
				v.Set("worker", 5)
				return v
			},
			wantErr:       true,
			containErrMsg: "is required",
		},
		{
			name: "Given root is 'Backup', worker is '5', provider.name is 'hi' and provider.cred is empty" +
				" should throw error and contain message `is required`",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("root", "Backup")
				v.Set("worker", 5)
				v.Set("provider.name", "hi")
				return v
			},
			wantErr:       true,
			containErrMsg: "is required",
		},
		{
			name: "Given root is 'Backup', worker is '5', provider.name is 'hi', provider.cred is /path/to/cred.json" +
				" and chunk is '-1' should throw error and contain message `provide positive number`",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("root", "Backup")
				v.Set("worker", 5)
				v.Set("provider.name", "hi")
				v.Set("provider.cred", "/path/to/cred.json")
				v.Set("chunk", -1)
				return v
			},
			wantErr:       true,
			containErrMsg: "provide positive number",
		},
		{
			name: "Given root is 'Backup', worker is '5', provider.name is 'hi', provider.cred is /path/to/cred.json" +
				" and chunk is empty should has no error",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("root", "Backup")
				v.Set("worker", 5)
				v.Set("provider.name", "hi")
				v.Set("provider.cred", "/path/to/cred.json")
				return v
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := tc.setup()
			err := config.Validate(conf)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.containErrMsg)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetRetainExpiry(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func() *viper.Viper
		sample uint
		expect uint
	}{
		{
			name: "Given retain is 0 and viper with retain value of 5 should return 5",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("retain", 5)
				return v
			},
			sample: 0,
			expect: 5,
		},
		{
			name: "Given retain is 4 and viper with retain value of 5 should return 4",
			setup: func() *viper.Viper {
				v := viper.New()
				v.Set("retain", 5)
				return v
			},
			sample: 4,
			expect: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := config.GetRetainExpiry(tc.setup(), tc.sample)
			assert.Equal(t, tc.expect, res)
		})
	}
}
