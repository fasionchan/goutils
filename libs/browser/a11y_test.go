package browser

import (
	"strings"
	"testing"

	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/require"
	"github.com/ysmood/gson"
)

func TestAccessibilityAxNodesString(t *testing.T) {
	nodes := AccessibilityAxNodes{
		{
			NodeID:           "1",
			Role:             axRole("RootWebArea"),
			Name:             axName("Example"),
			ChildIDs:         []proto.AccessibilityAXNodeID{"2", "3", "4", "5"},
			BackendDOMNodeID: 1,
		},
		{
			NodeID:           "2",
			ParentID:         "1",
			Role:             axRole("generic"),
			ChildIDs:         []proto.AccessibilityAXNodeID{"21"},
			BackendDOMNodeID: 2,
		},
		{
			NodeID:           "21",
			ParentID:         "2",
			Role:             axRole("button"),
			Name:             axName("Submit"),
			Properties:       []*proto.AccessibilityAXProperty{axBoolProp(proto.AccessibilityAXPropertyNameFocusable, true)},
			BackendDOMNodeID: 21,
		},
		{
			NodeID:           "3",
			ParentID:         "1",
			Role:             axRole("link"),
			Name:             axName("Docs"),
			Properties: []*proto.AccessibilityAXProperty{
				axStringProp(proto.AccessibilityAXPropertyNameURL, "https://example.com/docs"),
				axBoolProp(proto.AccessibilityAXPropertyNameFocusable, true),
			},
			ChildIDs:         []proto.AccessibilityAXNodeID{"31"},
			BackendDOMNodeID: 3,
		},
		{
			NodeID:           "31",
			ParentID:         "3",
			Role:             axInternalRole("StaticText"),
			Name:             axName("Docs"),
			ChildIDs:         []proto.AccessibilityAXNodeID{"32"},
			BackendDOMNodeID: 31,
		},
		{
			NodeID:           "32",
			ParentID:         "31",
			Role:             axInternalRole("InlineTextBox"),
			BackendDOMNodeID: 31,
		},
		{
			NodeID:           "4",
			ParentID:         "1",
			Role:             axRole("textbox"),
			Name:             axName("Query"),
			Value:            axValue("hello"),
			Properties: []*proto.AccessibilityAXProperty{
				axBoolProp(proto.AccessibilityAXPropertyNameFocused, true),
				axTokenProp(proto.AccessibilityAXPropertyNameInvalid, "false"),
			},
			BackendDOMNodeID: 4,
		},
		{
			NodeID:           "5",
			ParentID:         "1",
			Role:             axRole("image"),
			BackendDOMNodeID: 5,
		},
		{
			NodeID:   "6",
			ParentID: "1",
			Ignored:  true,
			Role:     axRole("generic"),
		},
	}

	got := nodes.String()
	require.Equal(t, strings.TrimSpace(`
- RootWebArea "Example" [ref=1]:
  - button "Submit" [ref=21]
  - link "Docs" [ref=3]:
    - /url: https://example.com/docs
  - textbox "Query" [active] [ref=4]:
    - text: hello
`), strings.TrimSpace(got))
	require.NotContains(t, got, "InlineTextBox")
	require.NotContains(t, got, "focusable")
	require.NotContains(t, got, `image`)
	require.NotContains(t, got, `[ref=2]`) // collapsed empty generic
}

func axRole(role string) *proto.AccessibilityAXValue {
	return &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeRole, Value: gson.New(role)}
}

func axInternalRole(role string) *proto.AccessibilityAXValue {
	return &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeInternalRole, Value: gson.New(role)}
}

func axName(name string) *proto.AccessibilityAXValue {
	return &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeComputedString, Value: gson.New(name)}
}

func axValue(value string) *proto.AccessibilityAXValue {
	return &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeString, Value: gson.New(value)}
}

func axBoolProp(name proto.AccessibilityAXPropertyName, value bool) *proto.AccessibilityAXProperty {
	return &proto.AccessibilityAXProperty{
		Name:  name,
		Value: &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeBooleanOrUndefined, Value: gson.New(value)},
	}
}

func axStringProp(name proto.AccessibilityAXPropertyName, value string) *proto.AccessibilityAXProperty {
	return &proto.AccessibilityAXProperty{
		Name:  name,
		Value: &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeString, Value: gson.New(value)},
	}
}

func axTokenProp(name proto.AccessibilityAXPropertyName, value string) *proto.AccessibilityAXProperty {
	return &proto.AccessibilityAXProperty{
		Name:  name,
		Value: &proto.AccessibilityAXValue{Type: proto.AccessibilityAXValueTypeToken, Value: gson.New(value)},
	}
}
