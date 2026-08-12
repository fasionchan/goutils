package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	sqrcode "github.com/skip2/go-qrcode"
)

// EncodeOptions 是生成图片的配置项；零值字段回退到默认值
// （Level=LevelM、Size=DefaultSize、Quality=DefaultQuality）。
type EncodeOptions struct {
	// Level 纠错级别，默认 LevelM。
	Level RecoveryLevel
	// Size 输出图片的像素边长（宽高相等），默认 DefaultSize。
	Size int
	// Quality jpeg 输出质量（1-100），仅 EncodeJPEG 生效，默认 DefaultQuality。
	Quality int
}

// normalized 把零值字段替换为默认值。
func (o EncodeOptions) normalized() EncodeOptions {
	if o.Level == 0 {
		o.Level = LevelM
	}
	if o.Size <= 0 {
		o.Size = DefaultSize
	}
	if o.Quality <= 0 {
		o.Quality = DefaultQuality
	}
	return o
}

// validate 校验归一化后的选项取值。
func (o EncodeOptions) validate() error {
	if !o.Level.Valid() {
		return fmt.Errorf("%w: invalid recovery level %d", ErrInvalidOption, o.Level)
	}
	if o.Size < 1 || o.Size > 4096 {
		return fmt.Errorf("%w: size must be in [1, 4096], got %d", ErrInvalidOption, o.Size)
	}
	if o.Quality < 1 || o.Quality > 100 {
		return fmt.Errorf("%w: quality must be in [1, 100], got %d", ErrInvalidOption, o.Quality)
	}
	return nil
}

// normalizeEncodeOptions 处理可变参数：不传时使用全部默认值。
func normalizeEncodeOptions(opts []EncodeOptions) EncodeOptions {
	if len(opts) == 0 {
		return (EncodeOptions{}).normalized()
	}
	return opts[0].normalized()
}

// EncodeImage 生成包含 content 的二维码 image.Image，宽高等于 Size 像素
// （当 Size 小于码图实际所需大小时，由底层库自动放大到码图尺寸）。
func EncodeImage(content string, opts ...EncodeOptions) (image.Image, error) {
	o := normalizeEncodeOptions(opts)
	if err := o.validate(); err != nil {
		return nil, err
	}
	if content == "" {
		return nil, ErrInvalidContent
	}
	code, err := sqrcode.New(content, o.Level.skip2())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidContent, err)
	}
	return code.Image(o.Size), nil
}

// EncodePNG 生成包含 content 的二维码 png 字节流。
func EncodePNG(content string, opts ...EncodeOptions) ([]byte, error) {
	img, err := EncodeImage(content, opts...)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOption, err)
	}
	return buf.Bytes(), nil
}

// EncodeJPEG 生成包含 content 的二维码 jpeg 字节流，质量默认 DefaultQuality。
func EncodeJPEG(content string, opts ...EncodeOptions) ([]byte, error) {
	o := normalizeEncodeOptions(opts)
	img, err := EncodeImage(content, opts...)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: o.Quality}); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOption, err)
	}
	return buf.Bytes(), nil
}
