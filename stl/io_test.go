package stl

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	assert.NoError(t, Close(nil))
	assert.NoError(t, Close(io.Reader(nil)))
	assert.NoError(t, Close(io.Closer(nil)))
	assert.NoError(t, Close((*struct{})(nil)))
	assert.NoError(t, Close(io.NopCloser(nil)))
	assert.NoError(t, Close(io.NopCloser(bytes.NewReader([]byte{}))))
}
