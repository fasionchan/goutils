package s3x

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	EnvNameAwsAccessKey = "AWS_ACCESS_KEY"
	EnvNameAwsSecretKey = "AWS_SECRET_KEY"
	EnvNameAwsEndpoint  = "AWS_ENDPOINT"
	EnvNameAwsRegion    = "AWS_REGION"
)

func NewConfigFromDefaultEnv() aws.Config {
	return NewConfigFromEnvPro(EnvNameAwsAccessKey, EnvNameAwsSecretKey, EnvNameAwsRegion, os.Getenv)
}

func NewConfigFromDefaultEnvPro(getenv func(string) string) aws.Config {
	return NewConfigFromEnvPro(EnvNameAwsAccessKey, EnvNameAwsSecretKey, EnvNameAwsRegion, getenv)
}

func NewConfigFromEnv(
	accessKeyEnvName,
	secretKeyEnvName,
	regionEnvName string,
) aws.Config {
	return NewConfigFromEnvPro(accessKeyEnvName, secretKeyEnvName, regionEnvName, os.Getenv)
}

func NewConfigFromEnvPro(
	accessKeyEnvName,
	secretKeyEnvName,
	regionEnvName string,
	getenv func(string) string,
) aws.Config {
	return NewConfig(
		getenv(accessKeyEnvName),
		getenv(secretKeyEnvName),
		getenv(regionEnvName),
	)
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
	config *Config
}

func NewClient(cfg aws.Config, opts ...func(*s3.Options)) *Client {
	return &Client{
		Client: s3.NewFromConfig(cfg, opts...),
	}
}

func (client *Client) GetConfig() *Config {
	return client.config
}

func (client *Client) TransferManager(opts ...func(*transfermanager.Options)) *transfermanager.Client {
	return transfermanager.New(client, opts...)
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