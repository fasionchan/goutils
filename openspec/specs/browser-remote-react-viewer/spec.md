# browser-remote-react-viewer

面向 goutils browser Remote 协议的可发布 React Viewer 包、画面渲染、鼠标键盘交互，以及仓库内本地 demo。

## Requirements

### Requirement: 可发布 React 远程 Viewer 包

仓库 MUST 在 `libs/browser/web/`（或设计文档约定的等价路径）提供可独立构建/发布的 React 包，导出可嵌入的远程浏览器画面组件（如 `BrowserRemoteViewer`）。该包 MUST 以 React 为 peer dependency，MUST 能通过传入 WebSocket URL 连接到 goutils browser Remote 端点。

#### Scenario: 通过 wsUrl 建立连接

- **WHEN** 应用渲染 viewer 并提供合法的 Remote WebSocket URL
- **THEN** 组件建立 WebSocket 连接并开始接收会话与帧消息

### Requirement: Screencast 画面渲染

Viewer MUST 消费服务端 `screencast.meta` 与 `screencast.frame` + BinaryMessage，将图像帧渲染到可视区域（canvas 或等价）。在收到 `screencast.frame` 后，MUST 将随后的 BinaryMessage 视为该帧图像数据。若连接仅提供裸 Binary 且已有 meta.format，组件 MAY 降级渲染以保持兼容。

#### Scenario: 正常帧流渲染

- **WHEN** 服务端依次发送 `screencast.meta`、`screencast.frame` 与对应 Binary 图像
- **THEN** viewer 在可视区域显示该帧画面

#### Scenario: 缺少 frame 标记时的降级

- **WHEN** 服务端在已发送 `screencast.meta` 后直接推送 Binary 图像（无 `screencast.frame`）
- **THEN** viewer MUST 仍尝试按 meta 中的 format 渲染该 Binary 帧

### Requirement: 鼠标与键盘交互

当 `interactive` 为启用（默认启用）时，viewer MUST 将用户在画面有效区域内的鼠标（move/down/up/wheel）与键盘（down/up）操作转换为 Remote 协议 JSON envelope 并发送。鼠标坐标 MUST 换算为页面 CSS viewport 像素后再发送，换算依据 `screencast.meta` 的 viewport 尺寸与画面 letterbox 布局。落在 letterbox 黑边外的指针事件 MUST 不发送鼠标消息。

#### Scenario: 点击换算为 viewport 坐标

- **WHEN** meta 标明 viewport 为 1280×720，画面 letterbox 显示且用户点击可见图像中心
- **THEN** 客户端发送的 `mouse.down`/`mouse.up` 的 `x`/`y` 接近 (640, 360)

#### Scenario: 只读模式不发送输入

- **WHEN** viewer 以 `interactive=false` 渲染
- **THEN** 用户在画面上的鼠标与键盘操作 MUST 不向 WebSocket 发送输入类消息

### Requirement: goutils 内 Demo

仓库 MUST 提供可运行的 demo（包内 demo 应用和/或由 `browserd` 托管的静态页），用于连接本地 Remote 端点并演示「看画面 + 鼠标键盘」。Demo 文档或 README MUST 说明启动 `browserd` 与打开 demo 的步骤。

#### Scenario: 本地端到端演示

- **WHEN** 开发者按文档启动 browserd 与 demo，并打开某一 tab 的 Remote URL
- **THEN** demo 页面能显示远程画面，且鼠标/键盘操作可作用于目标页面
