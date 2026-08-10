# 能力：chan

源码重点：`stl/chan.go`。

## 优先符号

| 需求 | 符号 |
|------|------|
| 类型别名 | `Chan`、`Push`、`Pull`、`PushPro` |
| 批量填充 | `PushDataToChanX`、`NewChanFromDatasX` |
| 可取消管道 | `ChanPipe`、`NewChanPipe`、`NewBufferedChanPipe` |
| 管道读写 | `Push` / `PushWithCtx`、`Pull` / `PullWithCtx`、`Cancel`、`Close` |

高层助手少于 slice/map——重点在超时与取消。
