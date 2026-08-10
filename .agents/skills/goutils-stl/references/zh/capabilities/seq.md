# 能力：seq

源码重点：`stl/seq.go`（Go `iter.Seq` / `Seq2` 桥接）。

## 优先符号

| 需求 | 符号 |
|------|------|
| 从值构造 | `DataSeq`、`IndexDataSeq`、`SingularDataSeq` |
| 空 / 拼接 | `EmptySeq`、`EmptySeq2`、`MultiSeq`、`MultiSeq2` |
| 映射 seq | `MapSeq`、`MapSeqToSlice`、`Seq2ToSeq` |
| 物化为切片 | `ReadSeq` |
| 写出 | `WriteSeq`、`WriteSeq2Key`、`WriteSeq2Value`、`WriteSeq2DataError` |
| 分组助手 | `Seqs`、`Seq2s`、`NewSeqs`、`NewSeq2s` |

面向 Go 1.23+ 迭代器；否则优先用 slice 助手。
