# 预期命中要点（对照 `flawed-snippet.go`）

人工走查时，按 Skill 输出模板至少应标出以下 **≥2 类**问题（严重度供参考，可微调）。

## 1. 命名与风格（naming-style）

- `get_user_info` 使用 `snake_case`，不符合 Go 惯用 `mixedCaps`。
- 参数 `userId` / `UserID` 风格不一致；导出函数参数名 `UserID` 的大小写口味问题可标 `nit`/`minor`。

## 2. 错误处理（error-handling）

- `os.Open` / `Read` 的 `error` 被直接丢弃（`_`）。
- 可预期的「无数据」用 `panic` 而非返回 `error`。

## 3. （可选）轻量安全（security-basics）

- `"/tmp/users/" + userId` 对外部输入做路径拼接，存在目录穿越风险（可标 `major`）。

## 对照结论

能稳定命中「命名」+「错误处理」两类即满足 AC-12；安全项为加分命中。
