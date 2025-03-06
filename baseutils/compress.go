/*
 * Author: fasion
 * Created time: 2025-03-06 13:45:29
 * Last Modified by: fasion
 * Last Modified time: 2025-03-06 13:47:41
 */

package baseutils

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Gzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Gunzip(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	return io.ReadAll(gz)
}
