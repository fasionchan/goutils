package datarender

import (
	"fmt"
	"strings"

	"github.com/fasionchan/goutils/libs/interpreter"
	"github.com/fasionchan/goutils/libs/interpreter/lua"
)

// LuaTextRender 将 templateData 注入 Lua 全局变量 data，执行脚本后得到文本再由 dataParser 解析。
type LuaTextRender[Data any] struct {
	code       string
	dataParser func(text string) (Data, error)
}

func ParseLuaTextRenderTemplate[Data any](text string, dataParser func(text string) (Data, error)) (DataRender[Data], error) {
	if !strings.HasPrefix(text, TemplateTypeMagicLua) {
		return nil, fmt.Errorf("lua template must start with %q", TemplateTypeMagicLua)
	}

	return &LuaTextRender[Data]{
		code:       strings.TrimPrefix(text, TemplateTypeMagicLua),
		dataParser: dataParser,
	}, nil
}

func (render *LuaTextRender[Data]) Render(data any) (Data, error) {
	var zero Data
	if render == nil {
		return zero, fmt.Errorf("LuaTextRender is nil")
	}

	interpreter_, err := lua.NewLuaInterpreter(interpreter.Config{})
	if err != nil {
		return zero, fmt.Errorf("NewLuaInterpreterFailed: %w", err)
	}
	defer interpreter_.Close()

	if err := interpreter_.SetGlobal("data", data); err != nil {
		return zero, fmt.Errorf("LuaSetGlobalDataFailed: %w", err)
	}

	var result, errmsg string
	if err := interpreter_.EvalAs(render.code,&result,&errmsg); err != nil {
		return zero, fmt.Errorf("LuaEvalFailed: %w", err)
	}

	if errmsg != "" {
		return zero, fmt.Errorf("LuaEvalFailed: %s", errmsg)
	}

	if render.dataParser == nil {
		return zero, fmt.Errorf("LuaTextRender dataParser is nil")
	}

	return render.dataParser(result)
}