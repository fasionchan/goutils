package lua

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/fasionchan/goutils/libs/interpreter"
	"github.com/fasionchan/goutils/types"
	"github.com/stretchr/testify/assert"
)

func TestNewLuaInterpreter(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	interpreter.RegisterTool("negative", func(value float64) (float64) {
		return -value
	})

	result, err := interpreter.Eval(`
		return negative(os.time()), nil
	`)
	if err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, float64(-time.Now().Unix()), result)

	t.Logf("Result: %v", result)
}

func TestNewLuaTypes(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	for _, test := range []struct {
		code string
		expected any
	}{
		{`return nil`, nil},
		{`return 1`, 1.0},
		{`return 1.0`, 1.0},
		{`return true`, true},
		{`return false`, false},
		{`return "hello"`, "hello"},
		{`return {1, 2, 3}`, []any{1.0, 2.0, 3.0}},
		{`return {a = 1, b = 2, c = 3}`, map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}},
		{`return {1, 2, 3, a = 4, b = 5, c = 6}`, map[string]any{"1": 1.0, "2": 2.0, "3": 3.0, "a": 4.0, "b": 5.0, "c": 6.0}},
	} {
		result, err := interpreter.Eval(test.code)
		if err != nil {
			t.Fatalf("Failed to evaluate Lua code: %v", err)
		}
		assert.Equal(t, test.expected, result, test.code)
	}
}

func TestJson(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	interpreter.RegisterTool("json", Json)

	result, err := interpreter.Eval(`return json({a = 1, b = 2, c = 3})`)
	if err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, `{"a":1,"b":2,"c":3}`, result)
}

func TestUnjson(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	interpreter.RegisterTool("unjson", Unjson)

	result, err := interpreter.Eval(`return unjson('{"a":1,"b":2,"c":3}')`)
	if err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, map[string]any{"a": 1.0, "b": 2.0, "c": 3.0}, result)
}

func TestJsonUnmarshal(t *testing.T) {
	m := map[string]string{
		"a": "a",
		"b": "b",
	}
	
	if err := json.Unmarshal([]byte(`{"b": "B", "c": "C"}`), &m); err != nil {
		t.Fatalf("Failed to unmarshal map: %v", err)
	}

	assert.Equal(t, map[string]string{"a": "a", "b": "B", "c": "C"}, m)

	if err := json.Unmarshal([]byte(`null`), &m); err != nil {
		t.Fatalf("Failed to unmarshal map: %v", err)
	}

	assert.Equal(t, map[string]string(nil), m)

	s := []int{
		1, 2, 3,
	}

	if err := json.Unmarshal([]byte(`[4, 5]`), &s); err != nil {
		t.Fatalf("Failed to unmarshal slice: %v", err)
	}

	assert.Equal(t, []int{4, 5}, s)

	if err := json.Unmarshal([]byte(`null`), &s); err != nil {
		t.Fatalf("Failed to unmarshal slice: %v", err)
	}

	assert.Equal(t, []int(nil), s)

	a := [3]int{
		1, 2, 3,
	}

	if err := json.Unmarshal([]byte(`[4, 5]`), &a); err != nil {
		t.Fatalf("Failed to unmarshal array: %v", err)
	}

	assert.Equal(t, [3]int{4, 5, 0}, a)

	if err := json.Unmarshal([]byte(`null`), &a); err != nil {
		t.Fatalf("Failed to unmarshal int: %v", err)
	}

	assert.Equal(t, [3]int{4, 5, 0}, a)

	if err := json.Unmarshal([]byte(`[4, 5, 6, 7]`), &a); err != nil {
		t.Fatalf("Failed to unmarshal array: %v", err)
	}

	assert.Equal(t, [3]int{4, 5, 6}, a)

	st := struct {
		Name string `json:"name"`
		Age int `json:"age"`
	} {
		Name: "John",
		Age: 30,
	}

	if err := json.Unmarshal([]byte(`{"name": "John", "age": 40}`), &st); err != nil {
		t.Fatalf("Failed to unmarshal struct: %v", err)
	}

	assert.Equal(t, st.Name, "John")
	assert.Equal(t, st.Age, 40)

	if err := json.Unmarshal([]byte(`null`), &st); err != nil {
		t.Fatalf("Failed to unmarshal struct: %v", err)
	}

	assert.Equal(t, st.Name, "John")
	assert.Equal(t, st.Age, 40)

	str := "hello"
	if err := json.Unmarshal([]byte(`"world"`), &str); err != nil {
		t.Fatalf("Failed to unmarshal string: %v", err)
	}

	assert.Equal(t, str, "world")

	if err := json.Unmarshal([]byte(`null`), &str); err != nil {
		t.Fatalf("Failed to unmarshal string: %v", err)
	}

	assert.Equal(t, "world", str)
}

func TestEvalAsHttpHeaders(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	var headers http.Header
	if err := interpreter.EvalAs(`
		return {
			["Content-Type"] = {"application/json"},
		}
	`, &headers); err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, http.Header{"Content-Type": {"application/json"}}, headers)
}

func TestEvalAsHttpHeaders2(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	var headers http.Header
	if err := interpreter.EvalAs(`
		return {
			["Foo"] = {"foo", "Foo"},
			["Bar"] = {"bar"},
			["Baz"] = {},
			["Qux"] = nil,
		}
	`, &headers); err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, http.Header{"Foo": {"foo", "Foo"}, "Bar": {"bar"}, "Baz": {}}, headers)
}

func TestEvalAsStringMap(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()
	
	var strs types.Strings
	if err := interpreter.EvalAs(`
		return {"foo", "bar", "baz"}
	`, &strs); err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, types.Strings{"foo", "bar", "baz"}, strs)
}

func TestEvalAsStruct(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()
	
	
	var st struct {
		Name string
		Age int
	}
	if err := interpreter.EvalAs(`
		return {Name = "John", Age = 30}
	`, &st); err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	assert.Equal(t, st.Name, "John")
	assert.Equal(t, st.Age, 30)
}

func TestCompile(t *testing.T) {
}