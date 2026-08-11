package lua

import (
	"fmt"
	"reflect"

	"github.com/fasionchan/goutils/stl"
	lua "github.com/yuin/gopher-lua"
)

var (
    // todo: move to goutils
    ReflectTypeAny = reflect.TypeOf((*any)(nil)).Elem()
)

func SubLTableByKeys(L *lua.LState, tb *lua.LTable, keys []string) *lua.LTable {
	result := L.NewTable()
	for _, key := range keys {
		if value := tb.RawGetString(key); value != lua.LNil {
			result.RawSetString(key, value)
		}
	}
	return result
}

func SubLTableByKeysX(L *lua.LState, tb *lua.LTable, keys ...string) *lua.LTable {
	return SubLTableByKeys(L, tb, keys)
}

// RemoveDangerousFunctions 移除危险函数，创建安全沙箱
func RemoveDangerousFunctions(L *lua.LState) {
    // 移除 os 模块的危险函数
    os, ok := L.GetGlobal("os").(*lua.LTable)
	if ok {
		L.SetGlobal("os", SubLTableByKeysX(L, os, "clock", "time"))
	}

    // 移除 io 模块（或限制）
    L.SetGlobal("io", lua.LNil)
    
    // 移除 package 模块（禁止 require）
    L.SetGlobal("package", lua.LNil)
    
    // 或者创建自定义的安全环境
    // env := L.NewTable()
    // L.SetEnv(env)
}

// convertToGo 将 Lua 值转换为 Go 值
func ConvertToGo(lv lua.LValue) any {
    switch v := lv.(type) {
    case *lua.LNilType:
        return nil
    case lua.LBool:
        return bool(v)
    case lua.LString:
        return string(v)
    case lua.LNumber:
        return float64(v)
    case *lua.LTable:
        // 转换为 map 或 slice
        return LTableToGoAsAny(v)
    case *lua.LFunction:
        return "[Lua Function]"
    case *lua.LUserData:
        return "[User Data]"
    default:
        return fmt.Sprintf("%v", v)
    }
}

func ConvertToGoTyped(lv lua.LValue, data any) error {
    return ConvertToGoReflect(lv, reflect.ValueOf(data))
}

func ConvertToGoReflect(lv lua.LValue, value reflect.Value) error {
    kind := value.Kind()

    // 处理 nil
    _, ok := lv.(*lua.LNilType)
    if ok {
        switch kind {
        case reflect.Ptr:
            for {
                if value.IsNil() {
                    return nil
                }

                elem := value.Elem()
                if elem.Kind() != reflect.Ptr {
                    break
                }

                value = elem
            }

            if !value.CanSet() {
                return fmt.Errorf("value is not settable: %s", value.Type().String())
            }

            value.SetZero()
        // 需要设置 nil 的类型
        case reflect.Map, reflect.Slice, reflect.Interface:
            if value.IsNil() {
                return nil
            }

            if !value.CanSet() {
                return fmt.Errorf("value is not settable: %s", value.Type().String())
            }

            value.SetZero()
        }

        return nil
    }

    if kind == reflect.Ptr {
        if value.IsNil() {
            if !value.CanSet() {
                return fmt.Errorf("value is not settable: %s", value.Type().String())
            }

            value.Set(reflect.New(value.Type().Elem()))
        }

        return ConvertToGoReflect(lv, value.Elem())
    }

    switch kind {
    case reflect.Map:
        table, ok := lv.(*lua.LTable)
        if !ok {
            return fmt.Errorf("go %s expect lua table, but got %s", value.Type().String(), lv.Type())
        }

        if value.IsNil() {
            if !value.CanSet() {
                return fmt.Errorf("value is not settable: %s", value.Type().String())
            }

            value.Set(reflect.MakeMap(value.Type()))
        }

        var err error
        table.ForEach(func(lk, lv lua.LValue) {
            if err != nil {
                return
            }

            key := reflect.New(value.Type().Key()).Elem()
            if err = ConvertToGoReflect(lk, key); err != nil {
                return
            }

            elem := reflect.New(value.Type().Elem()).Elem()
            if err = ConvertToGoReflect(lv, elem); err != nil {
                return
            }

            value.SetMapIndex(key, elem)
        })

        return err
    case reflect.Slice:
        table, ok := lv.(*lua.LTable)
        if !ok {
            return fmt.Errorf("go %s expect lua table, but got %s", value.Type().String(), lv.Type())
        }

        if !value.CanSet() {
            return fmt.Errorf("value is not settable: %s", value.Type().String())
        }

        value.Set(reflect.MakeSlice(value.Type(), table.Len(), table.Len()))
        for i := 0; i < table.Len(); i++ {
            if err := ConvertToGoReflect(table.RawGetInt(i+1), value.Index(i)); err != nil {
                return err
            }
        }

        return nil
    case reflect.Array:
        table, ok := lv.(*lua.LTable)
        if !ok {
            return fmt.Errorf("go %s expect lua table, but got %s", value.Type().String(), lv.Type())
        }

        n := min(value.Type().Len(), table.Len())
        for i := range n {
            if err := ConvertToGoReflect(table.RawGetInt(i+1), value.Index(i)); err != nil {
                return err
            }
        }

        return nil
    case reflect.Struct:
        table, ok := lv.(*lua.LTable)
        if !ok {
            return fmt.Errorf("go %s expect lua table, but got %s", value.Type().String(), lv.Type())
        }

        var err error
        table.ForEach(func(lk, lv lua.LValue) {
            if err != nil {
                return
            }

            lf, ok := lk.(lua.LString)
            if !ok {
                err = fmt.Errorf("go %s field expect lua string, but got %s", value.Type().String(), lk.Type())
                return
            }

            field := value.FieldByName(string(lf))
            if err = ConvertToGoReflect(lv, field); err != nil {
                return
            }
        })

        return nil
    case reflect.Interface:
        if value.Type() != ReflectTypeAny {
            return fmt.Errorf("go %s expect any, but got %s", value.Type().String(), lv.Type())
        }

        if !value.CanSet() {
            return fmt.Errorf("value is not settable: %s", value.Type().String())
        }

        value.Set(reflect.ValueOf(ConvertToGo(lv)))

        return nil
    case reflect.String:
        s, ok := lv.(lua.LString)
        if !ok {
            return fmt.Errorf("value is not a string: %s", lv.Type().String())
        }

        if !value.CanSet() {
            return fmt.Errorf("value is not settable: %s", value.Type().String())
        }

        value.SetString(string(s))
        return nil
    case reflect.Bool:
        b, ok := lv.(lua.LBool)
        if !ok {
            return fmt.Errorf("value is not a boolean: %s", lv.Type().String())
        }

        if !value.CanSet() {
            return fmt.Errorf("value is not settable: %s", value.Type().String())
        }

        value.SetBool(bool(b))
        return nil
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        i, ok := lv.(lua.LNumber)
        if !ok {
            return fmt.Errorf("value is not a number: %s", lv.Type().String())
        }

        if !value.CanSet() {
            return fmt.Errorf("value is not settable: %s", value.Type().String())
        }

        value.SetInt(int64(i))
        return nil
    case reflect.Float32, reflect.Float64:
        f, ok := lv.(lua.LNumber)
        if !ok {
            return fmt.Errorf("value is not a number: %s", lv.Type().String())
        }

        if !value.CanSet() {
            return fmt.Errorf("value is not settable: %s", value.Type().String())
        }

        value.SetFloat(float64(f))
        return nil
    }

    return fmt.Errorf("unsupported type: %s", value.Type().String())
}

// convertToLua 将 Go 值转换为 Lua 值
func ConvertToLua(L *lua.LState, value any) lua.LValue {
    if value == nil {
        return lua.LNil
    }

    return ConvertReflectValueToLua(L, reflect.ValueOf(value))
}

func ConvertReflectValueToLua(L *lua.LState, value reflect.Value) lua.LValue {
	switch value.Kind() {
	case reflect.Interface:
		if value.IsNil() {
			return lua.LNil
		}
		return ConvertReflectValueToLua(L, value.Elem())
	case reflect.Bool:
        return lua.LBool(value.Bool())
    case reflect.String:
        return lua.LString(value.String())
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return lua.LNumber(value.Int())
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return lua.LNumber(value.Uint())
    case reflect.Float32, reflect.Float64:
        return lua.LNumber(value.Float())
    case reflect.Map:
        table := L.NewTable()
        for _, key := range value.MapKeys() {
            table.RawSetString(key.String(), ConvertReflectValueToLua(L, value.MapIndex(key)))
        }
        return table
    case reflect.Slice:
        if value.Type() == stl.ReflectType[[]byte]() {
            return lua.LString(value.Bytes())
        }

        table := L.NewTable()
        for i := 0; i < value.Len(); i++ {
            table.RawSetInt(i+1, ConvertReflectValueToLua(L, value.Index(i)))
        }
        return table
    case reflect.Array:
        table := L.NewTable()
        for i := 0; i < value.Len(); i++ {
            table.RawSetInt(i+1, ConvertReflectValueToLua(L, value.Index(i)))
        }
        return table
    case reflect.Struct:
        table := L.NewTable()
        for i := 0; i < value.NumField(); i++ {
            table.RawSetString(value.Type().Field(i).Name, ConvertReflectValueToLua(L, value.Field(i)))
        }
        return table
    default:
        // 尝试转换为字符串
        return lua.LString(fmt.Sprintf("%T: %v", value.Interface(), value.Interface()))
    }
}

type LTableParser struct {
    m map[string]any
    a []any

    // feedStrKey func(key string, value lua.LValue)
    // feedNumKey func(index int, value lua.LValue)
}

func NewLTableParser() *LTableParser {
    return &LTableParser{
        m: make(map[string]any),
        a: make([]any, 0),
    }
}

func (p *LTableParser) Feed(key, value lua.LValue) {
    if strKey, ok := key.(lua.LString); ok {
        p.FeedStrKey(string(strKey), value)
    } else if numKey, ok := key.(lua.LNumber); ok {
        p.FeedNumKey(int(numKey), value)
    }
}

func (p *LTableParser) FeedStrKey(key string, value lua.LValue) {
    p.feedStrKey(key, value)

    if a := p.a; a != nil {
        for i, v := range p.a {
            p.m[fmt.Sprintf("%d", i+1)] = v
        }

        p.a = nil
    }
}

func (p *LTableParser) feedStrKey(key string, value lua.LValue) {
    p.m[key] = ConvertToGo(value)
}

func (p *LTableParser) FeedNumKey(index int, value lua.LValue) {
    if p.a == nil {
        p.feedStrKey(fmt.Sprintf("%d", index), value)
        return
    }

    if index == len(p.a) + 1 {
        p.a = append(p.a, ConvertToGo(value))
    } else {
        p.FeedStrKey(fmt.Sprintf("%d", index), value)
    }
}

func (p *LTableParser) GetResult() (map[string]any, []any) {
    if p == nil {
        return nil, nil
    }

    if p.a == nil {
        return p.m, nil
    }

    return nil, p.a
}

func (p *LTableParser) GetResultAsAny() any {
    if p == nil {
        return nil
    }

    if p.a == nil {
        return p.m
    }

    return p.a
}

func LTableToGo(tb *lua.LTable) (map[string]any, []any) {
    parser := NewLTableParser()

    tb.ForEach(parser.Feed)

    return parser.GetResult()
}

func LTableToGoAsAny(tb *lua.LTable) any {
    parser := NewLTableParser()

    tb.ForEach(parser.Feed)

    return parser.GetResultAsAny()
}