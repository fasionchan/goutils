// Package qrcode 提供纯 Go（无 cgo）实现的二维码（QR Code）处理能力：
//
//   - 识别（decode）：从 png / jpeg 图片中解析二维码内容，输入支持
//     文件路径、[]byte、io.Reader 与 image.Image 四种形态，并提供
//     Decode（单码）与 DecodeAll（多码）两组入口；
//   - 终端输出（terminal）：把字符串渲染成可扫描的终端二维码；
//   - 生成（encode）：把字符串生成 png / jpeg 图片字节流或 image.Image。
//
// 内部基于 gozxing（识别）、skip2/go-qrcode（生成）与 mdp/qrterminal（终端）
// 三个开源库封装，公开 API 不泄露底层库类型。
package qrcode

import (
	"errors"
	"fmt"
	"strings"
)

// 本包暴露的错误哨兵值（sentinel），可用 errors.Is / errors.As 断言。
var (
	// ErrNotFound 表示图片中未识别到任何二维码。
	ErrNotFound = errors.New("qrcode: no qr code found")

	// ErrDecodeFailed 表示图片中存在疑似二维码但解码失败（如损坏至不可纠）。
	ErrDecodeFailed = errors.New("qrcode: failed to decode qr code")

	// ErrUnsupportedFormat 表示输入字节流不是 png/jpeg 图片或无法解码。
	ErrUnsupportedFormat = errors.New("qrcode: unsupported image format (only png/jpeg)")

	// ErrInvalidContent 表示内容为空或超出二维码容量。
	ErrInvalidContent = errors.New("qrcode: invalid content (empty or too long)")

	// ErrInvalidOption 表示配置选项非法。
	ErrInvalidOption = errors.New("qrcode: invalid option")

	// ErrTooNarrow 表示终端列宽不足以完整显示二维码。
	ErrTooNarrow = errors.New("qrcode: output width exceeds MaxColumns")
)

// RecoveryLevel 表示二维码的纠错级别：LevelL < LevelM < LevelQ < LevelH。
// 纠错能力越强，可容纳的数据越少。
type RecoveryLevel int

const (
	// LevelL 约可恢复 7% 的数据损坏。
	LevelL RecoveryLevel = iota + 1
	// LevelM 约可恢复 15% 的数据损坏，为默认级别。
	LevelM
	// LevelQ 约可恢复 25% 的数据损坏。
	LevelQ
	// LevelH 约可恢复 30% 的数据损坏。
	LevelH
)

var recoveryLevelNames = map[RecoveryLevel]string{
	LevelL: "L",
	LevelM: "M",
	LevelQ: "Q",
	LevelH: "H",
}

// Valid 报告级别是否为合法取值之一。
func (l RecoveryLevel) Valid() bool {
	_, ok := recoveryLevelNames[l]
	return ok
}

// String 返回级别的单字母表示：L / M / Q / H。
func (l RecoveryLevel) String() string {
	if name, ok := recoveryLevelNames[l]; ok {
		return name
	}
	return fmt.Sprintf("RecoveryLevel(%d)", int(l))
}

// ParseRecoveryLevel 解析纠错级别字符串，接受 L / M / Q / H（不区分大小写）。
func ParseRecoveryLevel(s string) (RecoveryLevel, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "L":
		return LevelL, nil
	case "M":
		return LevelM, nil
	case "Q":
		return LevelQ, nil
	case "H":
		return LevelH, nil
	}
	return 0, fmt.Errorf("%w: unknown recovery level %q (expect L/M/Q/H)", ErrInvalidOption, s)
}

// 默认选项取值。
const (
	// DefaultSize 是生成图片的默认像素边长。
	DefaultSize = 256
	// DefaultQuality 是生成 jpeg 的默认质量（1-100）。
	DefaultQuality = 90
	// DefaultQuietZone 是终端输出默认的静区（quiet zone）模块数。
	DefaultQuietZone = 4
)
