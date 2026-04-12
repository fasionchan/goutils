package s3x

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/fasionchan/goutils/stl"
)

// S3 分片约束：除最后一片外每片至少 5MiB；最多 10000 片。
const (
	DefaultMultipartPartSize = 8 << 20 // 8 MiB
	MinMultipartPartSize     = 5 << 20 // 5 MiB
	MaxMultipartParts        = 10000   // 10000 片
)

type (
	CreateMultipartUploadOptions = func(*s3.CreateMultipartUploadInput)
	UploadPartOptions            = func(*s3.UploadPartInput)
)

func (client *Client) NewMultipartUploader(ctx context.Context, input *s3.CreateMultipartUploadInput, opts ...CreateMultipartUploadOptions) (*MultipartUploader, error) {
	return NewMultipartUploader(ctx, client.Client, input, opts...)
}

type MultipartUploader struct {
	*s3.Client

	ctx            context.Context
	createOutput   *s3.CreateMultipartUploadOutput
	completedParts []types.CompletedPart
	completed      bool
}

func NewMultipartUploader(ctx context.Context, client *s3.Client, input *s3.CreateMultipartUploadInput, opts ...func(*s3.CreateMultipartUploadInput)) (*MultipartUploader, error) {
	input = stl.NewOptions(opts...).Apply(input)

	created, err := client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("s3x.NewMultipartUploader.CreateMultipartUpload: %w", err)
	}

	return &MultipartUploader{
		Client: client,

		ctx:          ctx,
		createOutput: created,
	}, nil
}

func (uploader *MultipartUploader) UploadPart(opts ...UploadPartOptions) error {
	input := stl.NewOptions(opts...).Apply(&s3.UploadPartInput{
		Bucket:     uploader.createOutput.Bucket,
		Key:        uploader.createOutput.Key,
		UploadId:   uploader.createOutput.UploadId,
		PartNumber: aws.Int32(int32(len(uploader.completedParts) + 1)),
	})

	out, err := uploader.Client.UploadPart(uploader.ctx, input)
	if err != nil {
		return fmt.Errorf("s3x.MultipartUploader.UploadPart: %w", err)
	}

	uploader.completedParts = append(uploader.completedParts, types.CompletedPart{
		ETag:       out.ETag,
		PartNumber: input.PartNumber,
	})

	return nil
}

func (uploader *MultipartUploader) UploadResourceSegments(reader SegmentableResourceReader) error {
	for {
		segment, err := reader.Next(uploader.ctx)
		if err != nil {
			return err
		}

		if err = uploader.UploadPart(func(input *s3.UploadPartInput) {
			input.Body = segment
		}); err != nil {
			return err
		}

		if !segment.HasNext() {
			break
		}
	}

	return nil
}

func (uploader *MultipartUploader) Complete(ctx context.Context) (*s3.CompleteMultipartUploadOutput, error) {
	out, err := uploader.Client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   uploader.createOutput.Bucket,
		Key:      uploader.createOutput.Key,
		UploadId: uploader.createOutput.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: uploader.completedParts,
		},
	})
	uploader.completed = err == nil
	return out, err
}

func (uploader *MultipartUploader) Close() error {
	if uploader.completed {
		return nil
	}

	return uploader.abort(uploader.ctx)
}

func (uploader *MultipartUploader) abort(ctx context.Context) error {
	_, err := uploader.Client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   uploader.createOutput.Bucket,
		Key:      uploader.createOutput.Key,
		UploadId: uploader.createOutput.UploadId,
	})
	return err
}
