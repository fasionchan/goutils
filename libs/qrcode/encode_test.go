package qrcode

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncodePNGDecodable AC-G1：EncodePNG 产出可被 image/png.Decode 成功解码的字节流。
func TestEncodePNGDecodable(t *testing.T) {
	data, err := EncodePNG("png payload")
	require.NoError(t, err)
	_, err = png.Decode(bytes.NewReader(data))
	assert.NoError(t, err)
}

// TestEncodeJPEGDecodable AC-G2：EncodeJPEG 默认 quality=90，产出可被
// image/jpeg.Decode 成功解码的字节流；未显式传 quality 时行为等同 90。
func TestEncodeJPEGDecodable(t *testing.T) {
	data, err := EncodeJPEG("jpeg payload")
	require.NoError(t, err)
	img, err := jpeg.Decode(bytes.NewReader(data))
	assert.NoError(t, err)
	assert.Equal(t, DefaultSize, img.Bounds().Dx())
	assert.Equal(t, DefaultSize, img.Bounds().Dy())
}

// TestEncodeJPEGDefaultQualityEquals90 AC-G2：默认质量与显式 90 一致。
func TestEncodeJPEGDefaultQualityEquals90(t *testing.T) {
	content := "quality check"
	defaultData, err := EncodeJPEG(content)
	require.NoError(t, err)
	explicitData, err := EncodeJPEG(content, EncodeOptions{Quality: 90})
	require.NoError(t, err)
	assert.Equal(t, defaultData, explicitData, "default quality must equal explicit 90")
}

// TestEncodeRoundTrip AC-G3：png 与 jpeg(q=90) round-trip，Decode 均与原文一致。
func TestEncodeRoundTrip(t *testing.T) {
	content := "round trip content 0123456789"
	pngData, err := EncodePNG(content)
	require.NoError(t, err)
	got, err := DecodeBytes(pngData)
	require.NoError(t, err)
	assert.Equal(t, content, got)

	jpegData, err := EncodeJPEG(content)
	require.NoError(t, err)
	got, err = DecodeBytes(jpegData)
	require.NoError(t, err)
	assert.Equal(t, content, got)
}

// TestEncodeLevelsDistinguishable AC-G4：可设纠错 L/M/Q/H；不同级别输出可区分。
func TestEncodeLevelsDistinguishable(t *testing.T) {
	content := "level differentiation"
	outputs := make(map[string][]byte)
	for _, level := range []RecoveryLevel{LevelL, LevelM, LevelQ, LevelH} {
		data, err := EncodePNG(content, EncodeOptions{Level: level, Size: 128})
		require.NoError(t, err)
		outputs[level.String()] = data
	}
	for name, data := range outputs {
		for otherName, other := range outputs {
			if name != otherName {
				assert.NotEqual(t, data, other,
					"levels %s and %s must produce distinguishable output", name, otherName)
			}
		}
	}
}

// TestEncodeSizeExact AC-G5：可设像素边长；输出宽高等于请求值。
func TestEncodeSizeExact(t *testing.T) {
	for _, size := range []int{128, 256} {
		img, err := EncodeImage("size check", EncodeOptions{Size: size})
		require.NoError(t, err)
		assert.Equal(t, size, img.Bounds().Dx())
		assert.Equal(t, size, img.Bounds().Dy())

		data, err := EncodePNG("size check", EncodeOptions{Size: size})
		require.NoError(t, err)
		decoded, err := png.Decode(bytes.NewReader(data))
		require.NoError(t, err)
		assert.Equal(t, size, decoded.Bounds().Dx())
		assert.Equal(t, size, decoded.Bounds().Dy())
	}
}

// TestEncodeEmptyContent AC-G6：空字符串 → ErrInvalidContent，不写出半截文件。
func TestEncodeEmptyContent(t *testing.T) {
	_, err := EncodePNG("")
	assert.ErrorIs(t, err, ErrInvalidContent)
	_, err = EncodeJPEG("")
	assert.ErrorIs(t, err, ErrInvalidContent)
	_, err = EncodeImage("")
	assert.ErrorIs(t, err, ErrInvalidContent)
}

// TestEncodeTooLongContent AC-G6：超容量内容 → ErrInvalidContent。
func TestEncodeTooLongContent(t *testing.T) {
	long := strings.Repeat("a", 5000) // 超过 QR 最大容量（约 2953 bytes @ 最高版本）
	_, err := EncodePNG(long, EncodeOptions{Level: LevelL})
	assert.ErrorIs(t, err, ErrInvalidContent)
}

// TestEncodeInvalidOption AC 补充：非法选项 → ErrInvalidOption。
func TestEncodeInvalidOption(t *testing.T) {
	_, err := EncodePNG("x", EncodeOptions{Level: RecoveryLevel(99)})
	assert.ErrorIs(t, err, ErrInvalidOption)

	// Size<=0 按文档约定回退到默认值，不报错。
	_, err = EncodePNG("x", EncodeOptions{Size: -1})
	assert.NoError(t, err)

	_, err = EncodeJPEG("x", EncodeOptions{Quality: 101})
	assert.ErrorIs(t, err, ErrInvalidOption)
}

// TestParseRecoveryLevel AC 补充：级别字符串解析。
func TestParseRecoveryLevel(t *testing.T) {
	for _, tc := range []struct {
		in      string
		want    RecoveryLevel
		wantErr bool
	}{
		{in: "L", want: LevelL},
		{in: "m", want: LevelM},
		{in: "Q", want: LevelQ},
		{in: "H", want: LevelH},
		{in: " X ", wantErr: true},
		{in: "", wantErr: true},
	} {
		got, err := ParseRecoveryLevel(tc.in)
		if tc.wantErr {
			assert.ErrorIs(t, err, ErrInvalidOption, "input %q", tc.in)
			continue
		}
		require.NoError(t, err, "input %q", tc.in)
		assert.Equal(t, tc.want, got)
	}
}

// TestEncodeImageType AC 补充：EncodeImage 返回合法的 image.Image。
func TestEncodeImageType(t *testing.T) {
	img, err := EncodeImage("image type")
	require.NoError(t, err)
	var _ image.Image = img
}
