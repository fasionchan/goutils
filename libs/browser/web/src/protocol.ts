export const REMOTE_PROTOCOL_VERSION = 1;

export const RemoteType = {
  SessionPing: "session.ping",
  SessionPong: "session.pong",
  SessionReady: "session.ready",
  ScreencastMeta: "screencast.meta",
  ScreencastFrame: "screencast.frame",
  MouseMove: "mouse.move",
  MouseDown: "mouse.down",
  MouseUp: "mouse.up",
  MouseWheel: "mouse.wheel",
  KeyDown: "key.down",
  KeyUp: "key.up",
  KeyPress: "key.press",
  Ack: "ack",
  Error: "error",
} as const;

export type RemoteTypeName = (typeof RemoteType)[keyof typeof RemoteType];

export type RemoteEnvelope<T = unknown> = {
  v: number;
  id?: string;
  type: string;
  ts?: number;
  payload?: T;
};

export type RemoteSessionReadyPayload = {
  tab_id: string;
  protocol_version: number;
};

export type RemoteScreencastMetaPayload = {
  format?: string;
  viewport_width: number;
  viewport_height: number;
  frame_width?: number;
  frame_height?: number;
  device_scale_factor?: number;
};

export type RemoteScreencastFramePayload = {
  seq: number;
  format: string;
  ts?: number;
};

export type RemoteMousePayload = {
  x: number;
  y: number;
  button?: string;
  click_count?: number;
  delta_x?: number;
  delta_y?: number;
  modifiers?: string[];
};

export type RemoteKeyPayload = {
  key?: string;
  code?: string;
  text?: string;
  modifiers?: string[];
  auto_repeat?: boolean;
};

export type RemoteErrorPayload = {
  code: string;
  message: string;
  ref_type?: string;
};

export function parseRemoteEnvelope(raw: string): RemoteEnvelope {
  const value = JSON.parse(raw) as RemoteEnvelope;
  if (!value || typeof value !== "object" || typeof value.type !== "string") {
    throw new Error("invalid remote envelope");
  }
  if (value.v == null) {
    value.v = REMOTE_PROTOCOL_VERSION;
  }
  return value;
}

export function encodeRemoteEnvelope(
  type: string,
  payload?: unknown,
  id?: string,
): string {
  const envelope: RemoteEnvelope = {
    v: REMOTE_PROTOCOL_VERSION,
    type,
    ts: Date.now(),
  };
  if (id) {
    envelope.id = id;
  }
  if (payload !== undefined) {
    envelope.payload = payload;
  }
  return JSON.stringify(envelope);
}

export function mimeFromFormat(format?: string): string {
  switch ((format || "jpeg").toLowerCase()) {
    case "png":
      return "image/png";
    case "jpeg":
    case "jpg":
    default:
      return "image/jpeg";
  }
}

export function modifiersFromEvent(e: {
  altKey: boolean;
  ctrlKey: boolean;
  metaKey: boolean;
  shiftKey: boolean;
}): string[] {
  const mods: string[] = [];
  if (e.altKey) mods.push("alt");
  if (e.ctrlKey) mods.push("ctrl");
  if (e.metaKey) mods.push("meta");
  if (e.shiftKey) mods.push("shift");
  return mods;
}

export function mouseButtonFromEvent(button: number): string {
  switch (button) {
    case 1:
      return "middle";
    case 2:
      return "right";
    case 3:
      return "back";
    case 4:
      return "forward";
    default:
      return "left";
  }
}
