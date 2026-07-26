import { describe, expect, it } from "vitest";
import { computeLetterboxRect, mapPointerToViewport } from "./coords";

describe("computeLetterboxRect", () => {
  it("letterboxes a wide viewport inside a square container", () => {
    const box = computeLetterboxRect(400, 400, 1280, 720);
    expect(box.drawWidth).toBeCloseTo(400);
    expect(box.drawHeight).toBeCloseTo(225);
    expect(box.offsetX).toBeCloseTo(0);
    expect(box.offsetY).toBeCloseTo(87.5);
  });
});

describe("mapPointerToViewport", () => {
  it("maps the visible image center to viewport center", () => {
    const container = { left: 0, top: 0, width: 400, height: 400 };
    const point = mapPointerToViewport(
      200,
      200,
      container,
      { viewportWidth: 1280, viewportHeight: 720 },
    );
    expect(point).not.toBeNull();
    expect(point!.x).toBeCloseTo(640, 0);
    expect(point!.y).toBeCloseTo(360, 0);
  });

  it("returns null for letterbox bars", () => {
    const container = { left: 0, top: 0, width: 400, height: 400 };
    const point = mapPointerToViewport(
      200,
      10,
      container,
      { viewportWidth: 1280, viewportHeight: 720 },
    );
    expect(point).toBeNull();
  });
});
