## 1. 配置与基础设施

- [x] 1.1 在 `libs/einox/backend/git/` 新增 `config.go`：定义 `Config`、`Transport`、`AuthConfig`、`RefConfig` 及配置校验（URL 非空、ref 互斥、Transport 与 Auth 字段匹配）
- [x] 1.2 新增 `auth.go`：实现 HTTP Basic/Token 与 SSH 私钥认证到 `git.CloneOptions.Auth` 的转换，确保凭据不出现在 error 字符串中
- [x] 1.3 新增 `path.go`：实现 `normalizePath`、virtual path ↔ repo path 映射及 baseDir 前缀校验（防 `..` 穿越）
- [x] 1.4 新增 `errors.go`：定义 `errReadOnly` 等包级 sentinel error

## 2. 仓库加载与 tree 访问

- [x] 2.1 扩展 `git.go`：`GitBackend` 结构体持有 config、懒加载 `sync.Once`、缓存的 commit/tree
- [x] 2.2 实现 `initRepo`：bare clone（memory storage + memfs）、按 RefConfig 解析 branch/tag/commit/HEAD
- [x] 2.3 实现 tree 导航 helper：`resolveTreeEntry(repoPath)` 返回 file/tree entry；blob 读取 helper

## 3. 只读 Backend 方法

- [x] 3.1 实现 `LsInfo`：列出 baseDir 映射后目录的直接子项，返回相对 virtual path 的 `FileInfo`
- [x] 3.2 实现 `Read`：读取 blob 内容，支持 1-based `Offset` 与 `Limit` 行切片（语义对齐 InMemoryBackend）
- [x] 3.3 实现 `GrepRaw`：递归 walk 子树、按 Path/Glob/FileType 过滤、正则匹配并返回 `GrepMatch`
- [x] 3.4 实现 `GlobInfo`：递归 walk + `doublestar` 匹配，返回 `FileInfo` 列表
- [x] 3.5 实现 `Write`、`Edit`：返回 `errReadOnly`

## 4. 构造与导出

- [x] 4.1 实现 `NewGitBackend(cfg Config) (*GitBackend, error)` 并在构造阶段做配置校验
- [x] 4.2 保留 `var _ filesystem.Backend = (*GitBackend)(nil)` 编译期断言
- [x] 4.3 为公开类型与函数补充中文 Go doc 注释

## 5. 单元测试

- [x] 5.1 新增 `git_test.go`：测试 helper 创建 memory bare repo fixture（含 baseDir 子树、多文件目录）
- [x] 5.2 测试 `Read`：正常读取、offset/limit、文件不存在
- [x] 5.3 测试 `LsInfo`：根目录与子目录列表
- [x] 5.4 测试 `GrepRaw`：pattern 匹配、Path/Glob 过滤
- [x] 5.5 测试 `GlobInfo`：glob 匹配
- [x] 5.6 测试 baseDir 映射与路径穿越拒绝
- [x] 5.7 测试 branch/tag/commit 三种 ref 解析
- [x] 5.8 测试 `Write`/`Edit` 返回只读错误
- [x] 5.9 运行 `go test ./libs/einox/backend/git/...` 确保全部通过

## 6. 收尾

- [x] 6.1 确认 `go.mod` 无需新增依赖（复用 go-git、eino/doublestar）
- [x] 6.2 确认敏感字段无日志泄露，错误信息不含 token/密码/私钥
