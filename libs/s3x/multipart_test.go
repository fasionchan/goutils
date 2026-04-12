package s3x

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestMultipartUpload(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	uploader, err := client.NewMultipartUploader(context.Background(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String("fasionchan-cdn"),
		Key:    aws.String("test/multipart.txt"),
	})
	if err != nil {
		t.Fatalf("failed to create multipart uploader: %v", err)
	}
	defer uploader.Close()

	if err = uploader.UploadPart(func(input *s3.UploadPartInput) {
		input.Body = bytes.NewReader(bytes.Repeat([]byte("a"), DefaultMultipartPartSize))
	}); err != nil {
		t.Fatalf("failed to upload part: %v", err)
	}

	if err = uploader.UploadPart(func(input *s3.UploadPartInput) {
		input.Body = bytes.NewReader(([]byte("b")))
	}); err != nil {
		t.Fatalf("failed to upload part: %v", err)
	}

	complete, err := uploader.Complete(context.Background())
	if err != nil {
		t.Fatalf("failed to complete multipart upload: %v", err)
	}

	fmt.Println(complete)
}
