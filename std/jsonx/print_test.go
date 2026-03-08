/*
 * Author: fasion
 * Created time: 2026-03-07 13:11:59
 * Last Modified by: fasion
 * Last Modified time: 2026-03-07 15:24:51
 */

package jsonx

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	datas := []any{
		map[string]string{"foo": "bar"},
		map[string][]string{"foo": {"bar", "baz"}},
		map[string]map[string]string{"foo": {"bar": "baz"}},
		map[string]map[string][]string{"foo": {"bar": {"baz", "qux"}}},
		map[string]map[string]map[string]string{"foo": {"bar": {"baz": "qux"}}},
	}

	var buf bytes.Buffer
	FprintPro(&buf, "  ", "", "[%d] ", "", "", "\n", "", datas...)
	out := buf.String()

	// 每个元素应有 [0]、[1] 等前缀且为合法 JSON
	for i := range datas {
		marker := "[" + string(rune('0'+i)) + "]"
		if i >= 10 {
			// 仅校验前 10 个
			break
		}
		if !strings.Contains(out, marker) {
			t.Errorf("output should contain prefix %q", marker)
		}
	}
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) < len(datas) {
		t.Errorf("expected at least %d lines, got %d", len(datas), len(lines))
	}
}

func TestFprintPro_Simple(t *testing.T) {
	var buf bytes.Buffer
	arg := map[string]string{"a": "b"}
	FprintPro(&buf, "", "", "", "", "", "", "", arg)
	out := strings.TrimSpace(buf.String())
	var got map[string]string
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output should be valid JSON: %v\noutput: %s", err, out)
	}
	if got["a"] != "b" {
		t.Errorf("got %v, want map[a:b]", got)
	}
}

func TestFprintPro_WithIndent(t *testing.T) {
	var buf bytes.Buffer
	arg := map[string]string{"x": "y"}
	FprintPro(&buf, "  ", "", "", "", "", "", "", arg)
	out := buf.String()
	if !strings.Contains(out, "\n  ") {
		t.Errorf("output should be indented with 2 spaces, got: %q", out)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &got); err != nil {
		t.Fatalf("output should be valid JSON: %v", err)
	}
}

func TestFprintPro_WithPrefix(t *testing.T) {
	var buf bytes.Buffer
	FprintPro(&buf, "", ">> ", "", "", "", "", "", map[string]string{"k": "v"})
	out := buf.String()
	if !strings.HasPrefix(out, ">> ") {
		t.Errorf("output should start with prefix >> , got: %q", out)
	}
	if !strings.Contains(out, `"k":"v"`) && !strings.Contains(out, `"k": "v"`) {
		t.Errorf("output should contain JSON: %q", out)
	}
}

func TestFprintPro_WithPrefixFmt(t *testing.T) {
	var buf bytes.Buffer
	FprintPro(&buf, "", "", "[%d] ", "", "", "", "", "a", "b", "c")
	out := buf.String()
	for i := 0; i < 3; i++ {
		marker := "[" + string(rune('0'+i)) + "]"
		if !strings.Contains(out, marker) {
			t.Errorf("output should contain %q", marker)
		}
	}
}

func TestFprintPro_WithHeadAndTail(t *testing.T) {
	var buf bytes.Buffer
	FprintPro(&buf, "", "", "", "HEAD\n", "", "TAIL", "", map[string]string{"x": "1"})
	out := buf.String()
	if !strings.HasPrefix(out, "HEAD\n") {
		t.Errorf("output should start with HEAD\\n, got: %q", out[:min(20, len(out))])
	}
	if !strings.Contains(out, "TAIL") {
		t.Errorf("output should contain TAIL, got: %q", out)
	}
}

func TestFprintPro_WithHeadFmtAndTailFmt(t *testing.T) {
	var buf bytes.Buffer
	// head/headFmt: 第 5、6 个参数；tail/tailFmt: 第 7、8 个。格式串放 Fmt 里才会 Sprintf
	FprintPro(&buf, "", "", "", "", "=== %d ===\n", "", "\n=== end %d ===", "first", "second")
	out := buf.String()
	if !strings.Contains(out, "=== 0 ===") || !strings.Contains(out, "=== 1 ===") {
		t.Errorf("output should contain headFmt for index 0 and 1: %q", out)
	}
	if !strings.Contains(out, "=== end 0 ===") || !strings.Contains(out, "=== end 1 ===") {
		t.Errorf("output should contain tailFmt for index 0 and 1: %q", out)
	}
}

func TestFprintPro_EncodeError(t *testing.T) {
	var buf bytes.Buffer
	ch := make(chan int) // channel 无法被 JSON 编码
	FprintPro(&buf, "", "", "", "", "", "", "", ch)
	out := buf.String()
	if out == "" {
		t.Error("FprintPro should write error message to writer when Encode fails")
	}
	if !strings.Contains(out, "json") && !strings.Contains(out, "channel") && !strings.Contains(out, "unsupported") {
		// 至少应包含错误信息
		if len(out) < 5 {
			t.Errorf("expected error message in output, got: %q", out)
		}
	}
}

func TestFprintPro_EmptyArgs(t *testing.T) {
	var buf bytes.Buffer
	FprintPro(&buf, "  ", "", "", "", "", "", "")
	if buf.Len() != 0 {
		t.Errorf("empty args should produce no output, got %d bytes", buf.Len())
	}
}

func TestGetPrintFunc(t *testing.T) {
	fn := GetPrintFunc("  ", "", "[%d] ", "", "", "\n", "")
	if fn == nil {
		t.Fatal("GetPrintFunc should not return nil")
	}
	// 调用不 panic 即可（实际写 stdout，不便于断言）
	fn(map[string]string{"test": "value"})
}
