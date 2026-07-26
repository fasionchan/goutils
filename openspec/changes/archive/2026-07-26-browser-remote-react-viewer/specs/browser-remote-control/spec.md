## ADDED Requirements

### Requirement: Screencast 帧标记

服务端在推送每一帧 screencast 图像之前，MUST 先发送一条 TextMessage JSON envelope，且 `type` MUST 为 `screencast.frame`。该消息的 `payload` MUST 包含单调递增的 `seq`（非负整数），并 MUST 包含图像 `format`（与当前 screencast 输出一致，如 `jpeg` / `png`）。`payload` MAY 包含 `ts`（Unix 毫秒）。紧随该 TextMessage 的下一条与该帧对应的图像数据 MUST 仍以 BinaryMessage 推送，且载荷 MUST 为图像原始字节（不得在 Binary 载荷内附加自定义 header）。

#### Scenario: 帧标记后紧跟 Binary 图像

- **WHEN** screencast 产生一帧图像数据
- **THEN** 服务端先发送 `type=screencast.frame` 的 JSON TextMessage，再以 BinaryMessage 发送该帧原始字节

#### Scenario: seq 单调递增

- **WHEN** 同一远程控制连接上连续推送多帧
- **THEN** 各 `screencast.frame` 的 `payload.seq` 按推送顺序严格递增

## MODIFIED Requirements

### Requirement: WebSocket 双通道传输

远程控制连接 SHALL 使用 WebSocket。服务端 MUST 将 screencast 图像以 BinaryMessage 推送（载荷为图像原始字节，不含 JSON 包装或自定义 binary header）。在每一帧 BinaryMessage 之前，服务端 MUST 发送对应的 `screencast.frame` TextMessage（见「Screencast 帧标记」）。所有控制、会话与事件消息 MUST 使用 TextMessage，且正文为 JSON envelope。

#### Scenario: 推送 screencast 帧

- **WHEN** screencast 产生一帧图像数据
- **THEN** 服务端先发送 `screencast.frame` TextMessage，再以 BinaryMessage 将该帧字节写入 WebSocket，Binary 载荷不经过 JSON 编码

#### Scenario: 控制消息为 JSON Text

- **WHEN** 客户端发送鼠标或键盘控制指令
- **THEN** 指令 MUST 以 TextMessage 承载的 JSON envelope 发送
