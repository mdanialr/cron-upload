package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Init return ready to use AWS S3 bucket client that use the given IAM account
// token path as the credential also optionally set the region.
func Init(tokenPath string, region string) (*s3.Client, error) {
	// prepare the local config to retrieve the required values
	conf, err := newCredential(tokenPath)
	if err != nil {
		return nil, err
	}
	if err = conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate s3 bucket credential file: %s", err)
	}
	// create s3 config then init s3 client using it
	cfg, err := config.LoadDefaultConfig(context.Background(), func(opt *config.LoadOptions) error {
		opt.Credentials = credentials.NewStaticCredentialsProvider(conf.Key, conf.Secret, "")
		opt.Region = region
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init aws config: %s", err)
	}
	return s3.NewFromConfig(cfg), nil
}
