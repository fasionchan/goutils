# Capability: chan

Source focus: `stl/chan.go`.

## Prefer these symbols

| Need | Symbol |
|------|--------|
| Typed alias | `Chan`, `Push`, `Pull`, `PushPro` |
| Bulk fill | `PushDataToChanX`, `NewChanFromDatasX` |
| Cancelable pipe | `ChanPipe`, `NewChanPipe`, `NewBufferedChanPipe` |
| Pipe IO | `Push` / `PushWithCtx`, `Pull` / `PullWithCtx`, `Cancel`, `Close` |

Fewer high-level helpers than slice/map — focus on timeout and cancellation.
