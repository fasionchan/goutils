/*
 * Author: fasion
 * Created time: 2024-03-19 13:56:39
 * Last Modified by: fasion
 * Last Modified time: 2025-07-22 18:44:18
 */

package httpx

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

type MultipartFormItem interface {
	WriteTo(*multipart.Writer) error
}

type MultipartFormField struct {
	name  string
	value string
}

func NewMultipartFormField(name, value string) *MultipartFormField {
	return &MultipartFormField{name: name, value: value}
}

func (field *MultipartFormField) WriteTo(writer *multipart.Writer) error {
	return writer.WriteField(field.name, field.value)
}

type MultipartFormFile struct {
	fieldName   string
	fileName    string
	fileContent io.Reader
}

func NewMultipartFormFile(fieldName, fileName string, fileContent io.Reader) *MultipartFormFile {
	return &MultipartFormFile{
		fieldName:   fieldName,
		fileName:    fileName,
		fileContent: fileContent,
	}
}

func (file *MultipartFormFile) WriteTo(writer *multipart.Writer) error {
	part, err := writer.CreateFormFile(file.fieldName, file.fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file.fileContent)
	if err != nil {
		return err
	}

	return nil
}

type MultipartFormItems []MultipartFormItem

func (items MultipartFormItems) Append(others ...MultipartFormItem) MultipartFormItems {
	return append(items, others...)
}

func (items MultipartFormItems) AppendField(name, value string) MultipartFormItems {
	return items.Append(NewMultipartFormField(name, value))
}

func (items MultipartFormItems) AppendFile(fieldName string, file *os.File) MultipartFormItems {
	return items.AppendFileByReader(fieldName, file.Name(), file)
}

func (items MultipartFormItems) AppendFileByReader(fieldName, fileName string, file io.Reader) MultipartFormItems {
	return items.Append(NewMultipartFormFile(fieldName, fileName, file))
}

func (items MultipartFormItems) WriteTo(writer *multipart.Writer) error {
	for _, item := range items {
		err := item.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (items MultipartFormItems) WriteToBuffer(buffer *bytes.Buffer) (*multipart.Writer, error) {
	writer := multipart.NewWriter(buffer)
	defer writer.Close()

	return writer, items.WriteTo(writer)
}

func (items MultipartFormItems) Marshal() (buffer *bytes.Buffer, writer *multipart.Writer, err error) {
	buffer = bytes.NewBuffer(nil)
	writer, err = items.WriteToBuffer(buffer)
	return
}

func (items MultipartFormItems) MarshalToReader() (io.Reader, string, error) {
	buffer, writer, err := items.Marshal()
	if err != nil {
		return nil, "", err
	}

	return buffer, writer.FormDataContentType(), nil
}

func (items MultipartFormItems) MarshalToStreamingReadCloser() (io.ReadCloser, string, error) {
	pr, pw := io.Pipe()

	writer := multipart.NewWriter(pw)
	contentType := writer.FormDataContentType()

	go func() {
		var err error
		defer func() {
			pw.CloseWithError(err)
		}()

		defer writer.Close()

		err = items.WriteTo(writer)
	}()

	return pr, contentType, nil
}

type MultipartForm = MultipartFormItems

func NewMultipartForm() MultipartForm {
	return MultipartForm{}
}
