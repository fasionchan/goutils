// Command goutils-qrcode 是基于 libs/qrcode 的二维码命令行工具，
// 提供 decode（识别）、encode（生成）、print（终端输出）三个子命令。
//
// 退出码：0 成功 · 1 用法错误 · 2 输入/IO/格式错误 · 3 业务失败。
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	qrcode "github.com/fasionchan/goutils/libs/qrcode"
)

const program = "goutils-qrcode"

const rootUsage = `usage: goutils-qrcode <command> [options] [args]

commands:
  decode   识别图片中的二维码，输出 payload（支持 --all 解码全部）
  encode   把文本生成二维码图片（png / jpeg）
  print    把文本渲染为终端二维码

run 'goutils-qrcode <command> -h' for command help
`

const decodeUsage = `usage: goutils-qrcode decode [--all] [path|-]

decode 从图片中识别二维码内容，一行输出一个 payload。
path 缺省或为 "-" 时从 stdin 读取。

options:
  --all    解码图片中的全部二维码，每行一个 payload
  -h       显示帮助
`

const encodeUsage = `usage: goutils-qrcode encode [-f png|jpeg] [-o path|-] [--size N] [--level L|M|Q|H] [--quality N] [text]

encode 把 text（缺省时从 stdin 读取）生成二维码图片并写出。
默认输出文件 qrcode.png（-o - 时输出到 stdout）。

options:
  -f string        输出格式 png 或 jpeg（默认 png）
  -o string        输出路径；"-" 表示 stdout（默认 qrcode.<ext>）
  --size int       图片像素边长（默认 256）
  --level string   纠错级别 L/M/Q/H（默认 M）
  --quality int    jpeg 质量 1-100（默认 90）
  -h               显示帮助
`

const printUsage = `usage: goutils-qrcode print [--level L|M|Q|H] [--no-half-block] [text]

print 把 text（缺省时从 stdin 读取）渲染为终端二维码并输出到 stdout。

options:
  --level string      纠错级别 L/M/Q/H（默认 M）
  --no-half-block     使用全块字符代替默认的 half-block
  -h                  显示帮助
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run 执行 CLI，返回进程退出码。args 不含程序名。
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, rootUsage)
		return 1
	}

	switch args[0] {
	case "decode":
		return runDecode(args[1:], stdin, stdout, stderr)
	case "encode":
		return runEncode(args[1:], stdin, stdout, stderr)
	case "print":
		return runPrint(args[1:], stdin, stdout, stderr)
	case "help", "-h", "--help":
		fmt.Fprint(stdout, rootUsage)
		return 0
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n\n", program, args[0])
		fmt.Fprint(stderr, rootUsage)
		return 1
	}
}

// parseFlags 解析 flags；返回 false 表示应直接退出（-h 为 0，解析错误为 1）。
func parseFlags(fs *flag.FlagSet, args []string, stderr io.Writer) (int, bool) {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0, false
		}
		return 1, false
	}
	return 0, true
}

// readInput 从位置参数或 stdin 读取文本内容。
func readInput(fs *flag.FlagSet, stdin io.Reader) (string, error) {
	if fs.NArg() > 1 {
		return "", fmt.Errorf("too many arguments")
	}
	if fs.NArg() == 1 {
		return fs.Arg(0), nil
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	return string(data), nil
}

// exitCodeForError 把库错误映射为 CLI 退出码：
// 业务失败（哨兵错误）→ 3；输入/IO/格式（文件缺失、格式错误、未知错误）→ 2。
func exitCodeForError(err error) int {
	switch {
	case errors.Is(err, qrcode.ErrNotFound),
		errors.Is(err, qrcode.ErrDecodeFailed),
		errors.Is(err, qrcode.ErrInvalidContent),
		errors.Is(err, qrcode.ErrInvalidOption),
		errors.Is(err, qrcode.ErrTooNarrow):
		return 3
	default:
		return 2
	}
}

// reportError 向 stderr 打印错误并返回对应退出码。
func reportError(stderr io.Writer, cmd string, err error) int {
	code := exitCodeForError(err)
	fmt.Fprintf(stderr, "%s: %s: %v\n", program, cmd, err)
	return code
}

// parseLevelFlag 解析 --level 参数；非法时返回用法错误码 1。
func parseLevelFlag(stderr io.Writer, cmd, value string) (qrcode.RecoveryLevel, int, bool) {
	if value == "" {
		return 0, 0, true
	}
	level, err := qrcode.ParseRecoveryLevel(value)
	if err != nil {
		fmt.Fprintf(stderr, "%s: %s: %v\n", program, cmd, err)
		return 0, 1, false
	}
	return level, 0, true
}

func runDecode(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("decode", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, decodeUsage) }
	all := fs.Bool("all", false, "decode all qr codes, one payload per line")
	if code, ok := parseFlags(fs, args, stderr); !ok {
		return code
	}
	if fs.NArg() > 1 {
		fmt.Fprintf(stderr, "%s: decode: too many arguments\n\n", program)
		fmt.Fprint(stderr, decodeUsage)
		return 1
	}

	var (
		payloads []string
		err      error
	)
	if fs.NArg() == 1 && fs.Arg(0) != "-" {
		if *all {
			payloads, err = qrcode.DecodeAllFile(fs.Arg(0))
		} else {
			var one string
			one, err = qrcode.DecodeFile(fs.Arg(0))
			if err == nil {
				payloads = []string{one}
			}
		}
	} else {
		data, readErr := io.ReadAll(stdin)
		if readErr != nil {
			fmt.Fprintf(stderr, "%s: decode: read stdin: %v\n", program, readErr)
			return 2
		}
		if *all {
			payloads, err = qrcode.DecodeAllBytes(data)
		} else {
			var one string
			one, err = qrcode.DecodeBytes(data)
			if err == nil {
				payloads = []string{one}
			}
		}
	}
	if err != nil {
		return reportError(stderr, "decode", err)
	}
	for _, payload := range payloads {
		fmt.Fprintln(stdout, payload)
	}
	return 0
}

func runEncode(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("encode", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, encodeUsage) }
	format := fs.String("f", "png", "output format: png or jpeg")
	output := fs.String("o", "", "output file path; '-' means stdout")
	size := fs.Int("size", 0, "image width/height in pixels (default 256)")
	level := fs.String("level", "", "error correction level: L/M/Q/H (default M)")
	quality := fs.Int("quality", 0, "jpeg quality 1-100 (default 90)")
	if code, ok := parseFlags(fs, args, stderr); !ok {
		return code
	}
	content, err := readInput(fs, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "%s: encode: %v\n", program, err)
		return 2
	}

	parsedLevel, code, ok := parseLevelFlag(stderr, "encode", *level)
	if !ok {
		return code
	}

	opts := qrcode.EncodeOptions{Level: parsedLevel}
	if *size != 0 {
		opts.Size = *size
	}
	if *quality != 0 {
		opts.Quality = *quality
	}

	var data []byte
	switch strings.ToLower(*format) {
	case "png":
		data, err = qrcode.EncodePNG(content, opts)
	case "jpeg", "jpg":
		data, err = qrcode.EncodeJPEG(content, opts)
	default:
		fmt.Fprintf(stderr, "%s: encode: unsupported format %q (expect png or jpeg)\n", program, *format)
		return 1
	}
	if err != nil {
		return reportError(stderr, "encode", err)
	}

	if *output == "-" {
		if _, err := stdout.Write(data); err != nil {
			fmt.Fprintf(stderr, "%s: encode: write stdout: %v\n", program, err)
			return 2
		}
		return 0
	}
	path := *output
	if path == "" {
		path = "qrcode." + strings.ToLower(*format)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Fprintf(stderr, "%s: encode: write %s: %v\n", program, path, err)
		return 2
	}
	return 0
}

func runPrint(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("print", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, printUsage) }
	level := fs.String("level", "", "error correction level: L/M/Q/H (default M)")
	noHalfBlock := fs.Bool("no-half-block", false, "use full-block characters instead of half blocks")
	if code, ok := parseFlags(fs, args, stderr); !ok {
		return code
	}
	content, err := readInput(fs, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "%s: print: %v\n", program, err)
		return 2
	}

	parsedLevel, code, ok := parseLevelFlag(stderr, "print", *level)
	if !ok {
		return code
	}

	halfBlock := !*noHalfBlock
	opts := qrcode.TerminalOptions{
		Level:     parsedLevel,
		HalfBlock: &halfBlock,
	}
	if err := qrcode.EncodeToTerminal(content, stdout, opts); err != nil {
		return reportError(stderr, "print", err)
	}
	return 0
}
