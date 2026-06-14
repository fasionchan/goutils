/*
 * Author: fasion
 * Created time: 2026-06-14
 */

package mdutils

import (
	"strings"
	"testing"
)

func TestTruncateMarkdown_NoTruncation(t *testing.T) {
	content := "Hello, world!"
	result, truncated := TruncateMarkdown(content, 100)
	if truncated {
		t.Fatal("expected not truncated")
	}
	if result != content {
		t.Fatalf("expected %q, got %q", content, result)
	}
}

func TestTruncateMarkdown_EmptyInput(t *testing.T) {
	result, truncated := TruncateMarkdown("", 100)
	if truncated {
		t.Fatal("expected not truncated for empty input")
	}
	if result != "" {
		t.Fatalf("expected empty, got %q", result)
	}
}

func TestTruncateMarkdown_WhitespaceOnly(t *testing.T) {
	content := "   \n\t  "
	result, truncated := TruncateMarkdown(content, 100)
	if truncated {
		t.Fatal("expected not truncated")
	}
	if result != content {
		t.Fatalf("expected %q, got %q", content, result)
	}
}

func TestTruncateMarkdown_MultiParagraph(t *testing.T) {
	content := "Para1\n\nPara2\n\nPara3"
	result, truncated := TruncateMarkdown(content, 8)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 8 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
	if !strings.HasPrefix(result, "Para1") {
		t.Fatalf("expected to start with Para1, got %q", result)
	}
	if strings.HasSuffix(result, "...") {
		t.Fatalf("TruncateMarkdown should not append suffix, got %q", result)
	}
	if strings.Contains(result, "Para2") {
		t.Fatalf("Para2 should be truncated away, got %q", result)
	}
}

func TestTruncateMarkdown_CodeBlock(t *testing.T) {
	content := "Before\n\n```\nline1\nline2\nline3\n```\n\nAfter"
	result, truncated := TruncateMarkdown(content, 30)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 30 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
	if strings.Contains(result, "```") && !strings.Contains(result, "line1\nline2\nline3") {
		if strings.Count(result, "```") != 2 {
			t.Fatalf("incomplete code fence in result: %q", result)
		}
	}
}

func TestTruncateMarkdown_List(t *testing.T) {
	content := "- item1\n- item2\n- item3\n- item4"
	result, truncated := TruncateMarkdown(content, 25)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 25 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
	if strings.Contains(result, "item1") && strings.Contains(result, "item4") {
		t.Fatalf("should not contain both first and last items: %q", result)
	}
}

func TestTruncateMarkdown_OrderedList(t *testing.T) {
	content := "1. first\n2. second\n3. third"
	result, truncated := TruncateMarkdown(content, 20)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 20 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
}

func TestTruncateMarkdown_OversizedSingleNode(t *testing.T) {
	content := strings.Repeat("A", 100)
	result, truncated := TruncateMarkdown(content, 20)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 20 {
		t.Fatalf("result exceeds maxBytes: len=%d", len(result))
	}
	if strings.HasSuffix(result, "...") {
		t.Fatalf("TruncateMarkdown should not append suffix, got %q", result)
	}
	if len(result) != 20 {
		t.Fatalf("expected 20 bytes content, got %d: %q", len(result), result)
	}
}

func TestTruncateMarkdown_MaxBytesZero(t *testing.T) {
	result, truncated := TruncateMarkdown("hello", 0)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if result != "" {
		t.Fatalf("expected empty, got %q", result)
	}

	result, truncated = TruncateMarkdown("", 0)
	if truncated {
		t.Fatal("empty content should not be truncated")
	}
	if result != "" {
		t.Fatalf("expected empty, got %q", result)
	}
}

func TestTruncateMarkdown_MaxBytesNegative(t *testing.T) {
	result, truncated := TruncateMarkdown("hello", -1)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if result != "" {
		t.Fatalf("expected empty, got %q", result)
	}
}

func TestTruncateMarkdown_UTF8Safe(t *testing.T) {
	content := "中文测试内容很长"
	result, truncated := TruncateMarkdown(content, 15)
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 15 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
	if strings.HasSuffix(result, "...") {
		t.Fatalf("TruncateMarkdown should not append suffix, got %q", result)
	}
}

func TestTruncateMarkdownWithSuffix_DefaultSuffix(t *testing.T) {
	content := "Para1\n\nPara2\n\nPara3"
	result, truncated := TruncateMarkdownWithSuffix(content, 12, "")
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 12 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
	if !strings.HasSuffix(result, "...") {
		t.Fatalf("empty suffix should use default ..., got %q", result)
	}
}

func TestTruncateMarkdownWithSuffix_CustomSuffix(t *testing.T) {
	content := "Hello, world!"
	result, truncated := TruncateMarkdownWithSuffix(content, 10, "…")
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 10 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
	if !strings.HasSuffix(result, "…") {
		t.Fatalf("expected custom suffix, got %q", result)
	}
}

func TestTruncateMarkdownWithSuffix_SuffixExceedsBudget(t *testing.T) {
	result, truncated := TruncateMarkdownWithSuffix("hello world", 2, "...")
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 2 {
		t.Fatalf("result exceeds maxBytes: len=%d result=%q", len(result), result)
	}
}

func TestTruncateMarkdownWithSuffix_OversizedSingleNode(t *testing.T) {
	content := strings.Repeat("A", 100)
	result, truncated := TruncateMarkdownWithSuffix(content, 20, "...")
	if !truncated {
		t.Fatal("expected truncated")
	}
	if len(result) > 20 {
		t.Fatalf("result exceeds maxBytes: len=%d", len(result))
	}
	if !strings.HasSuffix(result, "...") {
		t.Fatalf("expected suffix, got %q", result)
	}
	contentPart := strings.TrimSuffix(result, "...")
	if len(contentPart) != 17 {
		t.Fatalf("expected 17 bytes content, got %d: %q", len(contentPart), contentPart)
	}
}

func TestTruncateUTF8Bytes(t *testing.T) {
	s := "你好世界"
	truncated := truncateUTF8Bytes(s, 7)
	if len([]byte(truncated)) > 7 {
		t.Fatalf("exceeds max bytes: %d", len([]byte(truncated)))
	}
}

func TestFallbackTruncate(t *testing.T) {
	content := strings.Repeat("x", 50)
	result := fallbackTruncate(content, 10)
	if len(result) != 10 {
		t.Fatalf("expected 10 bytes, got %d", len(result))
	}
}
