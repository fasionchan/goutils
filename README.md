# goutils

Golang utility library. Package docs on pkg.go.dev:

- [github.com/fasionchan/goutils](https://pkg.go.dev/github.com/fasionchan/goutils)
- [stl](https://pkg.go.dev/github.com/fasionchan/goutils/stl) — generics helpers for slices, maps, sets, seq, buffers, cache fetchers, graphs, channels

## Agent skill: goutils-stl (English)

For coding agents that need **when/how to use `stl`** (not a full API dump):

- In-repo path: `.agents/skills/goutils-stl/`
- Cursor symlink: `.cursor/skills/goutils-stl` → that directory
- Load **either** `references/en/` **or** `references/zh/` per session (never both)

### Install into another project

Copy or symlink the skill directory:

```bash
cp -R .agents/skills/goutils-stl /path/to/other/.agents/skills/goutils-stl
# or
ln -s /abs/path/to/goutils/.agents/skills/goutils-stl \
  /path/to/other/.agents/skills/goutils-stl
```

Import from GitHub (Multica):

```bash
multica skill import --url \
  https://github.com/fasionchan/goutils/tree/master/.agents/skills/goutils-stl
```

## 智能体 Skill：goutils-stl（中文）

面向编程智能体的 `stl` **选用与配方**说明（非全量 API 镜像）：

- 仓内路径：`.agents/skills/goutils-stl/`
- Cursor 软链：`.cursor/skills/goutils-stl`
- 单次会话只加载 `references/zh/` **或** `references/en/` 之一

### 安装到其他项目

复制或软链 Skill 目录；或用上方 GitHub URL 通过 `multica skill import` 导入。
人类可读包概述见 [pkg.go.dev/stl](https://pkg.go.dev/github.com/fasionchan/goutils/stl)。
