/*
 * Author: fasion
 * Created time: 2026-06-14
 */

package mdutils

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const defaultTruncateSuffix = "..."

// TruncateMarkdown 将 Markdown 内容截断到 maxBytes 字节以内，在块级 AST 节点边界处截断以保持语法结构完整。
// 仅执行截断，不追加任何后缀；maxBytes 全部用于内容预算。
// 返回值第二个 bool 表示是否发生了截断；未截断时原样返回内容。
func TruncateMarkdown(content string, maxBytes int) (string, bool) {
	if maxBytes <= 0 {
		if len(content) == 0 {
			return "", false
		}
		return "", true
	}

	if len(content) <= maxBytes {
		return content, false
	}

	truncatedContent, ok := truncateMarkdownByAST(content, maxBytes)
	if !ok || (truncatedContent == "" && len(content) > 0 && maxBytes > 0) {
		truncatedContent = fallbackTruncate(content, maxBytes)
	}

	return truncatedContent, true
}

// TruncateMarkdownWithSuffix 在 TruncateMarkdown 基础上追加省略后缀。
// suffix 为空时使用默认后缀 ...；后缀字节长度计入 maxBytes 总预算。
func TruncateMarkdownWithSuffix(content string, maxBytes int, suffix string) (string, bool) {
	if suffix == "" {
		suffix = defaultTruncateSuffix
	}

	if maxBytes <= 0 {
		if len(content) == 0 {
			return "", false
		}
		return "", true
	}

	if len(content) <= maxBytes {
		return content, false
	}

	suffixLen := len(suffix)
	contentBudget := maxBytes - suffixLen
	if contentBudget < 0 {
		contentBudget = 0
	}

	truncatedContent, truncated := TruncateMarkdown(content, contentBudget)
	if !truncated {
		return truncatedContent, false
	}

	if contentBudget == 0 {
		return truncateUTF8Bytes(suffix, maxBytes), true
	}

	return truncatedContent + suffix, true
}

// truncateMarkdownByAST 基于 goldmark AST 在块级节点边界累加内容。
// 第二个返回值 false 表示 AST 遍历失败，调用方应降级处理。
func truncateMarkdownByAST(content string, contentBudget int) (string, bool) {
	if contentBudget <= 0 {
		return "", true
	}

	md := goldmark.New()
	source := []byte(content)
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var result bytes.Buffer
	done := false

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || done {
			return ast.WalkContinue, nil
		}

		if n.Type() == ast.TypeInline {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindDocument {
			return ast.WalkContinue, nil
		}

		nodeContent, ok := nodeSourceContent(n, source)
		if !ok {
			return ast.WalkContinue, nil
		}
		nodeSize := len(nodeContent)
		currentSize := result.Len()

		// 累加此节点会超出预算，且已有内容：停止截断
		if currentSize+nodeSize > contentBudget && currentSize > 0 {
			done = true
			return ast.WalkStop, nil
		}

		// 单个节点本身超出预算
		if nodeSize > contentBudget {
			if currentSize > 0 {
				done = true
				return ast.WalkStop, nil
			}
			result.WriteString(truncateUTF8Bytes(string(nodeContent), contentBudget))
			done = true
			return ast.WalkSkipChildren, nil
		}

		result.Write(nodeContent)
		if n.Kind() == ast.KindListItem {
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})

	if err != nil {
		return "", false
	}

	return result.String(), true
}

// fallbackTruncate 在 AST 解析失败时按 UTF-8 安全边界做简单字节截断。
func fallbackTruncate(content string, contentBudget int) string {
	if contentBudget <= 0 {
		return ""
	}
	return truncateUTF8Bytes(content, contentBudget)
}

// nodeSourceContent 从 AST 节点提取对应源码片段。
func nodeSourceContent(n ast.Node, source []byte) ([]byte, bool) {
	switch n.Kind() {
	case ast.KindDocument, ast.KindList:
		return nil, false
	case ast.KindTextBlock:
		if parent := n.Parent(); parent != nil && parent.Kind() == ast.KindListItem {
			return nil, false
		}
	case ast.KindListItem:
		return listItemSourceContent(n, source)
	}

	lines := n.Lines()
	if lines.Len() == 0 {
		return nil, false
	}

	content := lines.Value(source)
	if len(content) == 0 {
		return nil, false
	}

	return content, true
}

// listItemSourceContent 提取列表项完整源码（含 marker 与换行）。
func listItemSourceContent(item ast.Node, source []byte) ([]byte, bool) {
	start, end := -1, -1

	err := ast.Walk(item, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Type() == ast.TypeInline {
			return ast.WalkContinue, nil
		}

		lines := n.Lines()
		if lines.Len() == 0 {
			return ast.WalkContinue, nil
		}

		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			if start < 0 || seg.Start < start {
				start = seg.Start
			}
			if seg.Stop > end {
				end = seg.Stop
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil || start < 0 || end <= start {
		return nil, false
	}

	lineStart := start
	for lineStart > 0 && source[lineStart-1] != '\n' {
		lineStart--
	}

	lineEnd := end
	for lineEnd < len(source) && source[lineEnd] != '\n' {
		lineEnd++
	}
	if lineEnd < len(source) && source[lineEnd] == '\n' {
		lineEnd++
	}

	return source[lineStart:lineEnd], true
}

// truncateUTF8Bytes 在 UTF-8 字符边界截断字符串，避免截断多字节字符中间。
func truncateUTF8Bytes(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}

	b := []byte(s)
	if len(b) <= maxBytes {
		return s
	}

	n := maxBytes
	for n > 0 && (b[n]&0xC0) == 0x80 {
		n--
	}
	return string(b[:n])
}
