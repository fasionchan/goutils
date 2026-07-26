import { describe, expect, it } from "vitest";
import {
  encodeRemoteEnvelope,
  mimeFromFormat,
  parseRemoteEnvelope,
  RemoteType,
} from "./protocol";

describe("protocol helpers", () => {
  it("round-trips envelope JSON", () => {
    const raw = encodeRemoteEnvelope(RemoteType.MouseDown, {
      x: 10,
      y: 20,
      button: "left",
    }, "id-1");
    const env = parseRemoteEnvelope(raw);
    expect(env.type).toBe(RemoteType.MouseDown);
    expect(env.id).toBe("id-1");
    expect((env.payload as { x: number }).x).toBe(10);
  });

  it("maps format to mime", () => {
    expect(mimeFromFormat("png")).toBe("image/png");
    expect(mimeFromFormat("jpeg")).toBe("image/jpeg");
  });
});
