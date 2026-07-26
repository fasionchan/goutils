## Context

`libs/browser/remote.go` 已有 `RemoteController`：升级 WebSocket、调用 `StartScreencast`、以 BinaryMessage 推送帧数据，但入站 JSON 被丢弃。`Browser` 仅有 selector 级交互（`Click`/`Type`/`Hover`），缺少坐标级鼠标/键盘注入，也无法向客户端提供 viewport/frame 元数据以完成坐标换算。

约束（已与需求方确认）：

- 协议与远程控制主逻辑落在 `remote.go`；允许小范围扩展 `Browser` / `RodBrowser`
- Screencast 图像走 Binary frame；控制与事件走 JSON Text frame
- 鼠标坐标约定为**页面 CSS 像素**（方案 A：客户端依据 meta 换算）
- 导航 / console 等事件本期不实现，但 JSON `type` 命名空间需预留扩展

## Goals / Non-Goals

**Goals:**

- 定义版本化、可扩展的 WebSocket JSON 消息契约，并在服务端完成解析与分发
- Binary 持续推送 screencast 图像字节；连接后下发 `session.ready` / `screencast.meta`
- 支持鼠标（move/down/up/wheel）与键盘（down/up/press）注入到目标 tab
- 对带 `id` 的请求返回 `ack` / `error`，便于客户端关联结果
- 以最小改动为 `Browser`/`RodBrowser` 增加坐标级输入与 meta 读取能力

**Non-Goals:**

- 本期不推送 navigation / console / dialog / download 等浏览器事件（仅预留 type）
- 不实现鉴权、多控制端冲突仲裁、音频、剪贴板文件传输
- 不改动 HTTP REST Tab API 行为；不做前端客户端实现
- 不引入新的外部依赖

## Decisions

### 1. 传输分层：Binary 帧 vs JSON 控制面

**选择**:

| WebSocket Opcode | 用途 |
|---|---|
| BinaryMessage | 仅 screencast 图像原始字节（jpeg/png，由 query `format` 决定） |
| TextMessage | 全部控制/会话/事件 JSON |

**理由**: 图像走二进制避免 base64 膨胀与 JSON 解析开销；控制面保持可读、可演进。

**备选**: 单 JSON 内嵌 base64 帧 —— 带宽与延迟更差，否决。

### 2. JSON Envelope（版本化信封）

**选择**: 所有 Text 消息统一为：

```json
{
  "v": 1,
  "id": "optional-correlation-id",
  "type": "namespace.action",
  "ts": 1710000000000,
  "payload": { }
}
```

| 字段 | 类型 | 说明 |
|---|---|---|
| `v` | int | 协议主版本；本期固定 `1`。不兼容变更时递增 |
| `id` | string, optional | 客户端请求关联 ID；服务端 `ack`/`error` 原样回传。服务端主动推送可省略 |
| `type` | string | `namespace.action` 点分命名；未知 type → `error`（`code=unsupported`） |
| `ts` | int64, optional | Unix 毫秒时间戳；服务端推送建议填充 |
| `payload` | object | 与 `type` 对应的载荷；未知字段 MUST 忽略（前向兼容） |

**扩展策略**:

- 新增能力优先加新 `type`，不复用旧 type 改语义
- 在既有 payload 上只**追加可选字段**；删除/改名视为 breaking，需升 `v`
- 命名空间预留：`session.*`、`screencast.*`、`mouse.*`、`key.*`、`nav.*`、`event.*`、`ack`、`error`

**备选**: 按消息种类拆多个顶层 schema、无统一 envelope —— 不利于版本与关联 ID，否决。

### 3. 本期实现的 type 清单

**客户端 → 服务端（C→S）**

| type | payload 要点 | 行为 |
|---|---|---|
| `session.ping` | 可选空对象 | 回 `session.pong`（可带同一 `id`） |
| `mouse.move` | `x`,`y`；可选 `modifiers` | 派发移动 |
| `mouse.down` / `mouse.up` | `x`,`y`,`button`；可选 `click_count`,`modifiers` | 按下/抬起 |
| `mouse.wheel` | `x`,`y`,`delta_x`,`delta_y`；可选 `modifiers` | 滚轮 |
| `key.down` / `key.up` | `key` 或 `code`；可选 `text`,`modifiers`,`auto_repeat` | 键按下/抬起 |
| `key.press` | 同 `key.down` 字段 | 服务端合成 down+up（便捷 API） |

**服务端 → 客户端（S→C）**

| type | payload 要点 | 时机 |
|---|---|---|
| `session.ready` | `tab_id`, `protocol_version` | 连接建立且 screencast 启动成功后 |
| `screencast.meta` | `format`, `viewport_width`, `viewport_height`, `frame_width`, `frame_height`, `device_scale_factor` | ready 时至少发一次；视口变化时可再发（本期若无法检测变化，仅初始一次） |
| `session.pong` | 可选空 | 响应 ping |
| `ack` | `ref_type`（原请求 type） | 对带 `id` 且处理成功的输入/会话请求 |
| `error` | `code`, `message`, 可选 `ref_type` | 解析失败、不支持、注入失败等 |

**预留（本期对显式客户端请求回 unsupported）**

| type | 方向 | 未来用途 |
|---|---|---|
| `nav.go` / `nav.back` / `nav.forward` / `nav.reload` | C→S | 导航控制 |
| `event.navigated` | S→C | URL/标题变化 |
| `event.console` | S→C | console API |
| `event.dialog` | S→C | alert/confirm/prompt |
| `event.tab_closed` | S→C | 目标 tab 关闭 |

### 4. 坐标约定（方案 A）

**选择**: 所有 `mouse.*` 的 `x`/`y` MUST 为**页面 CSS 像素**（与 CDP `Input.dispatchMouseEvent` 一致）。

客户端职责：

1. 读取 `screencast.meta` 中的 `viewport_*` 与展示元素尺寸
2. 将 pointer 事件从 display 坐标换算为 viewport CSS 像素后再发送

服务端职责：

1. 连接后提供准确 meta（viewport 来自页面 layout/metrics；frame 尺寸在可知时填充，否则可与 viewport 相同并文档说明）
2. **不再**根据 display 尺寸做二次映射

**备选**: 服务端换算 display/frame 坐标 —— 需客户端上报 display 布局，letterbox 易点偏，已否决。

### 5. Browser 接口最小扩展

**选择**: 在 `Browser` 增加专用方法：

- `DispatchMouseEvent(id string, event *MouseEvent) error`
- `DispatchKeyEvent(id string, event *KeyEvent) error`
- `GetScreencastSessionMeta(id string, opts *ScreencastOptions) (*ScreencastSessionMeta, error)`

`RodBrowser` 通过 CDP `Input.dispatchMouseEvent` / `Input.dispatchKeyEvent` 与 layout metrics 实现。`remote.go` 只依赖 `Browser` 接口，不 type-assert rod。

**不改动**: 现有 `Click`/`Type`/`Hover`/`StartScreencast` 签名与语义保持不变。

**备选**: 仅在 `remote.go` type-assert `*RodBrowser` —— 破坏抽象，否决。

### 6. `RemoteController` 运行模型

**选择**:

1. Upgrade → 解析 screencast query → `StartScreencast`
2. 获取并发送 `session.ready` + `screencast.meta`
3. goroutine：从 `frames.BytesChan` 写 BinaryMessage
4. 主循环：`ReadMessage`；Binary 入站忽略；Text 反序列化为 envelope 后 switch `type` 分发
5. 连接关闭 / 读错误 → stop screencast、退出

对高频 `mouse.move`：允许不强制 `id`；若无 `id` 则成功时不发 `ack`（降低带宽）。有 `id` 则必须 `ack`/`error`。

### 7. 字段与枚举对齐

- 鼠标 `button`: 复用现有常量语义（`left`/`middle`/`right`/`back`/`forward`/`none`）
- `modifiers`: `[]string`（`alt`/`ctrl`/`meta`/`shift`）
- 键盘：同时支持 `key`（如 `a`、`Enter`）与 `code`（如 `KeyA`）；`text` 用于可打印字符输入

### 8. HTTP 路由注册

在 chiopenapi Tab 路由下注册 `GET /Remote`（WebSocket 升级），由 `GetBrowserFromRequest.RegisterChiOpenApiRoutes` 挂载，覆盖 instance 与 pool：

- instance：`/Tabs/{tabId}/Remote`
- pool：`/Instances/{instanceId}/Tabs/{tabId}/Remote`

## Risks / Trade-offs

- [Risk] viewport 与 screencast 实际缩放不同步导致点偏 → Mitigation：meta 明确区分 viewport/frame；文档要求客户端用 viewport 换算；后续可在每 N 帧刷新 meta
- [Risk] 接口扩展迫使所有 `Browser` mock/实现补方法 → Mitigation：方法少、语义清晰；测试 mock 同步补齐
- [Risk] 高频 mouse.move 阻塞写循环 → Mitigation：读循环与写帧分 goroutine；输入调用尽快返回；必要时后续加合并（本期不做）
- [Trade-off] 预留 type 本期返回 unsupported vs 静默忽略 → 选择对**显式客户端请求**回 `error.code=unsupported`，避免误以为已生效；对纯服务端未来推送无影响

## Migration Plan

- 纯新增协议与接口方法，无数据迁移
- 旧客户端若只收 Binary 帧、不发 JSON，行为与现网推帧兼容
- 回滚：去掉远程分发逻辑即可，不影响 selector API

## Open Questions

（需求已确认，无阻塞项。实现阶段若 CDP meta 字段获取成本高，可先填 viewport，frame 尺寸暂与 viewport 相同并在注释标明。）
