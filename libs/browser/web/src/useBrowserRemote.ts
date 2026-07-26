import { useCallback, useEffect, useRef, useState } from "react";
import {
  encodeRemoteEnvelope,
  mimeFromFormat,
  parseRemoteEnvelope,
  RemoteScreencastFramePayload,
  RemoteScreencastMetaPayload,
  RemoteSessionReadyPayload,
  RemoteType,
} from "./protocol";

export type BrowserRemoteReadyInfo = {
  tabId: string;
  protocolVersion: number;
  meta: RemoteScreencastMetaPayload;
};

export type UseBrowserRemoteOptions = {
  wsUrl: string;
  enabled?: boolean;
  onReady?: (info: BrowserRemoteReadyInfo) => void;
  onError?: (error: Error) => void;
  onClose?: (event: CloseEvent) => void;
  onFrame?: (bitmap: ImageBitmap, meta: RemoteScreencastMetaPayload) => void;
};

export type BrowserRemoteSession = {
  status: "idle" | "connecting" | "ready" | "closed" | "error";
  meta: RemoteScreencastMetaPayload | null;
  sendJson: (type: string, payload?: unknown, id?: string) => void;
  sendRaw: (data: string) => void;
};

export function useBrowserRemote(
  options: UseBrowserRemoteOptions,
): BrowserRemoteSession {
  const { wsUrl, enabled = true, onReady, onError, onClose, onFrame } = options;
  const [status, setStatus] = useState<BrowserRemoteSession["status"]>("idle");
  const [meta, setMeta] = useState<RemoteScreencastMetaPayload | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const metaRef = useRef<RemoteScreencastMetaPayload | null>(null);
  const pendingFrameRef = useRef<RemoteScreencastFramePayload | null>(null);
  const readyInfoRef = useRef<Partial<BrowserRemoteReadyInfo>>({});
  const callbacksRef = useRef({ onReady, onError, onClose, onFrame });
  callbacksRef.current = { onReady, onError, onClose, onFrame };

  const sendRaw = useCallback((data: string) => {
    const ws = wsRef.current;
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(data);
    }
  }, []);

  const sendJson = useCallback(
    (type: string, payload?: unknown, id?: string) => {
      sendRaw(encodeRemoteEnvelope(type, payload, id));
    },
    [sendRaw],
  );

  useEffect(() => {
    if (!enabled || !wsUrl) {
      setStatus("idle");
      return;
    }

    let cancelled = false;
    setStatus("connecting");
    metaRef.current = null;
    pendingFrameRef.current = null;
    readyInfoRef.current = {};
    setMeta(null);

    const ws = new WebSocket(wsUrl);
    ws.binaryType = "arraybuffer";
    wsRef.current = ws;

    const handleFrameBytes = async (bytes: ArrayBuffer, format: string) => {
      const currentMeta = metaRef.current;
      if (!currentMeta) {
        return;
      }
      const blob = new Blob([new Uint8Array(bytes)], {
        type: mimeFromFormat(format),
      });
      try {
        const bitmap = await createImageBitmap(blob);
        if (cancelled) {
          bitmap.close();
          return;
        }
        callbacksRef.current.onFrame?.(bitmap, currentMeta);
      } catch (err) {
        callbacksRef.current.onError?.(
          err instanceof Error ? err : new Error(String(err)),
        );
      }
    };

    ws.onmessage = (event) => {
      if (typeof event.data === "string") {
        try {
          const envelope = parseRemoteEnvelope(event.data);
          switch (envelope.type) {
            case RemoteType.SessionReady: {
              const payload = envelope.payload as RemoteSessionReadyPayload;
              readyInfoRef.current.tabId = payload.tab_id;
              readyInfoRef.current.protocolVersion = payload.protocol_version;
              break;
            }
            case RemoteType.ScreencastMeta: {
              const payload = envelope.payload as RemoteScreencastMetaPayload;
              metaRef.current = payload;
              setMeta(payload);
              readyInfoRef.current.meta = payload;
              setStatus("ready");
              if (
                readyInfoRef.current.tabId != null &&
                readyInfoRef.current.protocolVersion != null
              ) {
                callbacksRef.current.onReady?.(
                  readyInfoRef.current as BrowserRemoteReadyInfo,
                );
              }
              break;
            }
            case RemoteType.ScreencastFrame: {
              pendingFrameRef.current =
                envelope.payload as RemoteScreencastFramePayload;
              break;
            }
            case RemoteType.Error: {
              const payload = envelope.payload as {
                message?: string;
                code?: string;
              };
              callbacksRef.current.onError?.(
                new Error(payload?.message || payload?.code || "remote error"),
              );
              break;
            }
            default:
              break;
          }
        } catch (err) {
          callbacksRef.current.onError?.(
            err instanceof Error ? err : new Error(String(err)),
          );
        }
        return;
      }

      if (event.data instanceof ArrayBuffer) {
        const pending = pendingFrameRef.current;
        pendingFrameRef.current = null;
        const format =
          pending?.format || metaRef.current?.format || "jpeg";
        void handleFrameBytes(event.data, format);
      }
    };

    ws.onerror = () => {
      setStatus("error");
      callbacksRef.current.onError?.(new Error("websocket error"));
    };

    ws.onclose = (ev) => {
      setStatus("closed");
      callbacksRef.current.onClose?.(ev);
    };

    return () => {
      cancelled = true;
      wsRef.current = null;
      ws.close();
    };
  }, [wsUrl, enabled]);

  return { status, meta, sendJson, sendRaw };
}
