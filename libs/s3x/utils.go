package s3x

import "github.com/aws/aws-sdk-go-v2/service/s3"

func PutObjectInputFromCreateMultipartUploadInput(input *s3.CreateMultipartUploadInput) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Bucket:               input.Bucket,
		Key:                  input.Key,
		ContentType:          input.ContentType,
		ContentDisposition:   input.ContentDisposition,
		ACL:                  input.ACL,
		Metadata:             input.Metadata,
		Expires:              input.Expires,
		CacheControl:         input.CacheControl,
		BucketKeyEnabled:     input.BucketKeyEnabled,
		StorageClass:         input.StorageClass,
		ServerSideEncryption: input.ServerSideEncryption,
		ChecksumAlgorithm:    input.ChecksumAlgorithm,
		ContentEncoding:      input.ContentEncoding,
		ContentLanguage:      input.ContentLanguage,
	}
}
