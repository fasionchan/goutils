# git-backend

基于 Git 仓库的只读 `filesystem.Backend` 实现，支持 HTTP/SSH 认证、ref 指定与 baseDir 路径映射。

## Requirements

### Requirement: 实现 filesystem.Backend 只读接口

`GitBackend` SHALL 实现 `github.com/cloudwego/eino/adk/filesystem.Backend` 接口中的只读方法：`LsInfo`、`Read`、`GrepRaw`、`GlobInfo`。所有方法 MUST 在指定 ref 对应的 commit tree 上操作，行为与 `InMemoryBackend` 语义对齐（路径以 `/` 为根、支持 offset/limit 读行、grep glob 等）。

#### Scenario: 列出目录内容

- **WHEN** 调用 `LsInfo`，路径指向 commit tree 中存在的目录
- **THEN** 返回该目录下直接子项的 `FileInfo` 列表（含文件名、是否目录、大小、修改时间）

#### Scenario: 读取文件内容

- **WHEN** 调用 `Read`，路径指向 commit tree 中的普通文件
- **THEN** 返回文件文本内容；支持 `Offset`（1-based 行号）与 `Limit`（最大行数），语义与 `InMemoryBackend.Read` 一致

#### Scenario: 搜索文件内容

- **WHEN** 调用 `GrepRaw`，提供合法正则 `Pattern` 及可选 `Path`、`Glob`、`FileType` 等过滤条件
- **THEN** 在 ref 对应 tree 的匹配文件中返回 `GrepMatch` 列表

#### Scenario: Glob 匹配文件

- **WHEN** 调用 `GlobInfo`，提供 glob `Pattern` 及可选基路径 `Path`
- **THEN** 返回匹配文件的 `FileInfo` 列表

### Requirement: 写操作与执行操作暂不支持

`Write`、`Edit` 以及 `Shell.Execute`（若实现 `Shell` 接口）MUST NOT 在本变更中实现。调用这些方法 SHALL 返回明确的「不支持」错误，不得静默成功或 panic。

#### Scenario: 写入文件被拒绝

- **WHEN** 调用 `Write` 或 `Edit`
- **THEN** 返回错误，表明 GitBackend 为只读模式

### Requirement: 支持 HTTP 与 SSH 两种仓库访问方式

`GitBackend` SHALL 支持通过 HTTP(S) URL 或 SSH URL 克隆/读取远程仓库。构造配置 MUST 明确指定传输协议类型，并根据协议应用对应认证。

#### Scenario: HTTP 访问公开仓库

- **WHEN** 配置 `Transport=HTTP` 且仓库 URL 为 HTTPS，无需认证
- **THEN** 成功 clone/fetch 并在指定 ref 上提供只读文件访问

#### Scenario: HTTP Token 认证

- **WHEN** 配置 `Transport=HTTP` 且提供 `Token`（Personal Access Token）
- **THEN** 使用该 token 作为 HTTP 认证凭据访问私有仓库

#### Scenario: HTTP 用户名密码认证

- **WHEN** 配置 `Transport=HTTP` 且提供 `Username` 与 `Password`
- **THEN** 使用 Basic Auth 访问私有仓库

#### Scenario: SSH 密钥认证

- **WHEN** 配置 `Transport=SSH` 且提供 SSH 私钥（及可选 passphrase）
- **THEN** 使用 SSH 公钥认证访问仓库

### Requirement: 认证信息结构化配置

认证相关字段 MUST 通过独立配置类型组织，不得硬编码在业务逻辑中。敏感字段（password、token、private key）SHALL 可通过构造函数注入，且不在日志或错误信息中明文输出。

#### Scenario: 认证配置与仓库 URL 分离

- **WHEN** 创建 `GitBackend`
- **THEN** 仓库 URL、传输协议、认证凭据、ref 选项、baseDir 分别通过配置字段指定

#### Scenario: 认证失败

- **WHEN** 提供的凭据无效或权限不足
- **THEN** 在 clone/fetch 或 tree 解析阶段返回明确错误，不暴露凭据内容

### Requirement: 支持指定 branch、tag 或 commit

`GitBackend` MUST 支持通过 ref 定位仓库快照。ref 类型为 branch、tag 或 commit hash 三者之一；若同时指定多个 ref 字段，SHALL 按优先级解析或返回配置错误。未指定 ref 时 MUST 使用仓库默认分支。

#### Scenario: 指定 branch

- **WHEN** 配置 `Branch="main"` 且未指定 tag/commit
- **THEN** 所有文件操作基于 `main` 分支 HEAD 对应的 tree

#### Scenario: 指定 tag

- **WHEN** 配置 `Tag="v1.0.0"` 且未指定 branch/commit
- **THEN** 所有文件操作基于该 tag 指向的 commit tree

#### Scenario: 指定 commit

- **WHEN** 配置 `Commit="<40-char-hash>"` 且未指定 branch/tag
- **THEN** 所有文件操作基于该 commit 的 tree

#### Scenario: ref 不存在

- **WHEN** 指定的 branch/tag/commit 在仓库中不存在
- **THEN** 返回明确错误，说明 ref 无法解析

### Requirement: baseDir 基路径映射

`GitBackend` SHALL 支持 `BaseDir` 配置，将对外虚拟根路径映射到仓库 tree 内的子目录。对外路径以 `/` 为根；内部实际路径为 `BaseDir` 与请求路径的拼接（去除多余 `/`）。`BaseDir` 为空时表示仓库 tree 根。

#### Scenario: baseDir 子树访问

- **WHEN** `BaseDir="users/tom"` 且调用 `Read`，`FilePath="/projects/todos/README.md"`
- **THEN** 实际读取仓库 tree 中 `users/tom/projects/todos/README.md`

#### Scenario: baseDir 根目录列表

- **WHEN** `BaseDir="users/tom"` 且调用 `LsInfo`，`Path="/"`
- **THEN** 列出 `users/tom` 目录下的直接子项，返回路径相对于虚拟根（不含 `users/tom` 前缀）

#### Scenario: 访问 baseDir 外路径

- **WHEN** `BaseDir="users/tom"` 且请求路径解析后超出 baseDir 子树（如 `../etc/passwd` 或 `../../`）
- **THEN** 返回路径非法或文件不存在错误，禁止目录穿越

#### Scenario: baseDir 不存在

- **WHEN** `BaseDir` 在指定 ref 的 tree 中不存在
- **THEN** 在首次访问或初始化时返回明确错误

### Requirement: 路径规范化

所有传入路径 MUST 经规范化处理：确保以 `/` 开头、使用 `filepath.Clean` 语义清理、禁止 `..` 逃逸出 baseDir 子树。行为与 `InMemoryBackend` 的 `normalizePath` 一致，并额外叠加 baseDir 映射。

#### Scenario: 相对路径规范化

- **WHEN** 请求路径为 `projects/todos`（无前导 `/`）
- **THEN** 规范化为 `/projects/todos` 后再映射到 baseDir 下

### Requirement: 仓库数据获取策略

Backend SHALL 通过 `go-git` 获取仓库对象并在内存中访问指定 commit 的 tree，无需 checkout 工作区。clone/fetch 策略 SHOULD 使用 bare + `NoCheckout` 以最小化 IO；同一 `GitBackend` 实例内 MAY 缓存已解析的 tree 以提升重复读性能。

#### Scenario: 首次访问触发 clone

- **WHEN** `GitBackend` 创建后首次调用只读方法
- **THEN** 执行 clone（或等效 fetch）并在指定 ref 上解析 tree，后续读操作复用已加载数据

#### Scenario: 远程不可达

- **WHEN** 仓库 URL 无法访问或网络失败
- **THEN** 返回包含根因的 error，不 panic
