/*
 * Author: fasion
 * Created time: 2026-03-07 00:20:15
 * Last Modified by: fasion
 * Last Modified time: 2026-03-09 00:45:26
 */

package baseutils

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestNewSpillBuffer(t *testing.T) {
	sb := NewSpillBuffer()
	if sb == nil {
		t.Fatal("NewSpillBuffer() returned nil")
	}
	if sb.memory == nil {
		t.Error("memory buffer should not be nil")
	}
	if sb.memoryWaterMark != SpillBufferMemoryWaterMark {
		t.Errorf("memoryWaterMark = %d, want %d", sb.memoryWaterMark, SpillBufferMemoryWaterMark)
	}
	if err := sb.Close(); err != nil {
		t.Errorf("Close() = %v", err)
	}
}

func TestSpillBuffer_WriteRead_SmallData(t *testing.T) {
	sb := NewSpillBuffer()
	defer sb.Close()

	const data = "hello hybrid file"
	n, err := sb.Write([]byte(data))
	if err != nil {
		t.Fatalf("Write() = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() wrote %d bytes, want %d", n, len(data))
	}

	// 未超过水位线，应仍在内存中，rfile 为 nil
	if sb.rfile != nil {
		t.Error("data under water mark should stay in memory, rfile should be nil")
	}

	buf := make([]byte, 64)
	m, err := sb.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read() = %v", err)
	}
	if got := string(buf[:m]); got != data {
		t.Errorf("Read() = %q, want %q", got, data)
	}
}

func TestSpillBuffer_WriteRead_OverWaterMark(t *testing.T) {
	// 使用较小的水位线便于测试
	sb := NewSpillBuffer().WithMemoryWaterMark(128)
	defer sb.Close()

	data := strings.Repeat("x", 256)
	n, err := sb.Write([]byte(data))
	if err != nil {
		t.Fatalf("Write() = %v", err)
	}
	if n != len(data) {
		t.Errorf("Write() wrote %d bytes, want %d", n, len(data))
	}

	// 超过水位线，应已切换到文件
	if sb.rfile == nil {
		t.Error("data over water mark should switch to file, rfile should not be nil")
	}

	buf := make([]byte, 512)
	m, err := sb.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read() = %v", err)
	}
	if got := string(buf[:m]); got != data {
		t.Errorf("Read() = %q, want %q (len %d vs %d)", got, data, m, len(data))
	}
}

func TestSpillBuffer_WriteRead_ExactWaterMark(t *testing.T) {
	waterMark := 64
	sb := NewSpillBuffer().WithMemoryWaterMark(waterMark)
	defer sb.Close()

	// 恰好等于水位线，不应切换（> 才切换）
	data := strings.Repeat("a", waterMark)
	_, err := sb.Write([]byte(data))
	if err != nil {
		t.Fatalf("Write() = %v", err)
	}
	if sb.rfile != nil {
		t.Error("data exactly at water mark should not switch, rfile should be nil")
	}

	// 再写 1 字节则超过，应切换
	_, err = sb.Write([]byte("b"))
	if err != nil {
		t.Fatalf("Write() = %v", err)
	}
	if sb.rfile == nil {
		t.Error("after exceeding water mark, rfile should not be nil")
	}

	all, err := io.ReadAll(sb)
	if err != nil {
		t.Fatalf("ReadAll() = %v", err)
	}
	want := data + "b"
	if string(all) != want {
		t.Errorf("ReadAll() = %q, want %q", string(all), want)
	}
}

func TestSpillBuffer_Read_Empty(t *testing.T) {
	sb := NewSpillBuffer()
	defer sb.Close()

	buf := make([]byte, 8)
	n, err := sb.Read(buf)
	if err != io.EOF {
		t.Errorf("Read() err = %v, want io.EOF", err)
	}
	if n != 0 {
		t.Errorf("Read() = %d bytes, want 0", n)
	}
}

func TestSpillBuffer_WriteThenRead_PartialReads(t *testing.T) {
	sb := NewSpillBuffer()
	defer sb.Close()

	const data = "abcdefghij"
	_, err := sb.Write([]byte(data))
	if err != nil {
		t.Fatalf("Write() = %v", err)
	}

	var got strings.Builder
	buf := make([]byte, 3)
	for {
		n, err := sb.Read(buf)
		if n > 0 {
			got.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() = %v", err)
		}
	}
	if got.String() != data {
		t.Errorf("Read() concatenated = %q, want %q", got.String(), data)
	}
}

func TestSpillBuffer_WithMemoryWaterMark(t *testing.T) {
	waterMark := 32
	sb := NewSpillBuffer().WithMemoryWaterMark(waterMark)
	defer sb.Close()

	if sb.memoryWaterMark != waterMark {
		t.Errorf("memoryWaterMark = %d, want %d", sb.memoryWaterMark, waterMark)
	}

	// 小数据不切换
	_, _ = sb.Write([]byte("short"))
	if sb.rfile != nil {
		t.Error("short write should not switch")
	}

	// 超过 32 字节切换
	_, _ = sb.Write([]byte(strings.Repeat("x", 40)))
	if sb.rfile == nil {
		t.Error("over water mark should switch to file")
	}
}

func TestSpillBuffer_WithMutex(t *testing.T) {
	sb := NewSpillBuffer().WithMutex().WithMemoryWaterMark(256)
	defer sb.Close()

	if sb.mutex == nil {
		t.Fatal("WithMutex() should set mutex")
	}

	// 并发写
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = sb.Write([]byte("concurrent write "))
		}()
	}
	wg.Wait()

	// 能正常读完即可
	_, err := io.ReadAll(sb)
	if err != nil {
		t.Errorf("ReadAll after concurrent write: %v", err)
	}
}

func TestSpillBuffer_Close_Idempotent(t *testing.T) {
	sb := NewSpillBuffer()
	_, _ = sb.Write([]byte("data"))

	if err := sb.Close(); err != nil {
		t.Errorf("first Close() = %v", err)
	}
	// 再次 Close 不应 panic，行为未定义但至少不崩溃
	sb.Close()
}

func TestSpillBuffer_Close_AfterSwitchToFile(t *testing.T) {
	sb := NewSpillBuffer().WithMemoryWaterMark(8)
	_, _ = sb.Write([]byte("enough to switch"))

	if err := sb.Close(); err != nil {
		t.Errorf("Close() after switch = %v", err)
	}
	// 临时文件应已被删除，无法再读
	if sb.file != nil {
		t.Error("file handle should be nil after Close")
	}
}

func TestSpillBuffer_Close_RemovesTempFileFromDisk(t *testing.T) {
	sb := NewSpillBuffer().WithMemoryWaterMark(32)
	_, _ = sb.Write([]byte(strings.Repeat("x", 64)))

	// 已溢出到文件，拿到临时文件路径（同包可访问 sb.file）
	if sb.file == nil {
		t.Fatal("expected spill to file, sb.file should not be nil")
	}
	path := sb.file.Name()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("temp file should exist before Close: %v", err)
	}

	if err := sb.Close(); err != nil {
		t.Fatalf("Close() = %v", err)
	}

	if _, err := os.Stat(path); err == nil {
		t.Error("temp file should be removed from disk after Close")
	} else if !os.IsNotExist(err) {
		t.Errorf("expected os.ErrNotExist after Close, got: %v", err)
	}
}

func TestSpillBuffer_WriteAfterSwitch(t *testing.T) {
	sb := NewSpillBuffer().WithMemoryWaterMark(32)
	defer sb.Close()

	// 触发切换
	_, _ = sb.Write([]byte(strings.Repeat("a", 64)))

	// 切换后继续写
	extra := " appended"
	n, err := sb.Write([]byte(extra))
	if err != nil {
		t.Fatalf("Write() after switch = %v", err)
	}
	if n != len(extra) {
		t.Errorf("Write() after switch wrote %d, want %d", n, len(extra))
	}

	all, err := io.ReadAll(sb)
	if err != nil {
		t.Fatalf("ReadAll() = %v", err)
	}
	if !strings.HasSuffix(string(all), extra) {
		t.Errorf("ReadAll() = %q, should end with %q", string(all), extra)
	}
}
