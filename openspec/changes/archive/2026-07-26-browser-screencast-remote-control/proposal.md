## Why

`libs/browser` 已具备 WebSocket screencast 推帧雏形（`remote.go`），但尚未定义完整的远程控制协议，也无法将客户端鼠标/键盘输入注入到目标页面。需要在可扩展的消息结构上打通「看画面 + 操作页面」，供远程操控与调试场景使用。

## What Changes

- 在 `libs/browser/remote.go` 落地 WebSocket 远程控制协议：Binary frame 推送 screencast 图像；Text/JSON 承载控制与事件消息
- 定义版本化 JSON envelope（`v`/`id`/`type`/`ts`/`payload`），覆盖鼠标、键盘、会话元数据、ack/error，并为导航/console 等后续事件预留 `type` 命名空间
- 约定鼠标坐标为**页面 CSS 像素**（客户端依据服务端下发的 viewport/frame meta 自行换算）
- 小范围扩展 `Browser` 接口及 `RodBrowser` 实现：增加坐标级鼠标/键盘派发，以及获取远程会话所需的 viewport/screencast meta（严格控制改动面，不改无关 API）
- 完善 `RemoteController.ServeHTTP`：解析入站 JSON、分发到输入/会话处理、下发 `session.ready`/`screencast.meta`/`ack`/`error`；本期不实现导航与 console 等事件推送

## Capabilities

### New Capabilities

- `browser-remote-control`: 基于 screencast 的浏览器远程控制 WebSocket 协议与服务端处理（帧推送、输入注入、会话 meta、可扩展 JSON 消息）

### Modified Capabilities

（无）

## Impact

- **代码**: 主要 `libs/browser/remote.go`；小范围改动 `libs/browser/browser.go`（接口方法）、`libs/browser/rod.go`（CDP 输入/meta 实现）；必要时补充同包测试
- **协议**: 新增面向客户端的 WebSocket 消息契约（非 HTTP REST 变更）
- **依赖**: 复用现有 `gorilla/websocket`、rod/CDP；无新增外部依赖预期
- **兼容性**: 纯新增能力；现有 selector 级 `Click`/`Type` 等 API 保持不变
