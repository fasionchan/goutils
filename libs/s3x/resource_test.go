package s3x

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestStoreSegmentableResource(t *testing.T) {
	client, err := newTestClient()
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, path := range []string{
		"../../std/mimex/testdata/file.txt",
		"../../std/mimex/testdata/file.png",
		"../../std/mimex/testdata/file.pdf",
		"../../std/mimex/testdata/file.docx",
		"../../std/mimex/testdata/file.xlsx",
		"../../std/mimex/testdata/file.pptx",
		"../../std/mimex/testdata/file.doc",
		"../../std/mimex/testdata/file.xls",
		"../../std/mimex/testdata/file.ppt",
	} {
		RunStoreSegmentableResourceTestCase(t, client, path)
	}
}

func RunStoreSegmentableResourceTestCase(t *testing.T, client *Client, path string) {
	fmt.Println("-------- test case", path, "--------")

	reader, err := OpenFileForSegmentableResource(path, MinMultipartPartSize)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	fileName := filepath.Base(path)

	bucket := os.Getenv("AWS_S3_BUCKET")
	key := "test/" + fileName

	resource, err := client.StoreSegmentableResource(context.Background(), reader, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = resource

	object, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("file path:", path)
	fmt.Println("content type:", aws.ToString(object.ContentType))
	fmt.Println("content disposition:", aws.ToString(object.ContentDisposition))
	fmt.Println("content length:", aws.ToInt64(object.ContentLength))
	fmt.Println()
}
