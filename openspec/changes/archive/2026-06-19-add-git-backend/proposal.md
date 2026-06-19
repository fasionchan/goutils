## Why

goutils 的 `libs/einox/backend/git` 已声明 `GitBackend` 实现 `filesystem.Backend`，但尚未实现任何方法。Eino ADK 的文件系统工具需要可插拔 Backend 来访问代码仓库内容；当前缺少对远程 Git 仓库的只读访问能力，无法让 Agent 直接读取指定 branch/tag/commit 下的文件，而不依赖本地 checkout。

## What Changes

- 在 `libs/einox/backend/git/git.go` 实现 `filesystem.Backend` 接口的**只读**方法：`LsInfo`、`Read`、`GrepRaw`、`GlobInfo`
- 支持通过 HTTP(S) 或 SSH 两种方式访问远程 Git 仓库
- 提供结构化的认证配置：HTTP 用户名/密码、Personal Access Token、SSH 私钥及可选 passphrase
- 支持指定 ref 类型：branch、tag 或 commit（三者互斥，至少指定其一）
- 支持 `baseDir` 基路径：对外暴露的路径以 `/` 为根，内部映射到仓库内 `baseDir` 子树（如 `baseDir=users/tom` 时，访问 `/projects/todos` 实际读取 `users/tom/projects/todos`）
- `Write`、`Edit` 及 `Shell.Execute` 等写/执行操作**暂不实现**，调用时返回明确的不支持错误
- 复用项目已有 `go-git` 依赖，补充单元测试（含本地 bare repo fixture）

## Capabilities

### New Capabilities

- `git-backend`: 基于 Git 仓库的只读 `filesystem.Backend` 实现，支持 HTTP/SSH 认证、ref 指定与 baseDir 路径映射

### Modified Capabilities

（无）

## Impact

- **代码**: `libs/einox/backend/git/` 新增/扩展 `git.go` 及配套文件（config、auth、path、read-only 实现、测试）
- **依赖**: 复用已有 `github.com/go-git/go-git/v6`、`github.com/cloudwego/eino`；无新增外部依赖
- **API**: 新增 `GitBackend` 构造与配置类型，纯新增，无 breaking change
- **调用方**: 可被 Eino ADK filesystem 工具及 opsys 相关 Agent 场景引用，用于远程仓库只读浏览与搜索
