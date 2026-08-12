package qrcode

import (
	"fmt"
	"io"

	"github.com/mdp/qrterminal/v3"
	rqr "rsc.io/qr"
)

// TerminalOptions 是终端输出的配置项。
type TerminalOptions struct {
	// Level 纠错级别，默认 LevelM。
	Level RecoveryLevel
	// HalfBlock 是否使用半块（half-block）字符渲染，默认 true；
	// false 时使用全块字符。nil 表示使用默认值。
	HalfBlock *bool
	// QuietZone 静区（quiet zone）宽度（模块数），默认 DefaultQuietZone。
	QuietZone int
	// MaxColumns 允许输出的最大终端列宽；大于 0 时若码图宽度超出
	// 该值则返回 ErrTooNarrow（不输出截断内容）。小于等于 0 表示不限制。
	MaxColumns int
}

// normalized 把零值字段替换为默认值。
func (o TerminalOptions) normalized() TerminalOptions {
	if o.Level == 0 {
		o.Level = LevelM
	}
	if o.QuietZone <= 0 {
		o.QuietZone = DefaultQuietZone
	}
	return o
}

// validate 校验归一化后的选项取值。
func (o TerminalOptions) validate() error {
	if !o.Level.Valid() {
		return fmt.Errorf("%w: invalid recovery level %d", ErrInvalidOption, o.Level)
	}
	if o.QuietZone < 1 || o.QuietZone > 64 {
		return fmt.Errorf("%w: quiet zone must be in [1, 64], got %d", ErrInvalidOption, o.QuietZone)
	}
	return nil
}

// halfBlock 返回 half-block 开关的实际取值。
func (o TerminalOptions) halfBlock() bool {
	if o.HalfBlock != nil {
		return *o.HalfBlock
	}
	return true
}

// normalizeTerminalOptions 处理可变参数：不传时使用全部默认值。
func normalizeTerminalOptions(opts []TerminalOptions) TerminalOptions {
	if len(opts) == 0 {
		return (TerminalOptions{}).normalized()
	}
	return opts[0].normalized()
}

// outputWidth 计算渲染结果的终端列宽：
//   - half-block：每个模块占 1 列，宽度 = 码图模块数 + 2 * QuietZone；
//   - 全块：默认字符每个模块占 2 列。
func (o TerminalOptions) outputWidth(codeSize int) int {
	width := codeSize + 2*o.QuietZone
	if !o.halfBlock() {
		width *= 2
	}
	return width
}

// EncodeToTerminal 把 content 渲染为终端二维码并写入 w。
//
// 输出字符集：
//   - half-block（默认）："▀ ▄ █"（每个字符 1 列）；
//   - 全块：ANSI 背景色块（每个模块 2 列）。
//
// 当 MaxColumns 大于 0 且输出宽度超出该值时，返回 ErrTooNarrow，
// 且不会向 w 写入任何内容。
func EncodeToTerminal(content string, w io.Writer, opts ...TerminalOptions) error {
	o := normalizeTerminalOptions(opts)
	if err := o.validate(); err != nil {
		return err
	}
	if content == "" {
		return ErrInvalidContent
	}

	// 先编码一次以校验内容容量并计算码图尺寸；
	// qrterminal 上游对编码错误不返回错误，故必须在此处拦截。
	code, err := rqr.Encode(content, o.Level.terminal())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidContent, err)
	}

	if o.MaxColumns > 0 {
		if width := o.outputWidth(code.Size); width > o.MaxColumns {
			return fmt.Errorf("%w: needs %d columns, max is %d", ErrTooNarrow, width, o.MaxColumns)
		}
	}

	config := qrterminal.Config{
		Level:      o.Level.terminal(),
		Writer:     w,
		HalfBlocks: o.halfBlock(),
		QuietZone:  o.QuietZone,
	}
	qrterminal.GenerateWithConfig(content, config)
	return nil
}
