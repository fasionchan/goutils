# 能力：buffer / io

源码重点：`stl/buffer.go`、`stl/buffer_pro.go`、`stl/io.go`、`stl/writer.go`、`stl/reader.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 泛型缓冲 | `NewBuffer`、`NewBufferFrom`、`Buffer.Write` / `Read` / `Datas` / `Reset` |
| 有界缓冲 | `NewBoundedBuffer` |
| 关闭助手 | `Close`、`CloseQuietly`、`Closers`、`NopCloser` |
| Writer 适配 | `Writer`、`NewSkipWriter`、`NewLimitWriter`、`NewPrinter` |

`Buffer` 语义对齐 `bytes.Buffer`，但元素类型可为任意 `Datas ~[]Data`，不限于 `[]byte`。
