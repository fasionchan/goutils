## Context

`libs/einox/backend/git/git.go` 已定义空的 `GitBackend` 并实现 `filesystem.Backend` 接口的类型断言，但所有方法均未实现。Eino ADK 的 `filesystem.Backend` 定义了 `LsInfo`、`Read`、`GrepRaw`、`GlobInfo`、`Write`、`Edit` 等统一文件操作协议；参考实现 `InMemoryBackend` 提供了路径规范化、行级 Read、正则 Grep、Glob 匹配的完整语义。

项目已有 `go-git/v6` 依赖（`libs/gitx` 中已有 clone 示例），适合用于 bare clone 后在内存中访问指定 commit 的 tree，无需本地工作区 checkout。

## Goals / Non-Goals

**Goals:**

- 实现 `GitBackend` 对 `filesystem.Backend` 的只读方法，语义与 `InMemoryBackend` 对齐
- 支持 HTTP(S) / SSH 两种传输协议及结构化认证配置
- 支持 branch / tag / commit 三种 ref 定位方式
- 支持 `BaseDir` 虚拟根路径映射，防止路径穿越
- 使用本地 bare repo fixture 编写单元测试，不依赖外部网络

**Non-Goals:**

- 不实现 `Write`、`Edit` 及任何 Shell/Execute 能力
- 不实现 `MultiModalReader`（PDF/图片多模态读取）
- 不做增量 fetch、多 ref 切换、仓库缓存共享等高级特性
- 不在本变更中集成到具体 Agent 或 opsys 调用方

## Decisions

### 1. Git 引擎：复用 go-git bare clone

**选择**: 使用 `github.com/go-git/go-git/v6`，以 `Bare: true, NoCheckout: true` 克隆到 `memory.Storage` + `memfs`（或等效内存 storer），解析 ref 后通过 `object.Tree` 遍历文件。

**理由**: 项目已有依赖与 `gitx` 使用先例；纯 Go 实现便于测试与跨平台；bare + NoCheckout 避免写磁盘。

**备选**:
- 调用系统 `git` CLI —— 依赖外部环境，难以在测试中控制
- `git2go` —— CGO 依赖，增加构建复杂度

### 2. 配置结构：分层 Config + Auth

**选择**:

```go
type Config struct {
    URL       string        // 仓库 URL（HTTPS 或 SSH）
    Transport Transport     // HTTP | SSH
    Auth      AuthConfig    // 按 Transport 选用子字段
    Ref       RefConfig     // Branch | Tag | Commit（互斥）
    BaseDir   string        // 仓库内基路径，如 "users/tom"
}

type AuthConfig struct {
    // HTTP
    Username string
    Password string
    Token    string // 优先于 Username/Password

    // SSH
    PrivateKey  []byte
    Passphrase  string
    KnownHosts  []byte // 可选，为空时使用 InsecureIgnoreHostKey（仅测试）或系统 known_hosts
}
```

**理由**: 认证与 ref、路径配置分离，调用方可从环境变量/密钥管理注入，避免散落 magic string。

**备选**: 单一 `CloneOptions` 透传 go-git —— 与 eino Backend 抽象耦合过紧，不便文档化。

### 3. Ref 解析优先级

**选择**: `Commit` > `Tag` > `Branch`；三者皆空时使用 `HEAD`（即 clone 后默认分支）。

**理由**: commit 最精确；与 Git 用户心智一致。冲突时 commit 优先避免歧义。

### 4. 路径映射：virtual path ↔ repo path

**选择**:

1. 对外路径经 `normalizePath`（与 InMemoryBackend 一致：补 `/` 前缀、`filepath.Clean`）
2. `repoPath = filepath.Join(baseDir, strings.TrimPrefix(virtualPath, "/"))`，统一使用 `/` 作为 tree 内分隔符
3. 校验 `repoPath` 仍以 `baseDir` 为前缀（Clean 后），否则拒绝（防 `..` 穿越）

**示例**: `BaseDir=users/tom`，`virtual=/projects/todos` → `repo=users/tom/projects/todos`

**理由**: 对外 API 保持与 InMemoryBackend 一致的 `/` 根路径体验；baseDir 对调用方透明。

### 5. 只读方法实现策略

**选择**: 在指定 commit tree 上构建「路径 → tree entry」视图：

- **LsInfo**: 定位 tree 节点，遍历直接子 entry，返回相对 virtual path 的名称
- **Read**: 通过 blob hash 读取内容，复用 InMemoryBackend 的行 offset/limit 逻辑（可提取为内部 helper 或直接复制精简版）
- **GrepRaw**: 递归收集 baseDir 子树下所有 blob 文件，逐文件正则匹配（参考 InMemoryBackend 的 filter/grep 流程）
- **GlobInfo**: 递归 walk tree，用 `doublestar` 匹配（与 InMemoryBackend 一致，eino 已间接依赖）

**Write/Edit**: 返回 `errReadOnly`（包级 sentinel error，如 `errors.New("git backend is read-only")`）

### 6. 认证实现

**选择**:
- HTTP + Token: `http.BasicAuth{Username: "token-or-user", Password: token}` 或 go-git 的 `transport/http.BasicAuth`
- HTTP + User/Pass: 标准 Basic Auth
- SSH: `ssh.NewPublicKeys` 或 `ssh.NewPublicKeysFromFile`，通过 `git.CloneOptions.Auth` 注入

敏感信息不在 `Error()` 字符串中回显；日志由调用方控制，Backend 内不使用 `fmt.Printf` 打印凭据。

### 7. 初始化与懒加载

**选择**: `NewGitBackend(cfg)` 校验配置；首次只读操作触发 `sync.Once` 执行 clone + resolve ref + 缓存 root tree/commit。

**理由**: 构造时不必阻塞网络；配置错误尽早返回；同一实例内 tree 只解析一次。

**备选**: 构造时立即 clone —— 对仅校验配置的场景不友好。

### 8. 测试策略

**选择**: 使用 `go-git` 在测试中创建 memory bare repo，写入已知目录结构（含 baseDir 子树），不访问真实 GitHub/GitLab。

**理由**: CI 稳定、无网络依赖；可覆盖 ref、baseDir、路径穿越、只读拒绝等场景。

## Risks / Trade-offs

- **[Risk] 大仓库 clone 耗时长、占内存** → 首版接受；文档说明适用只读浏览场景；后续可加 shallow clone
- **[Risk] go-git SSH/HTTP 与官方 git 行为差异** → 集成测试覆盖两种协议；问题场景回退文档说明
- **[Risk] Grep/Glob 在大 tree 上性能** → 首版全量 walk；与 InMemoryBackend 规模假设一致
- **[Risk] ModifiedAt 来自 commit author/committer 时间，非文件 mtime** → 文档说明；目录 ModifiedAt 取子项最新时间或 commit 时间
- **[Trade-off] 不支持写** → Agent 只读场景足够；写操作需另开 change

## Migration Plan

纯新增包内实现，无迁移。发布新版本 goutils 后，调用方通过 `git.NewGitBackend(cfg)` 注入 Eino filesystem 工具。

## Open Questions

- SSH `KnownHosts` 首版是否强制要求？建议：生产由调用方传入；测试允许 `InsecureIgnoreHostKey`
- HTTP Token 与 Username 同时存在时，Token 优先作为 password、Username 固定为 `"oauth2"` 或 `"x-access-token"`（按 Git 托管方惯例）—— 实现时在文档中明确各平台示例
