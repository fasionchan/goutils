package s3x

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fasionchan/goutils/std/mimex"
	"github.com/fasionchan/goutils/stl"
)

type ResourceInfo interface {
	GetContentType() string
	GetContentDisposition() string
	GetTotalLen() int64
}

type StaticResourceInfo struct {
	contentType        string
	contentDisposition string
	totalLen           int64
}

func (info *StaticResourceInfo) GetContentType() string {
	return info.contentType
}

func (info *StaticResourceInfo) GetContentDisposition() string {
	return info.contentDisposition
}

func (info *StaticResourceInfo) GetTotalLen() int64 {
	return info.totalLen
}

type Resource interface {
	io.ReadCloser
	ResourceInfo
}

type ResourceSegmentInfo interface {
	GetOffset() int64
	GetLen() int64
	HasNext() bool
}

type StaticResourceSegmentInfo struct {
	offset  int64
	len     int64
	hasNext bool
}

func (info *StaticResourceSegmentInfo) GetOffset() int64 {
	return info.offset
}

func (info *StaticResourceSegmentInfo) GetLen() int64 {
	return info.len
}

func (info *StaticResourceSegmentInfo) HasNext() bool {
	return info.hasNext
}

type ResourceSegment interface {
	Resource
	ResourceSegmentInfo
}

type SegmentableResourceReader interface {
	Next(ctx context.Context) (ResourceSegment, error)
	Close() error
}

type CommonResourceSegmentReader struct {
	io.ReadCloser
	ResourceInfo
	ResourceSegmentInfo
}

type SegmentableFileReader struct {
	blockSize int64
	*os.File
	totalLen           int64
	contentType        string
	contentDisposition string
	nextOffset         int64
}

func NewSegmentableFileReader(file *os.File, blockSize int64) (*SegmentableFileReader, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	contentDisposition := mime.FormatMediaType("attachment", map[string]string{"filename": file.Name()})

	return &SegmentableFileReader{
		blockSize:          blockSize,
		File:               file,
		totalLen:           info.Size(),
		contentType:        "application/octet-stream",
		contentDisposition: contentDisposition,
	}, nil
}

func OpenFileForSegmentableResource(path string, blockSize int64) (SegmentableResourceReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader, err := NewSegmentableFileReader(file, blockSize)
	if err != nil {
		file.Close()
		return nil, err
	}

	return reader, nil
}

func (reader *SegmentableFileReader) GetContentType() string {
	return reader.contentType
}

func (reader *SegmentableFileReader) GetContentDisposition() string {
	return reader.contentDisposition
}

func (reader *SegmentableFileReader) HasNext() bool {
	return reader.nextOffset < reader.totalLen
}

func (reader *SegmentableFileReader) GetTotalLen() int64 {
	return reader.totalLen
}

func (reader *SegmentableFileReader) Next(ctx context.Context) (ResourceSegment, error) {
	if !reader.HasNext() {
		return nil, io.EOF
	}

	offset := reader.nextOffset

	reader.nextOffset += reader.blockSize
	if reader.nextOffset > reader.totalLen {
		reader.nextOffset = reader.totalLen
	}

	len := reader.nextOffset - offset
	hasNext := reader.HasNext()

	dupFile, err := dupAndSeek(reader.File, offset)
	if err != nil {
		return nil, err
	}

	return &CommonResourceSegmentReader{
		ReadCloser:   dupFile,
		ResourceInfo: reader,
		ResourceSegmentInfo: &StaticResourceSegmentInfo{
			offset:  offset,
			len:     len,
			hasNext: hasNext,
		},
	}, nil
}

func dupFile(file *os.File) (*os.File, error) {
	// duplicate the file descriptor
	fd, err := syscall.Dup(int(file.Fd()))
	if err != nil {
		return nil, err
	}

	return os.NewFile(uintptr(fd), file.Name()), nil
}

func dupAndSeek(file *os.File, offset int64) (*os.File, error) {
	dupFile, err := dupFile(file)
	if err != nil {
		return nil, err
	}

	fileOffset, err := dupFile.Seek(offset, io.SeekStart)
	if err != nil {
		dupFile.Close()
		return nil, err
	}

	if fileOffset != offset {
		dupFile.Close()
		return nil, fmt.Errorf("failed to seek to offset %d: %w", offset, err)
	}

	return dupFile, nil
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
	defer segment.Close()

	contentType := segment.GetContentType()
	contentDisposition := segment.GetContentDisposition()

	// 预读探测 ContentType
	peekContentType, contentReader, _ := mimex.PeekContentTypeSmart(contentType, contentDisposition, segment, 0)
	if contentType == "" || contentType == "application/octet-stream" {
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
