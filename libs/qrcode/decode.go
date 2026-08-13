package qrcode

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	zxmulti "github.com/makiuchi-d/gozxing/multi/qrcode"
	zxqrcode "github.com/makiuchi-d/gozxing/qrcode"
	// 注册 png / jpeg 解码器，供 image.Decode 识别输入格式。
	_ "image/jpeg"
	_ "image/png"
)

// supportedFormats 是本库识别入口接受的图片格式。
var supportedFormats = map[string]bool{
	"png":  true,
	"jpeg": true,
}

// DecodeFile 从图片文件解码二维码内容，返回第一个识别到的 payload。
// 文件不存在等打开错误原样透传（可用 errors.Is(err, os.ErrNotExist) 断言）。
func DecodeFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return DecodeReader(f)
}

// DecodeBytes 从 png/jpeg 字节流解码二维码内容，返回第一个识别到的 payload。
func DecodeBytes(data []byte) (string, error) {
	return DecodeReader(bytes.NewReader(data))
}

// DecodeReader 从 io.Reader 读取图片并解码二维码内容，返回第一个识别到的 payload。
func DecodeReader(r io.Reader) (string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrUnsupportedFormat, err)
	}
	if !supportedFormats[strings.ToLower(format)] {
		return "", fmt.Errorf("%w: got %q", ErrUnsupportedFormat, format)
	}
	return DecodeImage(img)
}

// DecodeImage 从 image.Image 解码二维码内容，返回第一个识别到的 payload。
// 该入口跳过格式探测，直接送解码器。
func DecodeImage(img image.Image) (string, error) {
	text, err := decodeImage(img)
	if err != nil {
		return "", err
	}
	return text, nil
}

// DecodeAllFile 从图片文件解码全部二维码内容，按解码器探测顺序返回 payload 列表。
// 图片中没有任何二维码时返回 ErrNotFound。
func DecodeAllFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodeAllReader(f)
}

// DecodeAllBytes 从 png/jpeg 字节流解码全部二维码内容。
func DecodeAllBytes(data []byte) ([]string, error) {
	return DecodeAllReader(bytes.NewReader(data))
}

// DecodeAllReader 从 io.Reader 读取图片并解码全部二维码内容。
func DecodeAllReader(r io.Reader) ([]string, error) {
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedFormat, err)
	}
	if !supportedFormats[strings.ToLower(format)] {
		return nil, fmt.Errorf("%w: got %q", ErrUnsupportedFormat, format)
	}
	return DecodeAllImage(img)
}

// DecodeAllImage 从 image.Image 解码全部二维码内容。
func DecodeAllImage(img image.Image) ([]string, error) {
	return decodeAllImage(img)
}

// decodeImage 使用 gozxing 单码 Reader 解码一张图片。
func decodeImage(img image.Image) (string, error) {
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", wrapDecodeError(err)
	}
	result, err := zxqrcode.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		return "", wrapDecodeError(err)
	}
	return result.GetText(), nil
}

// decodeAllImage 使用 gozxing MultiReader 解码一张图片中的全部二维码。
func decodeAllImage(img image.Image) ([]string, error) {
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil, wrapDecodeError(err)
	}
	results, err := zxmulti.NewQRCodeMultiReader().DecodeMultiple(bitmap, nil)
	if err != nil {
		return nil, wrapDecodeError(err)
	}
	if len(results) == 0 {
		return nil, ErrNotFound
	}
	texts := make([]string, 0, len(results))
	for _, result := range results {
		texts = append(texts, result.GetText())
	}
	return texts, nil
}

// wrapDecodeError 把 gozxing 的异常映射为本包的错误哨兵：
//   - NotFoundException → ErrNotFound（未找到二维码）
//   - ChecksumException / FormatException → ErrDecodeFailed（找到但解码失败）
//   - 其它错误 → 包装为 ErrDecodeFailed 并保留原始信息
func wrapDecodeError(err error) error {
	var notFound gozxing.NotFoundException
	if errors.As(err, &notFound) {
		return ErrNotFound
	}
	var checksum gozxing.ChecksumException
	var format gozxing.FormatException
	if errors.As(err, &checksum) || errors.As(err, &format) {
		return fmt.Errorf("%w: %v", ErrDecodeFailed, err)
	}
	return fmt.Errorf("%w: %v", ErrDecodeFailed, err)
}
