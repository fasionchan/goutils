## Why

`libs/browser` 已具备 WebSocket 远程控制协议（Binary screencast + JSON 输入），但缺少可嵌入的前端 viewer，无法在浏览器里「看画面并操作」。调研后确认没有可直接对接本协议的成熟 React 组件（Atrium 最接近但协议与产品模型不匹配；noVNC 依赖 VNC 栈），因此需要自研轻量 React 组件，并在 goutils 内提供可运行 demo。

## What Changes

- 新增可发布的 React 前端包：连接 `/Tabs/{tabId}/Remote`，渲染 screencast，转发鼠标/键盘（一期不含导航 chrome、多 tab UI）
- 在 goutils 内提供 demo（例如挂到 `browserd` 静态页或 `demos/`），用真实 Remote 端点验证端到端
- 协议小幅增强（必要时）：每帧 Binary 前增加 `screencast.frame` JSON 标记（seq/ts/format），便于前端对齐与排障；Binary 载荷仍为原始图像字节（不加 binary header，避免前端解码复杂化）
- 更新 `browser-remote-control` 规范以覆盖帧标记语义；保持现有 `mouse.*` / `key.*` / meta / ack 行为兼容

## Capabilities

### New Capabilities

- `browser-remote-react-viewer`: 面向 goutils browser Remote 协议的轻量 React viewer（连接、帧渲染、坐标换算、鼠标键盘交互、可发布包 + goutils 内 demo）

### Modified Capabilities

- `browser-remote-control`: 增加 `screencast.frame` 帧标记（TextMessage，置于对应 Binary 帧之前）；明确客户端可用 meta + frame 标记完成渲染与诊断

## Impact

- **Go**: `libs/browser/remote.go`（及测试）发送 `screencast.frame`；`browserd` 或 demo 路由可能托管静态资源
- **Frontend**: 新 npm 包（React peer dependency）；协议 TypeScript 类型与组件 API
- **依赖**: 前端侧尽量无浏览器原生 WebSocket / canvas / Blob，不引入 noVNC/Atrium 整包
- **兼容**: 无生产前端客户端时，帧标记为可接受的协议演进；旧「只收 Binary」客户端仍可忽略 Text 标记（若存在）
