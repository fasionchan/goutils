# Capability: seq

Source focus: `stl/seq.go` (Go `iter.Seq` / `Seq2` bridges).

## Prefer these symbols

| Need | Symbol |
|------|--------|
| From values | `DataSeq`, `IndexDataSeq`, `SingularDataSeq` |
| Empty / multi | `EmptySeq`, `EmptySeq2`, `MultiSeq`, `MultiSeq2` |
| Map seq | `MapSeq`, `MapSeqToSlice`, `Seq2ToSeq` |
| Materialize | `ReadSeq` |
| Write out | `WriteSeq`, `WriteSeq2Key`, `WriteSeq2Value`, `WriteSeq2DataError` |
| Grouping helpers | `Seqs`, `Seq2s`, `NewSeqs`, `NewSeq2s` |

Use when working with Go 1.23+ iterators; otherwise prefer slice helpers.
