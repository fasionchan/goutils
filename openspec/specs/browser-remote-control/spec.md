# browser-remote-control

基于 screencast 的浏览器远程控制 WebSocket 协议与服务端处理（帧推送、输入注入、会话 meta、可扩展 JSON 消息）。

## Requirements

### Requirement: WebSocket 双通道传输

远程控制连接 SHALL 使用 WebSocket。服务端 MUST 将 screencast 图像以 BinaryMessage 推送（载荷为图像原始字节，不含 JSON 包装或自定义 binary header）。在每一帧 BinaryMessage 之前，服务端 MUST 发送对应的 `screencast.frame` TextMessage（见「Screencast 帧标记」）。所有控制、会话与事件消息 MUST 使用 TextMessage，且正文为 JSON envelope。

#### Scenario: 推送 screencast 帧

- **WHEN** screencast 产生一帧图像数据
- **THEN** 服务端先发送 `screencast.frame` TextMessage，再以 BinaryMessage 将该帧字节写入 WebSocket，Binary 载荷不经过 JSON 编码

#### Scenario: 控制消息为 JSON Text

- **WHEN** 客户端发送鼠标或键盘控制指令
- **THEN** 指令 MUST 以 TextMessage 承载的 JSON envelope 发送

### Requirement: Screencast 帧标记

服务端在推送每一帧 screencast 图像之前，MUST 先发送一条 TextMessage JSON envelope，且 `type` MUST 为 `screencast.frame`。该消息的 `payload` MUST 包含单调递增的 `seq`（非负整数），并 MUST 包含图像 `format`（与当前 screencast 输出一致，如 `jpeg` / `png`）。`payload` MAY 包含 `ts`（Unix 毫秒）。紧随该 TextMessage 的下一条与该帧对应的图像数据 MUST 仍以 BinaryMessage 推送，且载荷 MUST 为图像原始字节（不得在 Binary 载荷内附加自定义 header）。

#### Scenario: 帧标记后紧跟 Binary 图像

- **WHEN** screencast 产生一帧图像数据
- **THEN** 服务端先发送 `type=screencast.frame` 的 JSON TextMessage，再以 BinaryMessage 发送该帧原始字节

#### Scenario: seq 单调递增

- **WHEN** 同一远程控制连接上连续推送多帧
- **THEN** 各 `screencast.frame` 的 `payload.seq` 按推送顺序严格递增

### Requirement: 版本化 JSON Envelope

所有 TextMessage JSON MUST 包含字段：`v`（协议版本，本期为 `1`）、`type`（`namespace.action` 字符串）、`payload`（object）。可选字段：`id`（关联 ID）、`ts`（Unix 毫秒时间戳）。接收方 MUST 忽略 `payload` 中未识别的字段，以实现前向兼容。

#### Scenario: 合法鼠标消息

- **WHEN** 客户端发送包含 `v=1`、`type=mouse.down`、以及含 `x`/`y`/`button` 的 `payload` 的 JSON
- **THEN** 服务端能够解析并为目标 tab 派发对应鼠标事件

#### Scenario: 未知 payload 字段不导致失败

- **WHEN** 客户端在已知 `type` 的 `payload` 中附加额外可选字段
- **THEN** 服务端 MUST 忽略未知字段并继续处理已知字段

### Requirement: 会话就绪与 Screencast Meta

WebSocket 升级且 screencast 启动成功后，服务端 MUST 向客户端发送 `session.ready`，并至少发送一次 `screencast.meta`。`screencast.meta` 的 `payload` MUST 包含页面 viewport 尺寸（CSS 像素），并 SHOULD 包含帧尺寸、图像 format、device scale factor，供客户端做坐标换算。

#### Scenario: 连接成功下发 meta

- **WHEN** 客户端成功建立远程控制 WebSocket 且 screencast 启动成功
- **THEN** 客户端收到 `type=session.ready` 与 `type=screencast.meta` 的 JSON 消息

### Requirement: 鼠标坐标为页面 CSS 像素

所有 `mouse.*` 消息中的 `x`、`y` MUST 表示目标页面的 CSS 像素坐标（非客户端 display 坐标、非未换算的帧像素）。服务端 MUST 按该坐标系直接注入，不得假设客户端上报的是 display 坐标。

#### Scenario: 点击页面中心

- **WHEN** 页面 viewport 为 1280×720，客户端在换算后发送 `mouse.down` 且 `x=640`、`y=360`
- **THEN** 服务端在页面 CSS 坐标 (640, 360) 处派发鼠标按下事件

### Requirement: 鼠标输入消息

服务端 MUST 支持以下客户端 `type`：`mouse.move`、`mouse.down`、`mouse.up`、`mouse.wheel`。`mouse.down`/`mouse.up` MUST 接受 `button`（与现有浏览器鼠标按钮常量语义一致）。`mouse.wheel` MUST 接受 `delta_x`、`delta_y`。

#### Scenario: 鼠标按下与抬起

- **WHEN** 客户端依次发送合法的 `mouse.down` 与 `mouse.up`
- **THEN** 目标页面收到对应的鼠标按下与抬起注入

#### Scenario: 滚轮

- **WHEN** 客户端发送带非零 `delta_y` 的 `mouse.wheel`
- **THEN** 目标页面在指定坐标处收到滚轮注入

### Requirement: 键盘输入消息

服务端 MUST 支持以下客户端 `type`：`key.down`、`key.up`、`key.press`。`key.press` MUST 等价于依次派发 key down 与 key up。载荷 MUST 能表达按键身份（`key` 和/或 `code`），并可包含可选 `text` 与修饰键信息。

#### Scenario: 按下可打印字符

- **WHEN** 客户端发送 `key.press`，`payload` 含可打印字符所需字段
- **THEN** 目标页面收到对应的键盘注入

#### Scenario: 修饰键组合

- **WHEN** 客户端发送带 modifiers（如 ctrl/shift）的 `key.down`
- **THEN** 服务端在派发时保留修饰键状态信息

### Requirement: 请求关联的 ack 与 error

当客户端消息包含非空 `id` 时：处理成功 MUST 回复 `type=ack` 且回传同一 `id`；处理失败 MUST 回复 `type=error`，包含同一 `id`、`code` 与 `message`。当消息无 `id` 时，成功路径 MUST NOT 强制发送 `ack`（以降低高频 `mouse.move` 带宽）。

#### Scenario: 带 id 的成功输入

- **WHEN** 客户端发送带 `id` 的合法 `mouse.down` 且注入成功
- **THEN** 服务端回复 `type=ack` 且 `id` 与请求一致

#### Scenario: 非法 payload

- **WHEN** 客户端发送 `type=mouse.down` 但缺少必要坐标字段
- **THEN** 若请求含 `id`，服务端回复 `type=error`；注入不得以错误坐标执行

### Requirement: 未知与未实现 type 的处理

对于客户端发来的未知 `type`，或已预留但本期未实现的 type（如 `nav.*`），服务端 MUST 回复 `type=error` 且 `code` 表明不支持（例如 `unsupported`），不得静默当成成功。

#### Scenario: 预留导航命令本期未实现

- **WHEN** 客户端发送 `type=nav.go`（无论是否带 `id`）
- **THEN** 服务端回复 `error`（有 `id` 则回传），且不执行导航

### Requirement: 为后续浏览器事件预留命名空间

协议 MUST 保留 `event.*` 与 `nav.*` 命名空间供后续扩展（如 `event.navigated`、`event.console`、`event.dialog`、`nav.reload`）。本期实现 MUST NOT 占用这些 type 表达其他语义。服务端本期 MUST NOT 主动推送导航/console 等事件（可在后续变更中启用）。

#### Scenario: 本期不推送 console 事件

- **WHEN** 页面产生 console 输出
- **THEN** 本期远程控制连接 MUST NOT 因此发送 `event.console` 消息

### Requirement: Browser 坐标级输入能力

`Browser` 接口 MUST 提供坐标级鼠标与键盘派发能力，供 `RemoteController` 调用。`RodBrowser` MUST 通过 CDP 输入域实现上述能力。现有基于 selector 的 `Click`/`Type`/`Hover` 等方法的签名与语义 MUST 保持不变。

#### Scenario: 通过接口注入鼠标

- **WHEN** 调用方通过 `Browser` 新方法对某 tab 派发鼠标事件
- **THEN** `RodBrowser` 将该事件注入对应页面且不依赖 CSS selector

#### Scenario: 既有 Click API 不受影响

- **WHEN** 调用现有 `Browser.Click`
- **THEN** 行为与本变更前一致
