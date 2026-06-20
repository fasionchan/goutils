package examples

import (
	"bytes"
	"testing"
)

func TestRenderCPU(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderCPU(&buf); err != nil {
		t.Fatal(err)
	}
	assertPNG(t, buf.Bytes())
}

func TestRenderMemory(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderMemory(&buf); err != nil {
		t.Fatal(err)
	}
	assertPNG(t, buf.Bytes())
}

func TestRenderFilesystem(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderFilesystem(&buf); err != nil {
		t.Fatal(err)
	}
	assertPNG(t, buf.Bytes())
}

func TestRenderDiskIO(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderDiskIO(&buf); err != nil {
		t.Fatal(err)
	}
	assertPNG(t, buf.Bytes())
}

func TestRenderNetwork(t *testing.T) {
	var buf bytes.Buffer
	if err := RenderNetwork(&buf); err != nil {
		t.Fatal(err)
	}
	assertPNG(t, buf.Bytes())
}

func assertPNG(t *testing.T, data []byte) {
	t.Helper()
	if len(data) == 0 {
		t.Fatal("empty png output")
	}
	if len(data) < 4 || data[0] != 0x89 || data[1] != 0x50 || data[2] != 0x4e || data[3] != 0x47 {
		t.Fatalf("invalid png magic: % x", data[:min(4, len(data))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
