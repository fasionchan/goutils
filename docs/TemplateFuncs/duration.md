---
tags:
  - 日期时间
---

将时长格式化为本地化可读字符串。入参分两类，不可混用：

- `duration*`：接收 goutils `std/_time.Duration`
- `nativeDuration*`：接收标准库 `time.Duration`

`durationFromNative` 用于把 `time.Duration` 转成 `_time.Duration`。

支持的语言：`zh`（年/月/天/时/分/秒）、`en`（yr/mo/d/hr/min/sec）。
默认拆分到秒，不足 1 秒的部分不会输出。未知 locale 得到空字符串。

## durationFromNative

把标准库 `time.Duration` 转成 goutils `_time.Duration`，以便交给 `durationLocaleString` / `durationLocaleStringPro`。

### 用法

```go
func (d time.Duration) Duration
```

### 示例

以下输出：1时1分

```text
{{ durationLocaleString (durationFromNative .d) "zh" }}
```

其中 `.d` 为 `time.Hour + time.Minute + time.Second`。

## durationLocaleString

把 `_time.Duration` 格式化为本地化字符串。等价于：

```text
durationLocaleStringPro d locale true true 2
```

即去掉高位连续的零单位、去掉低位连续的零单位，最多保留 2 个单位。

### 用法

```go
func (d Duration, locale string) string
```

### 示例

以下输出：1时1分

```text
{{ durationLocaleString .d "zh" }}
```

以下输出：1hr1min

```text
{{ durationLocaleString .d "en" }}
```

其中 `.d` 为 `1h1m1s`（秒被 `limit=2` 截掉）。

## durationLocaleStringPro

`durationLocaleString` 的完整版，可控制是否剔除首尾零单位，以及最多输出几个单位。

### 用法

```go
func (d Duration, locale string, purgeZeroHead, purgeZeroTail bool, limit int) string
```

- `purgeZeroHead`：去掉高位连续的零单位（如 0年0月0天）
- `purgeZeroTail`：去掉低位连续的零单位（如 0分0秒）
- `limit`：最多保留几个单位；中间为零的单位仍占名额

### 示例

保留 3 个单位，并剔除首尾零。以下输出：1时1分1秒

```text
{{ durationLocaleStringPro .d "zh" true true 3 }}
```

其中 `.d` 为 `1h1m1s`。

保留中间的零单位。以下输出：1时0分1秒

```text
{{ durationLocaleStringPro .d "zh" true false 6 }}
```

其中 `.d` 为 `1h1s`。

## nativeDurationLocaleString

把标准库 `time.Duration` 格式化为本地化字符串，语义与 `durationLocaleString` 相同（默认 `purgeZeroHead=true`、`purgeZeroTail=true`、`limit=2`）。

### 用法

```go
func (d time.Duration, locale string) string
```

### 示例

以下输出：1时1分

```text
{{ nativeDurationLocaleString .d "zh" }}
```

以下输出：1hr1min

```text
{{ nativeDurationLocaleString .d "en" }}
```

其中 `.d` 为 `time.Hour + time.Minute + time.Second`。

## nativeDurationLocaleStringPro

`nativeDurationLocaleString` 的完整版，参数含义与 `durationLocaleStringPro` 相同，只是入参为 `time.Duration`。

### 用法

```go
func (d time.Duration, locale string, purgeZeroHead, purgeZeroTail bool, limit int) string
```

### 示例

保留中间的零单位。以下输出：1时0分1秒

```text
{{ nativeDurationLocaleStringPro .d "zh" true false 6 }}
```

其中 `.d` 为 `time.Hour + time.Second`。
