package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// UploadModel data structure of the upload data in config file.
type UploadModel struct {
	Path   string `mapstructure:"path"`
	Name   string `mapstructure:"name"`
	Retain uint   `mapstructure:"retain"`
}

// Init return new viper instance with the given filepath as
// the directory where the config file is.
func Init(filepath string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(filepath)
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

// Sanitize set default value for necessary fields.
func Sanitize(v *viper.Viper) {
	v.SetDefault("log", "/tmp")
	v.SetDefault("retain", 8640)
	v.SetDefault("worker", 2)
}

// Validate do validate for all required fields.
func Validate(v *viper.Viper) error {
	// make sure the root field is not empty
	if v.GetString("root") == "" {
		return fmt.Errorf("`root` field in config file is required")
	}
	// make sure the worker has minimum value of 1
	if v.GetInt("worker") < 1 {
		return fmt.Errorf("`worker` field should has minimum value of 1. if there is no worker then who on earth that will do the job")
	}
	// make sure provider name & auth is not empty
	if v.GetString("provider.name") == "" {
		return fmt.Errorf("`provider.name` field in config file is required")
	}
	if v.GetString("provider.cred") == "" {
		return fmt.Errorf("`provider.cred` field in config file is required")
	}
	// make sure the chunk size is positive number if provided
	if v.GetInt("chunk") < 0 {
		return fmt.Errorf("please provide positive number for the upload chunk size in field `chunk`")
	}
	return nil
}

// GetUploads return all data that will be uploaded to provider.
func GetUploads(v *viper.Viper) []UploadModel {
	var res []UploadModel
	v.UnmarshalKey("upload", &res)
	return res
}

// GetRetainExpiry get default retain expiry from config if the given retain
// value is zero.
func GetRetainExpiry(v *viper.Viper, retain uint) uint {
	if retain == 0 {
		return v.GetUint("retain")
	}
	return retain
}
