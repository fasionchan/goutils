package qrcode

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEncodeToTerminalNonEmpty AC-T1：写入 io.Writer 的输出非空、多行、
// 字符集符合默认 half-block 配置（▀▄█ 及空格）。
func TestEncodeToTerminalNonEmpty(t *testing.T) {
	var buf bytes.Buffer
	err := EncodeToTerminal("terminal payload", &buf)
	require.NoError(t, err)

	out := buf.String()
	assert.NotEmpty(t, out)
	assert.Greater(t, strings.Count(out, "\n"), 1, "output must be multi-line")
	for _, r := range out {
		switch r {
		case '▀', '▄', '█', ' ', '\n':
		default:
			t.Fatalf("unexpected character %q in half-block output", r)
		}
	}
}

// TestEncodeToTerminalDeterministic AC-T2：相同 content + Config → 输出字节序列一致。
func TestEncodeToTerminalDeterministic(t *testing.T) {
	var first, second bytes.Buffer
	require.NoError(t, EncodeToTerminal("deterministic", &first))
	require.NoError(t, EncodeToTerminal("deterministic", &second))
	assert.Equal(t, first.Bytes(), second.Bytes())
}

// TestEncodeToTerminalConfigEffect AC-T3：可配纠错级别与 half-block 开关；
// 改变配置可观测到输出差异。
func TestEncodeToTerminalConfigEffect(t *testing.T) {
	content := "config effect"
	var levelL, levelM bytes.Buffer
	require.NoError(t, EncodeToTerminal(content, &levelL, TerminalOptions{Level: LevelL}))
	require.NoError(t, EncodeToTerminal(content, &levelM, TerminalOptions{Level: LevelM}))
	assert.NotEqual(t, levelL.Bytes(), levelM.Bytes(),
		"different levels must produce different output")

	half := true
	var halfBuf, fullBuf bytes.Buffer
	require.NoError(t, EncodeToTerminal(content, &halfBuf, TerminalOptions{HalfBlock: &half}))
	full := false
	require.NoError(t, EncodeToTerminal(content, &fullBuf, TerminalOptions{HalfBlock: &full}))
	assert.NotEqual(t, halfBuf.Bytes(), fullBuf.Bytes(),
		"half-block vs full-block must produce different output")
	assert.NotEmpty(t, fullBuf.String(), "full-block output must be non-empty")
}

// TestEncodeToTerminalTooNarrow AC-T4：极窄列宽场景 → 返回明确错误，
// 且不向 writer 写入任何内容。
func TestEncodeToTerminalTooNarrow(t *testing.T) {
	var buf bytes.Buffer
	err := EncodeToTerminal("this content produces a wide qr code", &buf, TerminalOptions{MaxColumns: 10})
	assert.ErrorIs(t, err, ErrTooNarrow)
	assert.Empty(t, buf.Bytes(), "no output may be written when too narrow")

	// 足够宽的 MaxColumns 应成功。
	var okBuf bytes.Buffer
	err = EncodeToTerminal("narrow ok", &okBuf, TerminalOptions{MaxColumns: 1000})
	assert.NoError(t, err)
	assert.NotEmpty(t, okBuf.Bytes())
}

// TestEncodeToTerminalEmptyContent AC 补充：空内容 → ErrInvalidContent。
func TestEncodeToTerminalEmptyContent(t *testing.T) {
	var buf bytes.Buffer
	err := EncodeToTerminal("", &buf)
	assert.ErrorIs(t, err, ErrInvalidContent)
	assert.Empty(t, buf.Bytes())
}

// TestEncodeToTerminalInvalidOption AC 补充：非法级别 → ErrInvalidOption。
func TestEncodeToTerminalInvalidOption(t *testing.T) {
	var buf bytes.Buffer
	err := EncodeToTerminal("x", &buf, TerminalOptions{Level: RecoveryLevel(99)})
	assert.ErrorIs(t, err, ErrInvalidOption)
}
