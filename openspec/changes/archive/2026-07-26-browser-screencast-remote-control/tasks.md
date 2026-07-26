## 1. 协议类型与消息契约（remote.go）

- [x] 1.1 在 `libs/browser/remote.go` 定义 JSON envelope 与常量（`v`、`type` 命名空间、error code）
- [x] 1.2 定义 C→S payload 类型：`mouse.*`（move/down/up/wheel）、`key.*`（down/up/press）、`session.ping`
- [x] 1.3 定义 S→C payload 类型：`session.ready`、`screencast.meta`、`session.pong`、`ack`、`error`
- [x] 1.4 实现 envelope 解析/序列化辅助方法（未知 payload 字段忽略；非法必要字段返回明确错误）

## 2. Browser / RodBrowser 最小扩展

- [x] 2.1 在 `libs/browser/browser.go` 增加坐标级鼠标/键盘事件类型与 `Browser` 接口方法（不改现有 Click/Type 等签名）
- [x] 2.2 在 `Browser` 增加获取 screencast/session meta 的方法（viewport CSS 尺寸、format、可选 frame 尺寸）
- [x] 2.3 在 `libs/browser/rod.go` 用 CDP `Input.dispatchMouseEvent` / `Input.dispatchKeyEvent` 实现输入派发
- [x] 2.4 在 `RodBrowser` 实现 meta 读取（layout metrics + screencast options）；同步补齐测试 mock（若有）

## 3. RemoteController 接线

- [x] 3.1 完善 `ServeHTTP`：启动 screencast 后发送 `session.ready` 与 `screencast.meta`
- [x] 3.2 Binary 写帧 goroutine 与 Text 读循环分离；入站 Text 按 `type` 分发
- [x] 3.3 实现鼠标/键盘 handler，调用 `Browser` 新方法；支持 `key.press` = down+up
- [x] 3.4 实现 `session.ping` → `session.pong`；带 `id` 时成功回 `ack`、失败回 `error`
- [x] 3.5 对未知或预留未实现 type（如 `nav.*`）返回 `error.code=unsupported`，不静默成功

## 4. 测试与验证

- [x] 4.1 为 envelope 解析与 type 分发编写单测（合法消息、未知字段、缺字段、unsupported type）
- [x] 4.2 为 Rod 输入派发或可测的适配层补充测试（能跑则跑；依赖浏览器的可标集成/跳过策略与现有 `rod_test` 一致）
- [x] 4.3 运行相关 `go test`，确保既有 browser 测试未因接口扩展失败
