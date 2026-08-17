package ilink

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBotQrcode(t *testing.T) {
	client, err := NewClient()
	assert.NoError(t, err)

	result, err := client.GetBotQrcode(context.Background(), 3)
	assert.NoError(t, err)

	fmt.Println(result.Qrcode)
	fmt.Println(result.QrcodeImgContent)
}

func TestGetQrcodeStatus(t *testing.T) {
	client, err := NewClient()
	assert.NoError(t, err)

	result, err := client.GetQrcodeStatus(context.Background(), "e71054afb81e5844818a9396e75e9fce")
	assert.NoError(t, err)

	fmt.Println(result.Status)
}

func TestWaitQrcodeStatus(t *testing.T) {
	client, err := NewClient()
	assert.NoError(t, err)

	qrcodeResult, err := client.GetBotQrcode(context.Background(), 3)
	assert.NoError(t, err)

	fmt.Println(qrcodeResult.QrcodeImgContent)

	result, err := client.WaitQrcodeStatus(context.Background(), qrcodeResult.Qrcode)
	assert.NoError(t, err)

	fmt.Println(result)
	fmt.Println(result.Status)
}
