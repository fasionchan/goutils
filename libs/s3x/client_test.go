package s3x

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func newTestClient() (*Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return NewClient(aws.Config{
		Region:      os.Getenv("REGION"),
		Credentials: NewStaticCredentialsProvider(os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY")),
	}, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(os.Getenv("ENDPOINT"))
	}), nil
}
