# stl 概览（zh）

## 职责

`stl` 是以泛型为主的工具包，覆盖集合、轻量 IO 缓冲、缓存拉取、图与通道。
需要可组合的小助手、又不想在各处复制循环时优先使用。

导入：

```go
import "github.com/fasionchan/goutils/stl"
```

人类文档：https://pkg.go.dev/github.com/fasionchan/goutils/stl

## 何时使用

- 切片流水线：filter / map / reduce / 分块 / 去重
- 从切片构建或变换 `map[K]V`
- 集合运算与成员判断
- 在 `iter.Seq` / `iter.Seq2` 与切片或 Writer 之间桥接
- 元素不只是 `byte` 时用泛型 `Buffer[Datas, Data]`
- 带过期/刷新的缓存拉取、Kahn 拓扑排序、带超时的 channel 推送

## 何时不用

- 领域持久化、HTTP、SQL —— 看 goutils 其他包
- 标准库（`slices`、`maps`、`bytes.Buffer`）已完全匹配时不必强行替换
- 不要「以防万一」全量拉取 `*Pro` / `*Unary` / `*Binary` 变体

## 命名规律

| 后缀 | 含义 |
|------|------|
| （无） | 默认、最简 API |
| `Pro` | 需要额外上下文（下标、整片/整表、更丰富回调） |
| `Unary` / `Binary` / `Ternary` | 向回调绑定额外固定参数 |
| `X` | 可变参数便利形式 |

先选无后缀符号；只有回调真需要更多上下文时再升级。
