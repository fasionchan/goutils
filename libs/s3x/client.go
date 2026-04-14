package s3x

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	EnvNameAwsAccessKey = "AWS_ACCESS_KEY"
	EnvNameAwsSecretKey = "AWS_SECRET_KEY"
	EnvNameAwsEndpoint  = "AWS_ENDPOINT"
	EnvNameAwsRegion    = "AWS_REGION"
)

func NewConfigFromDefaultEnv() aws.Config {
	return NewConfigFromEnvPro(EnvNameAwsAccessKey, EnvNameAwsSecretKey, EnvNameAwsEndpoint, EnvNameAwsRegion, os.Getenv)
}

func NewConfigFromDefaultEnvPro(getenv func(string) string) aws.Config {
	return NewConfigFromEnvPro(EnvNameAwsAccessKey, EnvNameAwsSecretKey, EnvNameAwsEndpoint, EnvNameAwsRegion, getenv)
}

func NewConfigFromEnv(
	accessKeyEnvName,
	secretKeyEnvName,
	endpointEnvName,
	regionEnvName string,
) aws.Config {
	return NewConfigFromEnvPro(accessKeyEnvName, secretKeyEnvName, endpointEnvName, regionEnvName, os.Getenv)
}

func NewConfigFromEnvPro(
	accessKeyEnvName,
	secretKeyEnvName,
	endpointEnvName,
	regionEnvName string,
	getenv func(string) string,
) aws.Config {
	accessKey := getenv(accessKeyEnvName)
	secretKey := getenv(secretKeyEnvName)
	region := getenv(regionEnvName)

	return aws.Config{
		Region:      region,
		Credentials: NewStaticCredentialsProvider(accessKey, secretKey),
	}
}

func NewConfig(
	accessKey,
	secretKey,
	region string,
) aws.Config {
	return aws.Config{
		Region:      region,
		Credentials: NewStaticCredentialsProvider(accessKey, secretKey),
	}
}

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

type BucketClient struct {
	*Client
	Bucket string
}

func NewBucketClient(cfg aws.Config, bucket string, opts ...func(*s3.Options)) *BucketClient {
	return &BucketClient{
		Client: NewClient(cfg, opts...),
		Bucket: bucket,
	}
}