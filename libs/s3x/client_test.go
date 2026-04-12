package s3x

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func newTestClient() (*Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return NewClient(NewConfigFromDefaultEnv(), func(options *s3.Options) {
		options.BaseEndpoint = aws.String(os.Getenv("AWS_ENDPOINT"))
	}), nil
}

func TestGetBucketAcl(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	acl, err := client.GetBucketAcl(context.Background(), &s3.GetBucketAclInput{
		Bucket: aws.String("fasionchan-cdn"),
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(acl)

	for _, grant := range acl.Grants {
		fmt.Println(grant)
		fmt.Println(grant.Grantee.DisplayName, grant.Permission)
	}

	fmt.Println(acl.Owner)

	fmt.Println(acl.ResultMetadata)
}

func TestGetBucketPolicy(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	policy, err := client.GetBucketPolicy(context.Background(), &s3.GetBucketPolicyInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(policy)

	fmt.Println(policy.Policy)
}
