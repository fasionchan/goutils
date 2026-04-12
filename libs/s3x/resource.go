package s3x

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fasionchan/goutils/std/mimex"
	"github.com/fasionchan/goutils/stl"
)

type Resource interface {
	Read(p []byte) (n int, err error)
	GetContentType() string
	GetContentDisposition() string
}

type ResourceSegment interface {
	Resource
	HasNext() bool
	Len() int64
	TotalLen() int64
}

type SegmentableResourceReader interface {
	Next(ctx context.Context) (ResourceSegment, error)
}

// func (client *Client) StoreResource(ctx context.Context, bucket, key string, offset int64, length int64) (Resource, error) {
// 	return client.Read(ctx, key, offset, length)
// }

func (client *Client) StoreSegmentableResource(ctx context.Context, reader SegmentableResourceReader, input *s3.CreateMultipartUploadInput, opts ...func(*s3.CreateMultipartUploadInput)) (any, error) {
	input = stl.NewOptions(opts...).Apply(input)

	// 读取第一个分片
	segment, err := reader.Next(ctx)
	if err != nil {
		return nil, err
	}

	contentType := segment.GetContentType()
	contentDisposition := segment.GetContentDisposition()

	// 预读探测 ContentType
	peekContentType, contentReader, _ := mimex.PeekContentTypeSmart(contentType, contentDisposition, segment, 0)
	if contentType == "" {
		contentType = peekContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	// 如果只有一个分片，则直接上传
	if !segment.HasNext() {
		putObjectInput := stl.NewOption(func(o *s3.PutObjectInput) {
			o.Body = contentReader
			o.ContentType = aws.String(contentType)
			o.ContentDisposition = aws.String(contentDisposition)
		}).Apply(PutObjectInputFromCreateMultipartUploadInput(input))

		_, err = client.PutObject(ctx, putObjectInput)
		return nil, err
	}

	// 创建分片上传器
	uploader, err := client.NewMultipartUploader(ctx, input, func(input *s3.CreateMultipartUploadInput) {
		input.ContentType = aws.String(contentType)
		input.ContentDisposition = aws.String(segment.GetContentDisposition())
	})
	if err != nil {
		return nil, err
	}
	defer uploader.Close()

	// 上传第一个分片
	if err = uploader.UploadPart(func(input *s3.UploadPartInput) {
		input.Body = contentReader
	}); err != nil {
		return nil, err
	}

	// 上传剩余分片
	if err = uploader.UploadResourceSegments(reader); err != nil {
		return nil, err
	}

	// 完成分片上传
	_, err = uploader.Complete(ctx)
	return nil, err
}
