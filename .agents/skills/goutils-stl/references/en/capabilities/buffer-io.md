# Capability: buffer / io

Source focus: `stl/buffer.go`, `stl/buffer_pro.go`, `stl/io.go`, `stl/writer.go`, `stl/reader.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Generic buffer | `NewBuffer`, `NewBufferFrom`, `Buffer.Write` / `Read` / `Datas` / `Reset` |
| Bounded buffer | `NewBoundedBuffer` |
| Close helpers | `Close`, `CloseQuietly`, `Closers`, `NopCloser` |
| Writer adapters | `Writer`, `NewSkipWriter`, `NewLimitWriter`, `NewPrinter` |

`Buffer` mirrors `bytes.Buffer` semantics for arbitrary `Datas ~[]Data`, not only `[]byte`.
