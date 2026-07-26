export type LetterboxRect = {
  offsetX: number;
  offsetY: number;
  drawWidth: number;
  drawHeight: number;
};

export type ViewportSize = {
  viewportWidth: number;
  viewportHeight: number;
};

export type Point = {
  x: number;
  y: number;
};

/**
 * Compute the letterboxed image rect inside a display container
 * (object-fit: contain semantics).
 */
export function computeLetterboxRect(
  containerWidth: number,
  containerHeight: number,
  contentWidth: number,
  contentHeight: number,
): LetterboxRect {
  if (
    containerWidth <= 0 ||
    containerHeight <= 0 ||
    contentWidth <= 0 ||
    contentHeight <= 0
  ) {
    return { offsetX: 0, offsetY: 0, drawWidth: 0, drawHeight: 0 };
  }

  const scale = Math.min(
    containerWidth / contentWidth,
    containerHeight / contentHeight,
  );
  const drawWidth = contentWidth * scale;
  const drawHeight = contentHeight * scale;
  return {
    offsetX: (containerWidth - drawWidth) / 2,
    offsetY: (containerHeight - drawHeight) / 2,
    drawWidth,
    drawHeight,
  };
}

/**
 * Map a pointer position in container CSS pixels to viewport CSS pixels.
 * Returns null when the point falls in letterbox bars.
 */
export function mapPointerToViewport(
  clientX: number,
  clientY: number,
  containerRect: DOMRect | { left: number; top: number; width: number; height: number },
  viewport: ViewportSize,
  contentWidth?: number,
  contentHeight?: number,
): Point | null {
  const vw = viewport.viewportWidth;
  const vh = viewport.viewportHeight;
  if (vw <= 0 || vh <= 0) {
    return null;
  }

  const localX = clientX - containerRect.left;
  const localY = clientY - containerRect.top;
  const cw = contentWidth && contentWidth > 0 ? contentWidth : vw;
  const ch = contentHeight && contentHeight > 0 ? contentHeight : vh;

  const box = computeLetterboxRect(
    containerRect.width,
    containerRect.height,
    cw,
    ch,
  );
  if (box.drawWidth <= 0 || box.drawHeight <= 0) {
    return null;
  }

  const inX = localX - box.offsetX;
  const inY = localY - box.offsetY;
  if (inX < 0 || inY < 0 || inX > box.drawWidth || inY > box.drawHeight) {
    return null;
  }

  return {
    x: (inX / box.drawWidth) * vw,
    y: (inY / box.drawHeight) * vh,
  };
}
