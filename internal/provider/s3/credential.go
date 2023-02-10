package s3

import (
	"encoding/json"
	"fmt"
	"os"
)

// defaultCredential default structure of the credential json file used by
// this app to integrate with AWS S3 bucket service.
type defaultCredential struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func (c *defaultCredential) validate() error {
	if c.Key == "" {
		return fmt.Errorf("`key` is required")
	}
	if c.Secret == "" {
		return fmt.Errorf("`secret` is required")
	}
	return nil
}

// newCredential return ready to use credential that's necessary to init
// s3 bucket service.
func newCredential(filePath string) (*defaultCredential, error) {
	var awsConfig defaultCredential
	dt, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open s3 credential file: %s", err)
	}
	json.Unmarshal(dt, &awsConfig)
	return &awsConfig, nil
}
