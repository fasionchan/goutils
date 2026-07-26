## 1. 协议：screencast.frame 帧标记（Go）

- [x] 1.1 在 `libs/browser/remote.go` 增加 `RemoteTypeScreencastFrame` 与 `RemoteScreencastFramePayload`（`seq`、`format`，可选 `ts`）
- [x] 1.2 修改 `RemoteController` 推帧路径：每帧 Binary 前先 `writeEnvelope` 发送 `screencast.frame`（维护连接级 seq）
- [x] 1.3 补充/更新 `remote_test.go`（及必要时 `apiserver_test.go`）：断言 frame 标记顺序与 Binary 仍为原始图像字节
- [x] 1.4 运行相关 Go 测试并通过

## 2. 前端包脚手架（libs/browser/web）

- [x] 2.1 创建 `libs/browser/web/package.json`、tsconfig、构建配置（Vite/tsup 等），React 为 peerDependency
- [x] 2.2 定义协议 TypeScript 类型与编解码辅助（envelope、`screencast.meta` / `screencast.frame`、mouse/key payload）
- [x] 2.3 实现坐标换算纯函数（letterbox + viewport meta → CSS 像素）并加单测

## 3. React Viewer 组件

- [x] 3.1 实现 WebSocket 会话状态机：ready/meta/frame+binary、错误与关闭回调；支持无 frame 标记时的降级渲染
- [x] 3.2 实现帧渲染（Blob → canvas / createImageBitmap），正确处理 jpeg/png
- [x] 3.3 实现 interactive 输入：pointer/wheel/key → `mouse.*` / `key.*`；move 节流；黑边忽略
- [x] 3.4 导出 `BrowserRemoteViewer`（及可选 `useBrowserRemote`），编写包 README（props、wsUrl 示例）

## 4. goutils 内 Demo

- [x] 4.1 在 `libs/browser/web/demo`（或等价路径）提供 Vite demo 页：填写/选择 tab Remote URL，嵌入 Viewer
- [x] 4.2 编写启动说明（启动 `browserd`、创建 tab、打开 demo）
- [x] 4.3 （可选增强）`browserd` 托管构建后的 demo 静态文件，便于单端口体验

## 5. 验收

- [ ] 5.1 本地端到端：画面可见，点击/拖动/滚轮/键盘作用于远程页面
- [x] 5.2 确认只读模式（`interactive=false`）不发送输入消息
