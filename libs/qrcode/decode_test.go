package qrcode

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	sqrcode "github.com/skip2/go-qrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// encodePNGForTest 生成指定内容的 png 字节流。
func encodePNGForTest(t *testing.T, content string) []byte {
	t.Helper()
	data, err := EncodePNG(content)
	require.NoError(t, err)
	return data
}

// TestDecodePNGRoundTrip AC-R1：本库生成的合法 QR png，Decode 结果与原文一致。
func TestDecodePNGRoundTrip(t *testing.T) {
	contents := []string{
		"hello world",
		"https://example.com/path?query=1&x=2",
		"你好，世界！UTF-8 内容",
		"x",
	}
	for _, content := range contents {
		data := encodePNGForTest(t, content)
		got, err := DecodeBytes(data)
		require.NoError(t, err)
		assert.Equal(t, content, got, "png round trip mismatch for %q", content)
	}
}

// TestDecodeJPEGRoundTrip AC-R2：本库以 jpeg quality=90 生成的合法 QR，Decode 结果与原文一致。
func TestDecodeJPEGRoundTrip(t *testing.T) {
	data, err := EncodeJPEG("jpeg round trip payload")
	require.NoError(t, err)
	got, err := DecodeBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "jpeg round trip payload", got)
}

// TestDecodeNoCode AC-R3：纯色（无码）图片 → ErrNotFound，且不返回伪造成功文本。
func TestDecodeNoCode(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))

	_, err := DecodeBytes(buf.Bytes())
	assert.ErrorIs(t, err, ErrNotFound)

	_, err = DecodeAllBytes(buf.Bytes())
	assert.ErrorIs(t, err, ErrNotFound)
}

// TestDecodeCorruptedNoPanic AC-R4：损坏至不可纠的二维码 → 返回错误，不 panic。
func TestDecodeCorruptedNoPanic(t *testing.T) {
	content := "corrupt me please"
	data := encodePNGForTest(t, content)

	// 用大块黑色矩形覆盖码图中部，破坏足够多的模块使其不可纠。
	img, err := png.Decode(bytes.NewReader(data))
	require.NoError(t, err)
	bounds := img.Bounds()
	corrupted := image.NewRGBA(bounds)
	draw.Draw(corrupted, bounds, img, image.Point{}, draw.Src)
	bar := image.Rect(bounds.Min.X+bounds.Dx()/4, bounds.Min.Y+bounds.Dy()/4,
		bounds.Min.X+3*bounds.Dx()/4, bounds.Min.Y+3*bounds.Dy()/4)
	draw.Draw(corrupted, bar, &image.Uniform{color.Black}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, corrupted))

	_, err = DecodeBytes(buf.Bytes())
	assert.Error(t, err, "corrupted qr must fail to decode")
}

// multiCodePNG 构造包含两个二维码的 png 图片。
func multiCodePNG(t *testing.T) []byte {
	t.Helper()
	q1, err := sqrcode.New("payload-alpha", sqrcode.Medium)
	require.NoError(t, err)
	q2, err := sqrcode.New("payload-beta", sqrcode.Medium)
	require.NoError(t, err)

	const side = 150
	const gap = 40
	const margin = 40
	canvas := image.NewRGBA(image.Rect(0, 0, side*2+gap+margin*2, side+margin*2))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(margin, margin, margin+side, margin+side), q1.Image(side), image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(margin+side+gap, margin, margin+side+gap+side, margin+side), q2.Image(side), image.Point{}, draw.Src)

	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, canvas))
	return buf.Bytes()
}

// TestDecodeAllMultiCode AC-R5：含 N 个可解码 QR 的样例图：
// Decode 返回恰好 1 条；DecodeAll 返回 N 条且 payload 集合与样例一致。
func TestDecodeAllMultiCode(t *testing.T) {
	data := multiCodePNG(t)

	one, err := DecodeBytes(data)
	require.NoError(t, err)
	assert.Contains(t, []string{"payload-alpha", "payload-beta"}, one,
		"Decode must return exactly one of the two payloads")

	all, err := DecodeAllBytes(data)
	require.NoError(t, err)
	assert.Len(t, all, 2, "DecodeAll must return all codes")
	assert.ElementsMatch(t, []string{"payload-alpha", "payload-beta"}, all)
}

// TestDecodeAllSingleCode AC-R5 补充：单码图 DecodeAll 恰好返回 1 条。
func TestDecodeAllSingleCode(t *testing.T) {
	data := encodePNGForTest(t, "single")
	all, err := DecodeAllBytes(data)
	require.NoError(t, err)
	assert.Equal(t, []string{"single"}, all)
}

// TestDecodeUnsupportedFormat AC-R6：非 png/jpeg 字节流 → ErrUnsupportedFormat。
func TestDecodeUnsupportedFormat(t *testing.T) {
	// 1) 随机字节：无法被 image.Decode 识别。
	_, err := DecodeBytes([]byte("this is definitely not an image"))
	assert.ErrorIs(t, err, ErrUnsupportedFormat)

	// 2) 合法 gif：格式可识别但不是 png/jpeg。
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	require.NoError(t, gif.Encode(&buf, img, nil))
	_, err = DecodeBytes(buf.Bytes())
	assert.ErrorIs(t, err, ErrUnsupportedFormat)
}

// TestDecodeFourEntryPointsConsistency AC-R7：路径 / []byte / io.Reader / image.Image
// 四种入口对同一张合法 png 解码结果一致。
func TestDecodeFourEntryPointsConsistency(t *testing.T) {
	content := "four entry points"
	data := encodePNGForTest(t, content)

	// image.Image 由该 png 解码得到。
	img, err := png.Decode(bytes.NewReader(data))
	require.NoError(t, err)

	dir := t.TempDir()
	path := filepath.Join(dir, "qr.png")
	require.NoError(t, os.WriteFile(path, data, 0o644))

	byPath, err := DecodeFile(path)
	require.NoError(t, err)
	byBytes, err := DecodeBytes(data)
	require.NoError(t, err)
	byReader, err := DecodeReader(bytes.NewReader(data))
	require.NoError(t, err)
	byImage, err := DecodeImage(img)
	require.NoError(t, err)

	assert.Equal(t, content, byPath)
	assert.Equal(t, content, byBytes)
	assert.Equal(t, content, byReader)
	assert.Equal(t, content, byImage)

	// DecodeAll 四种入口同样一致。
	allPath, err := DecodeAllFile(path)
	require.NoError(t, err)
	allBytes, err := DecodeAllBytes(data)
	require.NoError(t, err)
	allReader, err := DecodeAllReader(bytes.NewReader(data))
	require.NoError(t, err)
	allImage, err := DecodeAllImage(img)
	require.NoError(t, err)
	for _, got := range [][]string{allPath, allReader, allImage} {
		assert.Equal(t, allBytes, got)
	}
}

// TestDecodeFileNotFound AC-I2：路径指向不存在文件 → 文件类错误（可与无码错误区分）。
func TestDecodeFileNotFound(t *testing.T) {
	_, err := DecodeFile(filepath.Join(t.TempDir(), "no-such-file.png"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrNotExist), "want os.ErrNotExist, got %v", err)
	assert.False(t, errors.Is(err, ErrNotFound))
}
