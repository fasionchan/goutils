import {
  type CSSProperties,
  type KeyboardEvent,
  type PointerEvent,
  type WheelEvent,
  useCallback,
  useEffect,
  useRef,
} from "react";
import { mapPointerToViewport } from "./coords";
import {
  mouseButtonFromEvent,
  modifiersFromEvent,
  RemoteType,
} from "./protocol";
import {
  type BrowserRemoteReadyInfo,
  useBrowserRemote,
} from "./useBrowserRemote";

export type BrowserRemoteViewerProps = {
  wsUrl: string;
  interactive?: boolean;
  className?: string;
  style?: CSSProperties;
  onReady?: (info: BrowserRemoteReadyInfo) => void;
  onError?: (error: Error) => void;
  onClose?: (event: CloseEvent) => void;
};

export function BrowserRemoteViewer(props: BrowserRemoteViewerProps) {
  const {
    wsUrl,
    interactive = true,
    className,
    style,
    onReady,
    onError,
    onClose,
  } = props;

  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const moveRafRef = useRef<number | null>(null);
  const pendingMoveRef = useRef<{ x: number; y: number; modifiers: string[] } | null>(
    null,
  );

  const { meta, sendJson } = useBrowserRemote({
    wsUrl,
    onReady,
    onError,
    onClose,
    onFrame: (bitmap) => {
      const canvas = canvasRef.current;
      if (!canvas) {
        bitmap.close();
        return;
      }
      if (canvas.width !== bitmap.width || canvas.height !== bitmap.height) {
        canvas.width = bitmap.width;
        canvas.height = bitmap.height;
      }
      const ctx = canvas.getContext("2d");
      if (!ctx) {
        bitmap.close();
        return;
      }
      ctx.drawImage(bitmap, 0, 0);
      bitmap.close();
    },
  });

  const resolveViewportPoint = useCallback(
    (clientX: number, clientY: number) => {
      const container = containerRef.current;
      if (!container || !meta) {
        return null;
      }
      return mapPointerToViewport(
        clientX,
        clientY,
        container.getBoundingClientRect(),
        {
          viewportWidth: meta.viewport_width,
          viewportHeight: meta.viewport_height,
        },
        meta.frame_width || meta.viewport_width,
        meta.frame_height || meta.viewport_height,
      );
    },
    [meta],
  );

  const flushMove = useCallback(() => {
    moveRafRef.current = null;
    const pending = pendingMoveRef.current;
    pendingMoveRef.current = null;
    if (!pending || !interactive) {
      return;
    }
    sendJson(RemoteType.MouseMove, {
      x: pending.x,
      y: pending.y,
      modifiers: pending.modifiers,
    });
  }, [interactive, sendJson]);

  const onPointerMove = useCallback(
    (e: PointerEvent) => {
      if (!interactive) return;
      const point = resolveViewportPoint(e.clientX, e.clientY);
      if (!point) return;
      pendingMoveRef.current = {
        x: point.x,
        y: point.y,
        modifiers: modifiersFromEvent(e),
      };
      if (moveRafRef.current == null) {
        moveRafRef.current = requestAnimationFrame(flushMove);
      }
    },
    [flushMove, interactive, resolveViewportPoint],
  );

  const onPointerDown = useCallback(
    (e: PointerEvent) => {
      if (!interactive) return;
      (e.currentTarget as HTMLElement).setPointerCapture?.(e.pointerId);
      const point = resolveViewportPoint(e.clientX, e.clientY);
      if (!point) return;
      sendJson(RemoteType.MouseDown, {
        x: point.x,
        y: point.y,
        button: mouseButtonFromEvent(e.button),
        click_count: e.detail || 1,
        modifiers: modifiersFromEvent(e),
      });
    },
    [interactive, resolveViewportPoint, sendJson],
  );

  const onPointerUp = useCallback(
    (e: PointerEvent) => {
      if (!interactive) return;
      const point = resolveViewportPoint(e.clientX, e.clientY);
      if (!point) return;
      sendJson(RemoteType.MouseUp, {
        x: point.x,
        y: point.y,
        button: mouseButtonFromEvent(e.button),
        click_count: e.detail || 1,
        modifiers: modifiersFromEvent(e),
      });
    },
    [interactive, resolveViewportPoint, sendJson],
  );

  const onWheel = useCallback(
    (e: WheelEvent) => {
      if (!interactive) return;
      e.preventDefault();
      const point = resolveViewportPoint(e.clientX, e.clientY);
      if (!point) return;
      sendJson(RemoteType.MouseWheel, {
        x: point.x,
        y: point.y,
        delta_x: e.deltaX,
        delta_y: e.deltaY,
        modifiers: modifiersFromEvent(e),
      });
    },
    [interactive, resolveViewportPoint, sendJson],
  );

  const onKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (!interactive) return;
      e.preventDefault();
      sendJson(RemoteType.KeyDown, {
        key: e.key,
        code: e.code,
        text: e.key.length === 1 ? e.key : "",
        modifiers: modifiersFromEvent(e),
        auto_repeat: e.repeat,
      });
    },
    [interactive, sendJson],
  );

  const onKeyUp = useCallback(
    (e: KeyboardEvent) => {
      if (!interactive) return;
      e.preventDefault();
      sendJson(RemoteType.KeyUp, {
        key: e.key,
        code: e.code,
        modifiers: modifiersFromEvent(e),
      });
    },
    [interactive, sendJson],
  );

  useEffect(() => {
    return () => {
      if (moveRafRef.current != null) {
        cancelAnimationFrame(moveRafRef.current);
      }
    };
  }, []);

  return (
    <div
      ref={containerRef}
      className={className}
      tabIndex={interactive ? 0 : undefined}
      onPointerMove={onPointerMove}
      onPointerDown={onPointerDown}
      onPointerUp={onPointerUp}
      onWheel={onWheel}
      onKeyDown={onKeyDown}
      onKeyUp={onKeyUp}
      style={{
        position: "relative",
        width: "100%",
        height: "100%",
        minHeight: 240,
        background: "#111",
        outline: "none",
        overflow: "hidden",
        touchAction: "none",
        ...style,
      }}
    >
      <canvas
        ref={canvasRef}
        style={{
          position: "absolute",
          inset: 0,
          width: "100%",
          height: "100%",
          objectFit: "contain",
          pointerEvents: "none",
        }}
      />
    </div>
  );
}
