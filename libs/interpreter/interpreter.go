package interpreter

import (
	"context"
	"time"
)

// Interpreter 解释器接口
type LuaInterpreter interface {
    Eval(code string) (any, error)
    EvalAs(code string, datas ... any) error
    EvalWithContext(ctx context.Context, code string) (any, error)
    EvalWithContextAs(ctx context.Context, code string, datas ... any) error
    RegisterTool(name string, fn any) error
    GetGlobal(name string) (any, error)
    SetGlobal(name string, value any) error
    Close() error
}

// Config 配置
type Config struct {
    Timeout     time.Duration
    MemoryLimit int64  // 字节
    MaxSteps    int64  // 最大执行步数（防止无限循环）
}