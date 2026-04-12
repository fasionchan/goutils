package s3x

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	*s3.Client
}

func NewClient(cfg aws.Config, opts ...func(*s3.Options)) *Client {
	return &Client{
		Client: s3.NewFromConfig(cfg, opts...),
	}
}

type StaticCredentialProvider struct {
	AccessKey string `json:"AccessKey"`
	SecretKey string `json:"SecretKey"`
}

func NewStaticCredentialsProvider(accessKey, secretKey string) aws.CredentialsProvider {
	return &StaticCredentialProvider{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func (provider *StaticCredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     provider.AccessKey,
		SecretAccessKey: provider.SecretKey,
	}, nil
}
