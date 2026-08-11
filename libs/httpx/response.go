/*
 * Author: wxhuangjuguan
 * Created time: 2025-10-29 14:30:00
 * Last Modified by: wxhuangjuguan
 * Last Modified time: 2025-10-29 14:30:00
 */

package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/fasionchan/goutils/stl"
)

type HttpResponseDigest struct {
	Status     string
	StatusCode int
	Proto      string
	Header     http.Header
	Body       any
}

func newHttpResponseDigestForRaw(response *http.Response) *HttpResponseDigest {
	return &HttpResponseDigest{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Proto:      response.Proto,
		Header:     response.Header,
	}

}

func (digest *HttpResponseDigest) WithBody(body any) *HttpResponseDigest {
	if digest == nil {
		return nil
	}

	digest.Body = body
	return digest
}

type HttpResponsePipe struct {
	io.WriteCloser
	header http.Header
	reader io.ReadCloser

	response *http.Response
}

func NewHttpResponsePipeWithBuffer(buffer *bytes.Buffer) *HttpResponsePipe {
	if buffer == nil {
		buffer = bytes.NewBuffer(nil)
	}

	return &HttpResponsePipe{
		WriteCloser: stl.NewNopCloseWriter(buffer),
		header:      make(http.Header),
		reader:      io.NopCloser(buffer),
	}
}

func NewHttpResponsePipe() *HttpResponsePipe {
	reader, writer := io.Pipe()
	return &HttpResponsePipe{
		WriteCloser: writer,
		header:      make(http.Header),
		reader:      reader,
	}
}

func (pipe *HttpResponsePipe) Header() http.Header {
	if pipe == nil {
		return nil
	}
	return pipe.header
}

func (pipe *HttpResponsePipe) WriteHeader(statusCode int) {
	if pipe == nil {
		return
	}

	pipe.response = &http.Response{
		StatusCode: statusCode,
		Header:     pipe.header,
		Body:       pipe.reader,
	}
}

func (pipe *HttpResponsePipe) CloseReader() error {
	reader := pipe.GetReadCloser()
	if reader == nil {
		return nil
	}
	return reader.Close()
}

func (pipe *HttpResponsePipe) CloseWriter() error {
	writer := pipe.GetWriteCloser()
	if writer == nil {
		return nil
	}
	return writer.Close()
}

func (pipe *HttpResponsePipe) ClosePipe() error {
	if pipe == nil {
		return nil
	}

	return stl.NewErrors(
		pipe.CloseWriter(),
		pipe.CloseReader(),
	).Simplify()
}

func (pipe *HttpResponsePipe) GetResponse() *http.Response {
	if pipe == nil {
		return nil
	}

	return pipe.response
}

func (pipe *HttpResponsePipe) GetReadCloser() io.ReadCloser {
	if pipe == nil {
		return nil
	}
	return pipe.reader
}

func (pipe *HttpResponsePipe) GetWriteCloser() io.WriteCloser {
	if pipe == nil {
		return nil
	}
	return pipe.WriteCloser
}

func CopyResponseBody(dst http.ResponseWriter, src io.Reader, bufferSize int, flush bool) (int, error) {
	buffer := make([]byte, bufferSize)

	var bytesCopied int
	for {
		bytesRead, err := src.Read(buffer)
		if bytesRead > 0 {
			bytesWritten, err := WriteResponseBody(dst, buffer[:bytesRead], flush)
			bytesCopied += bytesWritten
			if err != nil {
				return bytesCopied, fmt.Errorf("CopyHttpResponseBody.WriteHttpResponseBody: %w", err)
			}
		}

		if err != nil {
			if err == io.EOF {
				return bytesCopied, nil
			}

			return bytesCopied, fmt.Errorf("CopyHttpResponseBody.ReadSrc: %w", err)
		}

		if bytesRead == 0 {
			return bytesCopied, nil
		}
	}
}

func WriteResponseBody(writer http.ResponseWriter, data []byte, flush bool) (bytesWritten int, err error) {
	flusher, isFlusher := writer.(http.Flusher)

	var n int
	for len(data) > 0 {
		n, err = writer.Write(data)
		bytesWritten += n
		data = data[n:]

		if err != nil {
			return
		}

		if flush && isFlusher {
			flusher.Flush()
		}
	}

	return
}