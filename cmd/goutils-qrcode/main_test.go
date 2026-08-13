package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	qrcode "github.com/fasionchan/goutils/libs/qrcode"
	sqrcode "github.com/skip2/go-qrcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runCLI 以内存流执行 CLI，返回退出码与输出。
func runCLI(t *testing.T, stdin string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errBuf bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errBuf)
	return code, out.String(), errBuf.String()
}

// writePNGForTest 生成一张包含 content 的二维码 png 文件，返回路径。
func writePNGForTest(t *testing.T, dir, name, content string) string {
	t.Helper()
	data, err := qrcode.EncodePNG(content)
	require.NoError(t, err)
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, data, 0o644))
	return path
}

// writeMultiCodePNGForTest 生成包含两个二维码的 png 文件，返回路径。
func writeMultiCodePNGForTest(t *testing.T, dir, name string) string {
	t.Helper()
	q1, err := sqrcode.New("payload-alpha", sqrcode.Medium)
	require.NoError(t, err)
	q2, err := sqrcode.New("payload-beta", sqrcode.Medium)
	require.NoError(t, err)
	const side, gap, margin = 150, 40, 40
	canvas := image.NewRGBA(image.Rect(0, 0, side*2+gap+margin*2, side+margin*2))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(margin, margin, margin+side, margin+side), q1.Image(side), image.Point{}, draw.Src)
	draw.Draw(canvas, image.Rect(margin+side+gap, margin, margin+side+gap+side, margin+side), q2.Image(side), image.Point{}, draw.Src)
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, canvas))
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, buf.Bytes(), 0o644))
	return path
}

// TestCLIDecodeFile AC-I3：给定合法 png/jpeg 路径，stdout 一行 payload，exit 0。
func TestCLIDecodeFile(t *testing.T) {
	dir := t.TempDir()
	pngPath := writePNGForTest(t, dir, "qr.png", "hello cli")

	code, stdout, stderr := runCLI(t, "", "decode", pngPath)
	assert.Equal(t, 0, code)
	assert.Equal(t, "hello cli\n", stdout)
	assert.Empty(t, stderr, "success path must not write noise to stderr")

	// jpeg 同样支持。
	jpegBytes, err := qrcode.EncodeJPEG("jpeg cli")
	require.NoError(t, err)
	jpegPath := filepath.Join(dir, "qr.jpg")
	require.NoError(t, os.WriteFile(jpegPath, jpegBytes, 0o644))
	code, stdout, _ = runCLI(t, "", "decode", jpegPath)
	assert.Equal(t, 0, code)
	assert.Equal(t, "jpeg cli\n", stdout)
}

// TestCLIDecodeStdin AC-I3 补充：stdin 输入同样支持（路径缺省或 "-"）。
func TestCLIDecodeStdin(t *testing.T) {
	data, err := qrcode.EncodePNG("from stdin")
	require.NoError(t, err)

	code, stdout, _ := runCLI(t, string(data), "decode", "-")
	assert.Equal(t, 0, code)
	assert.Equal(t, "from stdin\n", stdout)

	code, stdout, _ = runCLI(t, string(data), "decode")
	assert.Equal(t, 0, code)
	assert.Equal(t, "from stdin\n", stdout)
}

// TestCLIDecodeAll AC-I4：多码样例 stdout 行数 = 码数，payload 集合正确，exit 0。
func TestCLIDecodeAll(t *testing.T) {
	dir := t.TempDir()
	path := writeMultiCodePNGForTest(t, dir, "multi.png")

	code, stdout, _ := runCLI(t, "", "decode", "--all", path)
	assert.Equal(t, 0, code)
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	assert.Len(t, lines, 2)
	assert.ElementsMatch(t, []string{"payload-alpha", "payload-beta"}, lines)
}

// TestCLIEncodePNG AC-I5：encode -f png 产出文件可被标准库解码，且可被本库还原。
func TestCLIEncodePNG(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.png")
	code, _, stderr := runCLI(t, "", "encode", "-f", "png", "-o", out, "cli encode png")
	assert.Equal(t, 0, code, "stderr: %s", stderr)

	data, err := os.ReadFile(out)
	require.NoError(t, err)
	_, err = png.Decode(bytes.NewReader(data))
	assert.NoError(t, err, "output must be valid png")
	got, err := qrcode.DecodeBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "cli encode png", got)
}

// TestCLIEncodeJPEG AC-I5：encode -f jpeg（默认 quality 90）产出文件可被解码还原。
func TestCLIEncodeJPEG(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.jpg")
	code, _, stderr := runCLI(t, "", "encode", "-f", "jpeg", "-o", out, "cli encode jpeg")
	assert.Equal(t, 0, code, "stderr: %s", stderr)

	data, err := os.ReadFile(out)
	require.NoError(t, err)
	img, err := jpeg.Decode(bytes.NewReader(data))
	assert.NoError(t, err, "output must be valid jpeg")
	assert.Equal(t, 256, img.Bounds().Dx())
	got, err := qrcode.DecodeBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "cli encode jpeg", got)
}

// TestCLIEncodeStdout AC-I5 补充：-o - 输出到 stdout。
func TestCLIEncodeStdout(t *testing.T) {
	code, stdout, _ := runCLI(t, "", "encode", "-o", "-", "stdout encode")
	assert.Equal(t, 0, code)
	data := []byte(stdout)
	_, err := png.Decode(bytes.NewReader(data))
	assert.NoError(t, err, "stdout must be valid png bytes")
	got, err := qrcode.DecodeBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "stdout encode", got)
}

// TestCLIEncodeDefaultFile AC-I5 补充：未指定 -o 时默认写出 qrcode.<ext>。
func TestCLIEncodeDefaultFile(t *testing.T) {
	oldwd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(oldwd)
	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))

	code, _, stderr := runCLI(t, "", "encode", "default file name")
	assert.Equal(t, 0, code, "stderr: %s", stderr)
	data, err := os.ReadFile(filepath.Join(dir, "qrcode.png"))
	require.NoError(t, err)
	got, err := qrcode.DecodeBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "default file name", got)
}

// TestCLIPrint AC-I6：print 输出非空多行终端码，exit 0。
func TestCLIPrint(t *testing.T) {
	code, stdout, stderr := runCLI(t, "", "print", "terminal text")
	assert.Equal(t, 0, code)
	assert.NotEmpty(t, stdout)
	assert.Greater(t, strings.Count(stdout, "\n"), 1)
	assert.Empty(t, stderr)

	// stdin 输入。
	code, stdout, _ = runCLI(t, "from stdin", "print")
	assert.Equal(t, 0, code)
	assert.NotEmpty(t, stdout)
}

// TestCLIErrors AC-I7：错误输入下 exit != 0，且 stderr 含错误信息。
func TestCLIErrors(t *testing.T) {
	dir := t.TempDir()

	t.Run("bad path", func(t *testing.T) {
		code, _, stderr := runCLI(t, "", "decode", filepath.Join(dir, "missing.png"))
		assert.Equal(t, 2, code)
		assert.NotEmpty(t, stderr)
	})

	t.Run("no code", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
		var buf bytes.Buffer
		require.NoError(t, png.Encode(&buf, img))
		path := filepath.Join(dir, "blank.png")
		require.NoError(t, os.WriteFile(path, buf.Bytes(), 0o644))

		code, _, stderr := runCLI(t, "", "decode", path)
		assert.Equal(t, 3, code)
		assert.Contains(t, stderr, "no qr code found")

		code, _, stderr = runCLI(t, "", "decode", "--all", path)
		assert.Equal(t, 3, code, "--all with no code must exit 3")
		assert.Contains(t, stderr, "no qr code found")
	})

	t.Run("unsupported format", func(t *testing.T) {
		path := filepath.Join(dir, "not-image.txt")
		require.NoError(t, os.WriteFile(path, []byte("not an image"), 0o644))
		code, _, stderr := runCLI(t, "", "decode", path)
		assert.Equal(t, 2, code)
		assert.NotEmpty(t, stderr)
	})

	t.Run("invalid level", func(t *testing.T) {
		code, _, stderr := runCLI(t, "", "encode", "--level", "X", "-o", "-", "x")
		assert.Equal(t, 1, code)
		assert.NotEmpty(t, stderr)
	})

	t.Run("invalid format", func(t *testing.T) {
		code, _, stderr := runCLI(t, "", "encode", "-f", "gif", "-o", "-", "x")
		assert.Equal(t, 1, code)
		assert.NotEmpty(t, stderr)
	})

	t.Run("unknown command", func(t *testing.T) {
		code, _, stderr := runCLI(t, "", "frobnicate")
		assert.Equal(t, 1, code)
		assert.NotEmpty(t, stderr)
	})

	t.Run("empty content encode", func(t *testing.T) {
		code, _, stderr := runCLI(t, "", "encode", "-o", "-", "")
		assert.Equal(t, 3, code)
		assert.Contains(t, stderr, "invalid content")
	})
}

// TestCLIBinarySmoke 冒烟测试：真实编译出二进制并用子进程执行 decode / encode / print。
func TestCLIBinarySmoke(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "goutils-qrcode")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = "."
	out, err := build.CombinedOutput()
	require.NoError(t, err, "build failed: %s", out)

	dir := t.TempDir()
	pngPath := writePNGForTest(t, dir, "qr.png", "smoke test")

	// decode
	cmd := exec.Command(bin, "decode", pngPath)
	got, err := cmd.Output()
	require.NoError(t, err, "decode exit: %v", err)
	assert.Equal(t, "smoke test\n", string(got))

	// encode -o -
	cmd = exec.Command(bin, "encode", "-o", "-", "smoke encode")
	data, err := cmd.Output()
	require.NoError(t, err)
	decoded, err := qrcode.DecodeBytes(data)
	require.NoError(t, err)
	assert.Equal(t, "smoke encode", decoded)

	// print
	cmd = exec.Command(bin, "print", "smoke print")
	got, err = cmd.Output()
	require.NoError(t, err)
	assert.NotEmpty(t, got)

	// 错误路径：exit code 非 0
	cmd = exec.Command(bin, "decode", filepath.Join(dir, "missing.png"))
	runErr := cmd.Run()
	var exitErr *exec.ExitError
	require.ErrorAs(t, runErr, &exitErr)
	assert.NotEqual(t, 0, exitErr.ExitCode())
}
