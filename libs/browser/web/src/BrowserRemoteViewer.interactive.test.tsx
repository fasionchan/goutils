import { act, useEffect, useRef } from "react";
import { createRoot, Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { BrowserRemoteViewer } from "./BrowserRemoteViewer";
import { encodeRemoteEnvelope, RemoteType } from "./protocol";

class MockWebSocket {
  static OPEN = 1;
  static instances: MockWebSocket[] = [];
  readyState = MockWebSocket.OPEN;
  binaryType = "arraybuffer";
  sent: string[] = [];
  onmessage: ((ev: MessageEvent) => void) | null = null;
  onerror: (() => void) | null = null;
  onclose: ((ev: CloseEvent) => void) | null = null;

  constructor(public url: string) {
    MockWebSocket.instances.push(this);
    queueMicrotask(() => {
      this.onmessage?.({
        data: encodeRemoteEnvelope(RemoteType.SessionReady, {
          tab_id: "t1",
          protocol_version: 1,
        }),
      } as MessageEvent);
      this.onmessage?.({
        data: encodeRemoteEnvelope(RemoteType.ScreencastMeta, {
          format: "jpeg",
          viewport_width: 1280,
          viewport_height: 720,
          frame_width: 1280,
          frame_height: 720,
        }),
      } as MessageEvent);
    });
  }

  send(data: string) {
    this.sent.push(data);
  }

  close() {
    this.readyState = 3;
  }
}

function Harness({ interactive }: { interactive: boolean }) {
  const hostRef = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    // ensure layout size for hit-testing
    if (hostRef.current) {
      Object.defineProperty(hostRef.current, "getBoundingClientRect", {
        value: () => ({
          left: 0,
          top: 0,
          width: 400,
          height: 400,
          right: 400,
          bottom: 400,
          x: 0,
          y: 0,
          toJSON: () => ({}),
        }),
      });
    }
  }, []);
  return (
    <div ref={hostRef} style={{ width: 400, height: 400 }}>
      <BrowserRemoteViewer
        wsUrl="ws://example.test/remote"
        interactive={interactive}
        style={{ width: 400, height: 400 }}
      />
    </div>
  );
}

describe("BrowserRemoteViewer interactive flag", () => {
  let root: Root;
  let container: HTMLDivElement;
  const OriginalWebSocket = globalThis.WebSocket;

  beforeEach(() => {
    // jsdom lacks PointerEvent; MouseEvent is enough for our handlers.
    if (typeof globalThis.PointerEvent === "undefined") {
      // @ts-expect-error polyfill
      globalThis.PointerEvent = class PointerEvent extends MouseEvent {
        constructor(type: string, props?: MouseEventInit) {
          super(type, props);
        }
      };
    }
    // @ts-expect-error react testing act flag
    globalThis.IS_REACT_ACT_ENVIRONMENT = true;

    MockWebSocket.instances = [];
    // @ts-expect-error mock
    globalThis.WebSocket = MockWebSocket;
    container = document.createElement("div");
    document.body.appendChild(container);
    root = createRoot(container);
  });

  afterEach(() => {
    act(() => {
      root.unmount();
    });
    container.remove();
    globalThis.WebSocket = OriginalWebSocket;
    vi.restoreAllMocks();
  });

  it("does not send input messages when interactive=false", async () => {
    await act(async () => {
      root.render(<Harness interactive={false} />);
    });
    await act(async () => {
      await Promise.resolve();
    });

    const ws = MockWebSocket.instances[0];
    expect(ws).toBeTruthy();

    const target = container.querySelector("canvas")?.parentElement;
    expect(target).toBeTruthy();

    await act(async () => {
      target!.dispatchEvent(
        new PointerEvent("pointerdown", {
          bubbles: true,
          clientX: 200,
          clientY: 200,
          button: 0,
        }),
      );
      target!.dispatchEvent(
        new KeyboardEvent("keydown", {
          bubbles: true,
          key: "a",
          code: "KeyA",
        }),
      );
    });

    expect(ws.sent.filter((m) => m.includes("mouse.") || m.includes("key."))).toEqual([]);
  });
});
