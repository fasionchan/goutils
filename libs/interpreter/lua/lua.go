package lua

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/libs/interpreter"
	"github.com/fasionchan/goutils/stl"
	lua "github.com/yuin/gopher-lua"
)

type luaInterpreter struct {
    L        *lua.LState
    config   interpreter.Config
    tools    map[string]lua.LGFunction
    mu       sync.RWMutex
    closed   bool
}

// NewLuaInterpreter 创建新的 Lua 解释器
func NewLuaInterpreter(config interpreter.Config) (interpreter.LuaInterpreter, error) {
    L := lua.NewState()

    // 应用安全限制
    // if config.MaxSteps > 0 {
    //     L.SetMx(0, config.MaxSteps) // 设置最大执行步数
    // }
    
    // 移除危险函数（安全沙箱）
    RemoveDangerousFunctions(L)
    
    return &luaInterpreter{
        L:      L,
        config: config,
        tools:  make(map[string]lua.LGFunction),
    }, nil
}

// Eval 执行 Lua 代码
func (li *luaInterpreter) Eval(code string) (any, error) {
    li.mu.Lock()
    defer li.mu.Unlock()
    
    if li.closed {
        return nil, fmt.Errorf("interpreter closed")
    }

    luaValues, err := Eval(li.L, code)
    if err != nil {
        return nil, err
    }

    values := stl.Map(luaValues, ConvertToGo)

    if len(values) == 0 {
        return nil, nil
    }

    if len(values) == 1 {
        return values[0], nil
    }

    last := values[len(values)-1]
    if last != nil {
        return nil, fmt.Errorf("lua return error: %v", last)
    }

    values = values[:len(values)-1]
    if len(values) == 1 {
        return values[0], nil
    }

    return values, nil
}

func Eval(L *lua.LState, code string) ([]lua.LValue, error) {
    // 执行代码
    if err := L.DoString(code); err != nil {
        return nil, fmt.Errorf("lua execution error: %w", err)
    }
    
    // 获取返回值（如果有）
    top := L.GetTop()
    if top == 0 {
        return nil, nil
    }

    // 转换返回值
    values := stl.CountAndMap(0, top, 1, func (i int) lua.LValue {
        return L.Get(-top+i)
    })
    
    // 清理栈
    L.Pop(top)

    return values, nil
}

func (li *luaInterpreter) EvalAs(code string, datas ... any) error {
    luaValues, err := Eval(li.L, code)
    if err != nil {
        return err
    }

    if len(datas) != len(luaValues) {
        return fmt.Errorf("lua return error: %v", luaValues[len(luaValues)-1])
    }

    for i, data := range datas {
        if err := ConvertToGoReflect(luaValues[i], reflect.ValueOf(data)); err != nil {
            return err
        }
    }

    return nil
}

// EvalWithContext 带上下文的执行（支持超时）
func (li *luaInterpreter) EvalWithContext(ctx context.Context, code string) (any, error) {
    return nil, baseutils.NewNotImplementedError("luaInterpreter.EvalWithContext")
}

func (li *luaInterpreter) EvalWithContextAs(ctx context.Context, code string, datas ... any) error {
    return baseutils.NewNotImplementedError("luaInterpreter.EvalWithContextAs")
}

func (li *luaInterpreter) RegisterTool(name string, fn any) error {
    fnValue := reflect.ValueOf(fn)
    if fnValue.Kind() != reflect.Func {
        return baseutils.NewBlankBadTypeError().
            WithExpected("function").
            WithGivenReflectType(fnValue.Type())
    }

    fnType := fnValue.Type()

    // 将 Go 函数包装为 Lua 函数
	return li.RegisterLuaTool(name, func(L *lua.LState) int {
        // 获取参数数量
        nargs := L.GetTop()
		if nargs != fnType.NumIn() {
			L.RaiseError("tool %s has %d arguments, but got %d", name, fnType.NumIn(), nargs)
			return 0
		}

        args := make([]reflect.Value, nargs)
        
        // 转换参数
        for i := 0; i < nargs; i++ {
            luaArg := L.Get(i + 1)

            arg := reflect.New(fnType.In(i)).Elem()
            if err := ConvertToGoReflect(luaArg, arg); err != nil {
                L.RaiseError("tool %s argument %d of type %v is not convertible to %v: %v", name, i, luaArg.Type(), fnType.In(i), err)
                return 0
            }

            args[i] = arg
        }

        // 调用 Go 函数
		outs := fnValue.Call(args)
		for _, out := range outs {
			L.Push(ConvertToLua(L, out.Interface()))
		}

		return len(outs)
    })
}

func (li *luaInterpreter) RegisterLuaTool(name string, fn lua.LGFunction) error {
	li.mu.Lock()
	defer li.mu.Unlock()
	
    // 注册到全局
	li.L.SetGlobal(name, li.L.NewFunction(fn))
	li.tools[name] = fn

	return nil
}

// GetGlobal 获取全局变量
func (li *luaInterpreter) GetGlobal(name string) (any, error) {
    li.mu.RLock()
    defer li.mu.RUnlock()
    
    value := li.L.GetGlobal(name)
    if value == lua.LNil {
        return nil, fmt.Errorf("global variable '%s' not found", name)
    }
    
    return ConvertToGo(value), nil
}

// SetGlobal 设置全局变量
func (li *luaInterpreter) SetGlobal(name string, value any) error {
    li.mu.Lock()
    defer li.mu.Unlock()
    
    li.L.SetGlobal(name, ConvertToLua(li.L, value))
    return nil
}

// Close 关闭解释器
func (li *luaInterpreter) Close() error {
    li.mu.Lock()
    defer li.mu.Unlock()
    
    if li.closed {
        return nil
    }
    
    li.closed = true
    li.L.Close()

    return nil
}