# @fasionchan/browser-remote-react

面向 goutils `libs/browser` Remote WebSocket 协议的轻量 React viewer：渲染 screencast 画面，并转发鼠标/键盘输入。

## Install

```bash
npm install @fasionchan/browser-remote-react react react-dom
```

## Usage

```tsx
import { BrowserRemoteViewer } from "@fasionchan/browser-remote-react";

export function App() {
  return (
    <BrowserRemoteViewer
      wsUrl="ws://localhost:8080/Tabs/<tabId>/Remote?format=jpeg&quality=80"
      interactive
      style={{ width: 960, height: 540 }}
      onReady={(info) => console.log("ready", info)}
      onError={(err) => console.error(err)}
    />
  );
}
```

### Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `wsUrl` | `string` | — | Remote WebSocket URL |
| `interactive` | `boolean` | `true` | `false` 时只看画面，不发送输入 |
| `onReady` / `onError` / `onClose` | callbacks | — | 生命周期 |
| `className` / `style` | — | — | 容器样式 |

Also exported: `useBrowserRemote` hook, protocol helpers, and `mapPointerToViewport`.

## Protocol notes

- Text JSON envelope: `session.ready`, `screencast.meta`, `screencast.frame`, `mouse.*`, `key.*`
- Binary frames are raw JPEG/PNG bytes; each is preceded by `screencast.frame`
- Mouse coordinates are page CSS viewport pixels (letterbox mapping handled in the component)

## Local demo

1. Start browserd from goutils:

```bash
cd libs/browser/cmd/browserd
go run . instance
```

2. Create a tab (example):

```bash
curl -s -X POST http://localhost:8080/Tabs -H 'Content-Type: application/json' \
  -d '{"Url":"https://example.com"}'
```

记下返回的 tab id。

3. Run the demo UI:

```bash
cd libs/browser/web
npm install
npm run demo
```

Open the printed Vite URL, paste `ws://localhost:8080/Tabs/<tabId>/Remote`, then connect.

### Optional: serve demo from browserd

```bash
cd libs/browser/web
npm run build
npx vite build --config demo/vite.config.ts --outDir dist-demo

# from repo, start browserd with static dir
# from libs/browser/web after demo:build
DEMO_STATIC_DIR=$PWD/dist-demo go run ../cmd/browserd instance
```

Then open `http://localhost:8080/demo/`.

## Scripts

- `npm run build` — build library to `dist/`
- `npm test` — unit tests
- `npm run demo` — Vite demo app
