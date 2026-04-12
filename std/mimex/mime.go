package mimex

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

// 智能探测 ContentType，优先根据 ContentType 和 ContentDisposition 探测，内容探测为兜底
func PeekContentTypeSmart(contentType string, contentDisposition string, contentReader io.Reader, peekSize int) (string, io.Reader, error) {
	if contentType != "" {
		mimeType, err := ParseMimeType(contentType)
		if err == nil && mimeType.GetType() != "application/octet-stream" {
			return contentType, contentReader, nil
		}
	}

	if contentDisposition != "" {
		mimeType, err := ParseDispositionMimeType(contentDisposition)
		if err == nil && mimeType.GetType() != "application/octet-stream" {
			return mimeType.String(), contentReader, nil
		}
	}

	mimeType, reader, err := PeekContentType(contentType, contentDisposition, contentReader, peekSize)
	if err != nil {
		return "", reader, err
	}

	return mimeType.String(), reader, nil
}

// 预读探测 ContentType
func PeekContentType(contentType string, contentDisposition string, contentReader io.Reader, peekSize int) (mimeType *MimeType, reader io.Reader, err error) {
	if peekSize <= 0 {
		peekSize = 8 << 10
	}

	reader = contentReader

	var buffer bytes.Buffer
	n, err := io.CopyN(&buffer, contentReader, int64(peekSize))
	if n > 0 {
		reader = io.MultiReader(&buffer, contentReader)
	} else {
		err = errors.New("no content to peek type")
		return
	}

	ftype, err := filetype.Match(buffer.Bytes())
	if err != nil {
		return
	}
	if ftype == filetype.Unknown {
		err = errors.New("unknown file type")
		return
	}

	mimeType = &MimeType{
		Type:    ftype.MIME.Value,
		TopType: ftype.MIME.Type,
		SubType: ftype.MIME.Subtype,
	}

	return
}

func ParseDispositionMimeType(contentDisposition string) (*MimeType, error) {
	if contentDisposition == "" {
		return nil, errors.New("content disposition is empty")
	}

	mimeType, err := ParseMimeType(contentDisposition)
	if err != nil {
		return nil, err
	}

	filename := mimeType.GetParams()["filename"]
	if filename == "" {
		return nil, errors.New("content disposition filename is empty")
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		return nil, errors.New("content disposition extension is empty")
	}

	return ParseMimeType(mime.TypeByExtension(ext))
}

type MimeType struct {
	Type    string
	TopType string
	SubType string
	Params  map[string]string
}

func ParseMimeType(mimeType string) (*MimeType, error) {
	t, params, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(t, "/")
	topType := parts[0]
	subType := ""
	if len(parts) > 1 {
		subType = parts[1]
	}

	return &MimeType{
		Type:    t,
		TopType: topType,
		SubType: subType,
		Params:  params,
	}, nil
}

func (mt *MimeType) GetType() string {
	return mt.Type
}

func (mt *MimeType) GetTopType() string {
	return mt.TopType
}

func (mt *MimeType) GetSubType() string {
	return mt.SubType
}

func (mt *MimeType) GetParams() map[string]string {
	return mt.Params
}

func (mt *MimeType) String() string {
	return mime.FormatMediaType(mt.Type, mt.Params)
}
