## Context

`RemoteController`（`libs/browser/remote.go`）已实现：WebSocket 升级、Binary 推帧、JSON envelope 输入（鼠标/键盘）、`session.ready` / `screencast.meta`。缺少前端 viewer。调研结论：不采用 Atrium / noVNC 整包；自研轻量 React 组件，复用 canvas + letterbox 坐标映射等常见模式。

已确认约束：

1. 自研轻量组件（不绑定第三方远程浏览器产品）
2. 交付：**可发布前端包** + **goutils 内 demo** 都要
3. 一期：仅看画面 + 鼠标/键盘（无地址栏 / nav chrome）
4. 协议可加帧标记；优先 JSON `screencast.frame`，不加 binary header

## Goals / Non-Goals

**Goals:**

- 可发布 React 包：给定 `wsUrl`（或 base + tabId）即可嵌入远程画面并交互
- goutils 内 demo：本地跑 `browserd` + 静态页即可验证端到端
- 服务端在每帧 Binary 前发送 `screencast.frame`（seq/ts/format）
- 客户端按 `screencast.meta` 的 viewport 做 CSS 像素坐标换算后发送 `mouse.*` / `key.*`

**Non-Goals:**

- 导航栏、前进/后退/刷新 UI（服务端 `nav.*` 仍 reserved）
- 多 tab 条、剪贴板、音频、鉴权 UI、控制权仲裁
- 改造成 Atrium / VNC / RFB 协议
- Binary 帧内嵌自定义 header

## Decisions

### 1. 包与目录布局

**选择**:

| 路径 | 角色 |
|---|---|
| `libs/browser/web/` | 可发布前端包根目录（`package.json`、`src/`、构建配置） |
| `libs/browser/web/src/` | React 组件、协议类型、hooks |
| `libs/browser/cmd/browserd/` 或 `libs/browser/web/demo/` | goutils 内 demo 页；开发时 Vite/静态，可由 browserd 同端口或独立端口托管 |

包名建议：`@fasionchan/browser-remote-react`（实现时按仓库 npm scope 惯例微调）。

**备选**: 放在仓库顶层 `packages/` —— 当前 goutils 无 packages 惯例，否决。

### 2. 组件 API（最小面）

**选择**: 导出主组件 `BrowserRemoteViewer`（名称可微调），核心 props：

- `wsUrl: string` — 完整 WebSocket URL（含 query，如 `format`/`quality`）
- `interactive?: boolean` — 默认 `true`；`false` 时只看不操作
- `onReady?` / `onError?` / `onClose?` — 生命周期回调
- `className?` / `style?` — 容器样式

内部职责：

1. 建连 → 解析 `session.ready` / `screencast.meta` / `screencast.frame`
2. Binary → `Blob` → canvas/`createImageBitmap` 渲染
3. pointer / wheel / key → 换算后发 JSON envelope
4. 可选低频 `session.ping`

**备选**: 拆成无 UI 的 `useBrowserRemote` + 纯展示组件 —— 一期可在包内同时导出 hook，主路径仍是一体组件。

### 3. 帧标记：`screencast.frame`（不用 binary header）

**选择**: 每推一帧图像时，顺序为：

1. TextMessage：`type=screencast.frame`，`payload` 含至少 `seq`（单调递增）、`format`（或与 meta 一致的 mime 提示），可选 `ts`
2. BinaryMessage：该帧原始 JPEG/PNG 字节

客户端状态机：收到 `screencast.frame` 后，下一条 Binary 归属该 frame；乱序/缺失时丢弃或跳过并打日志。

**理由**:

- 与现有「Binary = 纯图像」约束兼容，前端 `Blob`/`createImageBitmap` 无需解析自定义二进制头
- 提供 seq/诊断能力，接近 Atrium「JSON frame + binary」模式但不改 envelope 体系
- 不加 binary header，避免跨语言长度/endian 约定成本

**备选**: 仅靠 meta、无 per-frame 标记 —— 排障与多帧对齐较弱，否决（本期按需求加上）。  
**备选**: binary header —— 否决（前端解码负担大，收益有限）。

协议版本：保持 `v=1`（新增 type，属前向扩展）；文档标明新客户端应处理 `screencast.frame`。

### 4. 坐标换算

**选择**: 沿用服务端约定——鼠标坐标为 **viewport CSS 像素**。

客户端：

1. 用 `screencast.meta` 的 `viewport_width` / `viewport_height`
2. 画面以 letterbox（`object-fit: contain`）绘入容器
3. 指针落在「可见图像矩形」内时，映射到 viewport 坐标再发送；落在黑边则忽略

参考 Atrium / 同类 CDP viewer 的 letterbox 映射，但不引入其协议。

### 5. 输入消息映射

| DOM | 协议 |
|---|---|
| pointermove | `mouse.move`（可节流；默认不带 `id`） |
| pointerdown / pointerup | `mouse.down` / `mouse.up` |
| wheel | `mouse.wheel` |
| keydown / keyup | `key.down` / `key.up`（可打印字符带 `text`） |

`mouse.move` 默认无 `id`（减少 ack 流量）。需要可靠确认的操作可带 `id`。

### 6. Demo 集成

**选择**:

- 包内 `demo/`（Vite）开发时热更新，指向本地 `browserd`（如 `ws://localhost:8080/.../Tabs/{id}/Remote`）
- `browserd` 可选：在非 API 路径 serve 构建后的 demo 静态文件，一键体验

一期至少保证：**前端包可独立 `npm run demo` + browserd 已启动** 即可验收。

### 7. 测试策略

- Go：`remote` 单测断言帧前发送 `screencast.frame`，且 Binary 仍为原始字节
- 前端：协议解析与坐标换算单测（纯函数）；组件可用 mock WebSocket 做轻量测试
- 手工：demo 页点按、拖动、滚轮、打字

## Risks / Trade-offs

- [Risk] 帧标记与 Binary 之间被其他 Text 插入 → Mitigation：客户端仅在「期望 binary」状态下消费 Binary；非预期消息重置/忽略
- [Risk] letterbox 换算偏差 → Mitigation：单测覆盖常见缩放；meta 区分 viewport/frame
- [Risk] 高频 move 压垮服务端 → Mitigation：客户端节流（如 rAF）；服务端已对无 id 的 move 不发 ack
- [Trade-off] 包放在 `libs/browser/web` vs 独立仓库 → 选同仓便于协议联调；发布仍走独立 package.json
- [Trade-off] 加 frame 标记轻微增加带宽 → 可接受（小 JSON）

## Migration Plan

1. 服务端先发 `screencast.frame`（与 viewer 同步上线）
2. 更新 `openspec/specs/browser-remote-control`（archive 时 sync）
3. 发布前端包；demo 文档说明启动步骤
4. 回滚：去掉 frame 发送即可；新客户端可在无 frame 时降级为「裸 Binary + meta.format」兼容模式（建议实现降级）

## Open Questions

- npm 包最终 scope/名称（实现前按仓库惯例确认一次即可）
- demo 是否必须由 `browserd` 同进程托管静态文件，或仅独立 Vite 即可（默认：独立 Vite 足够，同进程为增强项）
