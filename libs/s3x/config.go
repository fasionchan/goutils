package s3x

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	S3QueryNameAccessKey = "accessKey"
	S3QueryNameSecretKey = "secretKey"
	S3QueryNameBucket    = "bucket"
	S3QueryNameRegion    = "region"
	S3QueryNamePathStyle = "pathStyle"
	S3QueryNameBlockSize = "blockSize"
	S3QueryNameAcl       = "acl"
	S3QueryNameKeyPrefix = "keyPrefix"

	DefaultS3Region = "default"

	MinS3BlockSize  = 5 * 1024 * 1024

	DefaultS3Acl = "private"
)

type Config struct {
	Endpoint string
	EndpointHost string
	EndpointUrl *url.URL
	Region   string
	AccessKey string
	SecretKey string
	Bucket    string
	KeyPrefix string
	PathStyle bool
	BlockSize int
	Acl       string
}

func ParseConfigFromUrl(url_ string) (*Config, error) {
	parsedURL, err := url.Parse(url_)
	if err != nil {
		return nil, fmt.Errorf("S3ParseUrlFailed, err: %w", err)
	}

	query := parsedURL.Query()

	accessKey := query.Get(S3QueryNameAccessKey)
	if accessKey == "" {
		return nil, fmt.Errorf("S3AccessKeyNil")
	}

	secretKey := query.Get(S3QueryNameSecretKey)
	if secretKey == "" {
		return nil, fmt.Errorf("S3SecretKeyNil")
	}

	bucket := query.Get(S3QueryNameBucket)

	keyPrefix := query.Get(S3QueryNameKeyPrefix)

	region := query.Get(S3QueryNameRegion)
	if region == "" {
		region = DefaultS3Region
	}

	var pathStyle bool
	pathStyleStr := query.Get(S3QueryNamePathStyle)
	if pathStyleStr != "" {
		var err error
		pathStyle, err = strconv.ParseBool(pathStyleStr)
		if err != nil {
			return nil, fmt.Errorf("S3PathStyleParseFailed, err: %w", err)
		}
	}

	var blockSize int64
	blockSizeStr := query.Get(S3QueryNameBlockSize)
	if blockSizeStr != "" {
		var err error
		blockSize, err = strconv.ParseInt(blockSizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("S3BlockSizeParseFailed, err: %w", err)
		}
		if blockSize <= MinS3BlockSize {
			blockSize = MinS3BlockSize
		}
	}

	acl := query.Get(S3QueryNameAcl)
	if acl == "" {
		acl = DefaultS3Acl
	}

	parsedURL.RawQuery = ""

	return &Config{
		Endpoint: parsedURL.String(),
		EndpointHost: parsedURL.Host,
		EndpointUrl: parsedURL,
		Region:   region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		KeyPrefix: keyPrefix,
		PathStyle: pathStyle,
		BlockSize: int(blockSize),
		Acl:       acl,
	}, nil
}

func (config *Config) ApplyToOptions(options *s3.Options) {
	options.BaseEndpoint = aws.String(config.Endpoint)
}

func (config *Config) AwsConfig() aws.Config {
	return aws.Config{
		Region:      config.Region,
		Credentials: NewStaticCredentialsProvider(config.AccessKey, config.SecretKey),
	}
}

func (config *Config) GetEndpoint() string {
	if config == nil {
		return ""
	}
	return config.Endpoint
}

func (config *Config) GetBucket() string {
	if config == nil {
		return ""
	}
	return config.Bucket
}

func (config *Config) NewClient() *Client {
	return &Client{
		Client: s3.NewFromConfig(config.AwsConfig(), config.ApplyToOptions),
		config: config,
	}
}