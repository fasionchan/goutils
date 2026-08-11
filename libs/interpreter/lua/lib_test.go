package lua

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/fasionchan/goutils/libs/interpreter"
	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	interpreter, err := NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		t.Fatalf("Failed to create Lua interpreter: %v", err)
	}
	defer interpreter.Close()

	interpreter.RegisterTool("http", Http)
	interpreter.RegisterTool("json", Json)

	var statusCode int
	var respHeaders http.Header
	var respBody string

	err = interpreter.EvalAs(`
		local headers = {
			["Content-Type"] = {"application/json"},
		}
		local body = json({
			message = "hello world",
		})
		local statusCode, respHeaders, respBody = http("https://httpbin.org/post", "POST", nil, nil, nil, body)
		return statusCode, respHeaders, respBody
	`, &statusCode, &respHeaders, &respBody)
	if err != nil {
		t.Fatalf("Failed to evaluate Lua code: %v", err)
	}

	fmt.Println("statusCode", statusCode)
	fmt.Println("respHeaders", respHeaders)
	fmt.Println("respBody", respBody)

	assert.Equal(t, 200, statusCode)
}