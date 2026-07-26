import { StrictMode, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRemoteViewer } from "@fasionchan/browser-remote-react";

function DemoApp() {
  const [wsUrlInput, setWsUrlInput] = useState(
    "ws://localhost:8080/Tabs/<tabId>/Remote?format=jpeg&quality=80",
  );
  const [activeUrl, setActiveUrl] = useState("");
  const [interactive, setInteractive] = useState(true);
  const [status, setStatus] = useState("idle");

  const viewer = useMemo(() => {
    if (!activeUrl) {
      return (
        <div style={{ padding: 24, color: "#666" }}>
          填写 Remote WebSocket URL 后点击 Connect
        </div>
      );
    }
    return (
      <BrowserRemoteViewer
        key={`${activeUrl}:${interactive ? "1" : "0"}`}
        wsUrl={activeUrl}
        interactive={interactive}
        style={{ width: "100%", height: "100%" }}
        onReady={(info) =>
          setStatus(`ready tab=${info.tabId} ${info.meta.viewport_width}x${info.meta.viewport_height}`)
        }
        onError={(err) => setStatus(`error: ${err.message}`)}
        onClose={() => setStatus("closed")}
      />
    );
  }, [activeUrl, interactive]);

  return (
    <div
      style={{
        fontFamily: "ui-sans-serif, system-ui, sans-serif",
        height: "100vh",
        display: "flex",
        flexDirection: "column",
        margin: 0,
      }}
    >
      <header
        style={{
          display: "flex",
          gap: 8,
          alignItems: "center",
          padding: 12,
          borderBottom: "1px solid #ddd",
          flexWrap: "wrap",
        }}
      >
        <strong>browser-remote-react demo</strong>
        <input
          value={wsUrlInput}
          onChange={(e) => setWsUrlInput(e.target.value)}
          style={{ flex: 1, minWidth: 280, padding: "6px 8px" }}
          placeholder="ws://localhost:8080/Tabs/{tabId}/Remote"
        />
        <label style={{ display: "flex", alignItems: "center", gap: 4 }}>
          <input
            type="checkbox"
            checked={interactive}
            onChange={(e) => setInteractive(e.target.checked)}
          />
          interactive
        </label>
        <button type="button" onClick={() => setActiveUrl(wsUrlInput.trim())}>
          Connect
        </button>
        <button type="button" onClick={() => setActiveUrl("")}>
          Disconnect
        </button>
        <span style={{ color: "#555", fontSize: 13 }}>{status}</span>
      </header>
      <main style={{ flex: 1, minHeight: 0 }}>{viewer}</main>
    </div>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <DemoApp />
  </StrictMode>,
);
