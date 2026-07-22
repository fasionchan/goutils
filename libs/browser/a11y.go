package browser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"github.com/go-rod/rod/lib/proto"
)

type AccessibilityAxNodes []*proto.AccessibilityAXNode

func (nodes AccessibilityAxNodes) PurgeNil() AccessibilityAxNodes {
	return stl.PurgeZero(nodes)
}

// String 将无访问性树渲染为面向 AI agent 的 YAML 快照。
// 节点引用直接使用 backendDOMNodeId（[ref=<id>]），不做额外 ref 映射。
func (nodes AccessibilityAxNodes) String() string {
	nodes = nodes.PurgeNil()
	if len(nodes) == 0 {
		return ""
	}

	byID := make(map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode, len(nodes))
	for _, node := range nodes {
		byID[node.NodeID] = node
	}

	var roots AccessibilityAxNodes
	for _, node := range nodes {
		if node.ParentID == "" {
			roots = append(roots, node)
			continue
		}
		if _, ok := byID[node.ParentID]; !ok {
			roots = append(roots, node)
		}
	}
	if len(roots) == 0 {
		roots = AccessibilityAxNodes{nodes[0]}
	}

	var b strings.Builder
	for _, root := range roots {
		renderAccessibilityAxNode(&b, root, byID, 0, "")
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderAccessibilityAxNode(
	b *strings.Builder,
	node *proto.AccessibilityAXNode,
	byID map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode,
	depth int,
	parentName string,
) {
	if node == nil {
		return
	}

	role := accessibilityAxValueString(node.Role)
	name := normalizeAccessibilityText(accessibilityAxValueString(node.Name))
	value := normalizeAccessibilityText(accessibilityAxValueString(node.Value))
	states, props := collectAccessibilityStatesAndProps(node)

	if shouldSkipAccessibilityNode(node, role, name, value, states, props) {
		renderAccessibilityChildren(b, node, byID, depth, parentName)
		return
	}

	if role == "StaticText" || role == "LegacyAccessible" {
		text := name
		if text == "" {
			text = value
		}
		if text == "" || text == parentName || !textContributesBeyondName(parentName, text) {
			return
		}
		writeAccessibilityIndent(b, depth)
		b.WriteString("- text: ")
		b.WriteString(yamlEscapeScalar(text))
		b.WriteByte('\n')
		return
	}

	if shouldCollapseAccessibilityNode(role, name, value, states, props) {
		renderAccessibilityChildren(b, node, byID, depth, parentName)
		return
	}

	writeAccessibilityIndent(b, depth)
	b.WriteString("- ")
	b.WriteString(role)
	if name != "" {
		b.WriteByte(' ')
		b.WriteString(strconv.Quote(truncateAccessibilityText(name, 300)))
	}
	for _, state := range states {
		b.WriteByte(' ')
		b.WriteString(state)
	}
	if node.BackendDOMNodeID != 0 {
		fmt.Fprintf(b, " [ref=%d]", node.BackendDOMNodeID)
	}

	childrenParentName := name
	if childrenParentName == "" {
		childrenParentName = parentName
	}

	var body strings.Builder
	for _, prop := range props {
		writeAccessibilityIndent(&body, depth+1)
		body.WriteString("- /")
		body.WriteString(prop.name)
		body.WriteString(": ")
		body.WriteString(yamlEscapeScalar(prop.value))
		body.WriteByte('\n')
	}
	if value != "" && value != name && isValueBearingRole(role) {
		writeAccessibilityIndent(&body, depth+1)
		body.WriteString("- text: ")
		body.WriteString(yamlEscapeScalar(value))
		body.WriteByte('\n')
	}
	renderAccessibilityChildren(&body, node, byID, depth+1, childrenParentName)

	if body.Len() == 0 {
		b.WriteByte('\n')
		return
	}
	b.WriteString(":\n")
	b.WriteString(body.String())
}

func renderAccessibilityChildren(
	b *strings.Builder,
	node *proto.AccessibilityAXNode,
	byID map[proto.AccessibilityAXNodeID]*proto.AccessibilityAXNode,
	depth int,
	parentName string,
) {
	for _, childID := range node.ChildIDs {
		child, ok := byID[childID]
		if !ok || child == nil {
			continue
		}
		renderAccessibilityAxNode(b, child, byID, depth, parentName)
	}
}

func shouldSkipAccessibilityNode(
	node *proto.AccessibilityAXNode,
	role, name, value string,
	states []string,
	props []accessibilityProp,
) bool {
	if node.Ignored {
		return true
	}
	switch role {
	case "", "none", "presentation", "InlineTextBox":
		return true
	case "image", "img":
		// 无障碍名/描述的装饰图对 agent 无操作价值
		if name == "" && value == "" && len(props) == 0 {
			return true
		}
	}
	return false
}

func shouldCollapseAccessibilityNode(role, name, value string, states []string, props []accessibilityProp) bool {
	switch role {
	case "generic", "group", "LabelText":
	default:
		return false
	}
	if name != "" || value != "" || len(states) > 0 || len(props) > 0 {
		return false
	}
	return true
}

type accessibilityProp struct {
	name  string
	value string
}

func collectAccessibilityStatesAndProps(node *proto.AccessibilityAXNode) (states []string, props []accessibilityProp) {
	for _, property := range node.Properties {
		if property == nil || property.Value == nil {
			continue
		}
		raw := accessibilityAxValueString(property.Value)
		switch property.Name {
		case proto.AccessibilityAXPropertyNameChecked:
			switch raw {
			case "true":
				states = append(states, "[checked]")
			case "mixed":
				states = append(states, "[checked=mixed]")
			}
		case proto.AccessibilityAXPropertyNamePressed:
			switch raw {
			case "true":
				states = append(states, "[pressed]")
			case "mixed":
				states = append(states, "[pressed=mixed]")
			}
		case proto.AccessibilityAXPropertyNameSelected:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[selected]")
			}
		case proto.AccessibilityAXPropertyNameExpanded:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[expanded]")
			}
		case proto.AccessibilityAXPropertyNameDisabled:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[disabled]")
			}
		case proto.AccessibilityAXPropertyNameRequired:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[required]")
			}
		case proto.AccessibilityAXPropertyNameReadonly:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[readonly]")
			}
		case proto.AccessibilityAXPropertyNameFocused:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[active]")
			}
		case proto.AccessibilityAXPropertyNameInvalid:
			if raw != "" && raw != "false" {
				if raw == "true" {
					states = append(states, "[invalid]")
				} else {
					states = append(states, "[invalid="+raw+"]")
				}
			}
		case proto.AccessibilityAXPropertyNameLevel:
			if raw != "" && raw != "0" {
				states = append(states, "[level="+raw+"]")
			}
		case proto.AccessibilityAXPropertyNameMultiline:
			if isTruthyAccessibilityToken(raw) {
				states = append(states, "[multiline]")
			}
		case proto.AccessibilityAXPropertyNameEditable:
			if raw != "" && raw != "false" {
				states = append(states, "[editable]")
			}
		case proto.AccessibilityAXPropertyNameURL:
			if raw != "" {
				props = append(props, accessibilityProp{name: "url", value: raw})
			}
		}
	}
	if desc := normalizeAccessibilityText(accessibilityAxValueString(node.Description)); desc != "" {
		props = append(props, accessibilityProp{name: "description", value: desc})
	}
	return states, props
}

func accessibilityAxValueString(v *proto.AccessibilityAXValue) string {
	if v == nil || v.Value.Nil() {
		return ""
	}
	switch raw := v.Value.Val().(type) {
	case string:
		return raw
	case bool:
		if raw {
			return "true"
		}
		return "false"
	case float64:
		if raw == float64(int64(raw)) {
			return strconv.FormatInt(int64(raw), 10)
		}
		return strconv.FormatFloat(raw, 'f', -1, 64)
	case int:
		return strconv.Itoa(raw)
	case int64:
		return strconv.FormatInt(raw, 10)
	default:
		s := v.Value.Str()
		if s == "<nil>" {
			return ""
		}
		return s
	}
}

func isTruthyAccessibilityToken(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true", "yes", "on", "1":
		return true
	default:
		return false
	}
}

func isValueBearingRole(role string) bool {
	switch role {
	case "textbox", "searchbox", "combobox", "slider", "spinbutton", "progressbar", "scrollbar":
		return true
	default:
		return false
	}
}

func normalizeAccessibilityText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func truncateAccessibilityText(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func textContributesBeyondName(name, text string) bool {
	if name == "" {
		return true
	}
	if text == name {
		return false
	}
	// 子文本几乎被父 name 覆盖时视为冗余
	if strings.Contains(name, text) {
		return false
	}
	return true
}

func writeAccessibilityIndent(b *strings.Builder, depth int) {
	for i := 0; i < depth; i++ {
		b.WriteString("  ")
	}
}

func yamlEscapeScalar(s string) string {
	if s == "" {
		return `""`
	}
	if needsYAMLQuote(s) {
		return strconv.Quote(s)
	}
	return s
}

func needsYAMLQuote(s string) bool {
	if s == "" {
		return true
	}
	switch strings.ToLower(s) {
	case "true", "false", "null", "yes", "no", "on", "off":
		return true
	}
	if strings.ContainsAny(s, "\n\t#{}[]&*!|>'\"%@`") {
		return true
	}
	// 允许 URL 中的 ':'，仅对 ": " 等歧义场景加引号
	if strings.Contains(s, ": ") {
		return true
	}
	if strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
		return true
	}
	if strings.HasPrefix(s, "-") || strings.HasPrefix(s, "?") || strings.HasPrefix(s, ":") {
		return true
	}
	return false
}
