/*
 * Author: fasion
 * Created time: 2026-01-25 23:37:47
 * Last Modified by: fasion
 * Last Modified time: 2026-01-25 23:37:55
 */

package types

import (
	"bufio"
	"testing"
)

// TestBytesBoundedBuffer_Basic 测试基本功能
func TestBytesBoundedBuffer_Basic(t *testing.T) {
	buf := NewBytesBoundedBuffer(10)

	// 初始状态
	if buf.IsFull() {
		t.Error("初始状态应该未满")
	}
	if len(buf.Datas()) != 0 {
		t.Errorf("初始数据应该为空，实际: %v", buf.Datas())
	}

	// 写入数据
	data := []byte("hello")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != len(data) {
		t.Errorf("写入字节数不匹配，期望: %d, 实际: %d", len(data), n)
	}
	if string(buf.Datas()) != "hello" {
		t.Errorf("数据不匹配，期望: hello, 实际: %s", string(buf.Datas()))
	}
	if buf.IsFull() {
		t.Error("写入5字节后，10字节缓冲区应该未满")
	}
}

// TestBytesBoundedBuffer_Full 测试缓冲区满的情况
func TestBytesBoundedBuffer_Full(t *testing.T) {
	buf := NewBytesBoundedBuffer(5)

	// 写满缓冲区
	data1 := []byte("hello")
	n, err := buf.Write(data1)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != 5 {
		t.Errorf("写入字节数不匹配，期望: 5, 实际: %d", n)
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}

	// 尝试继续写入，应该失败
	data2 := []byte("world")
	n, err = buf.Write(data2)
	if err != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}
	if string(buf.Datas()) != "hello" {
		t.Errorf("数据不应该改变，期望: hello, 实际: %s", string(buf.Datas()))
	}
}

// TestBytesBoundedBuffer_PartialWrite 测试部分写入
func TestBytesBoundedBuffer_PartialWrite(t *testing.T) {
	buf := NewBytesBoundedBuffer(5)

	// 先写入部分数据
	data1 := []byte("hi")
	n, err := buf.Write(data1)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != 2 {
		t.Errorf("写入字节数不匹配，期望: 2, 实际: %d", n)
	}

	// 尝试写入超过剩余空间的数据
	data2 := []byte("world")
	n, err = buf.Write(data2)
	if err != nil {
		t.Errorf("部分写入不应该返回错误，实际: %v", err)
	}
	if n != 3 {
		t.Errorf("应该只写入3字节，实际: %d", n)
	}
	if string(buf.Datas()) != "hiwor" {
		t.Errorf("数据不匹配，期望: hiwor, 实际: %s", string(buf.Datas()))
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}
}

// TestBytesBoundedBuffer_EmptyWrite 测试写入空数据
func TestBytesBoundedBuffer_EmptyWrite(t *testing.T) {
	buf := NewBytesBoundedBuffer(5)

	n, err := buf.Write(nil)
	if err != nil {
		t.Errorf("写入空数据不应该返回错误，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}

	n, err = buf.Write([]byte{})
	if err != nil {
		t.Errorf("写入空数据不应该返回错误，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}
}

// TestBytesBoundedBuffer_ZeroSize 测试零大小缓冲区
func TestBytesBoundedBuffer_ZeroSize(t *testing.T) {
	buf := NewBytesBoundedBuffer(0)

	if !buf.IsFull() {
		t.Error("零大小缓冲区应该始终为满")
	}

	data := []byte("hello")
	n, err := buf.Write(data)
	if err != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}
}

// TestBytesRingBuffer_Basic 测试基本功能
func TestBytesRingBuffer_Basic(t *testing.T) {
	buf := NewBytesRingBuffer(10)

	// 初始状态
	if buf.IsFull() {
		t.Error("初始状态应该未满")
	}
	if len(buf.Datas()) != 0 {
		t.Errorf("初始数据应该为空，实际: %v", buf.Datas())
	}
	if buf.TotalWritten() != 0 {
		t.Errorf("初始总写入数应该为0，实际: %d", buf.TotalWritten())
	}
	if buf.IsTruncated() {
		t.Error("初始状态不应该被截断")
	}

	// 写入数据
	data := []byte("hello")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != len(data) {
		t.Errorf("写入字节数不匹配，期望: %d, 实际: %d", len(data), n)
	}
	if string(buf.Datas()) != "hello" {
		t.Errorf("数据不匹配，期望: hello, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != int64(len(data)) {
		t.Errorf("总写入数不匹配，期望: %d, 实际: %d", len(data), buf.TotalWritten())
	}
}

// TestBytesRingBuffer_Overwrite 测试循环覆盖
func TestBytesRingBuffer_Overwrite(t *testing.T) {
	buf := NewBytesRingBuffer(5)

	// 先写满缓冲区
	data1 := []byte("hello")
	n, err := buf.Write(data1)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != 5 {
		t.Errorf("写入字节数不匹配，期望: 5, 实际: %d", n)
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}
	if buf.IsTruncated() {
		t.Error("写入5字节到5字节缓冲区不应该被截断")
	}

	// 继续写入，应该覆盖旧数据
	data2 := []byte("world")
	n, err = buf.Write(data2)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != 5 {
		t.Errorf("写入字节数不匹配，期望: 5, 实际: %d", n)
	}
	if string(buf.Datas()) != "world" {
		t.Errorf("数据应该被覆盖，期望: world, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != 10 {
		t.Errorf("总写入数不匹配，期望: 10, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}
}

// TestBytesRingBuffer_PartialOverwrite 测试部分覆盖
func TestBytesRingBuffer_PartialOverwrite(t *testing.T) {
	buf := NewBytesRingBuffer(5)

	// 先写满缓冲区
	buf.Write([]byte("hello"))

	// 写入3字节，应该覆盖前3字节
	// writePos 会变成 3，所以返回的是 buffer[3:] + buffer[:3] = "lo" + "abc" = "loabc"
	buf.Write([]byte("abc"))
	result := buf.Datas()
	if string(result) != "loabc" {
		t.Errorf("部分覆盖结果不匹配，期望: loabc, 实际: %s", string(result))
	}
}

// TestBytesRingBuffer_LargeWrite 测试写入大于缓冲区大小的数据
func TestBytesRingBuffer_LargeWrite(t *testing.T) {
	buf := NewBytesRingBuffer(5)

	// 写入大于缓冲区大小的数据
	data := []byte("helloworld")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != len(data) {
		t.Errorf("写入字节数不匹配，期望: %d, 实际: %d", len(data), n)
	}
	if string(buf.Datas()) != "world" {
		t.Errorf("应该保留最后5字节，期望: world, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != int64(len(data)) {
		t.Errorf("总写入数不匹配，期望: %d, 实际: %d", len(data), buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}
}

// TestBytesRingBuffer_MultipleWrites 测试多次写入
func TestBytesRingBuffer_MultipleWrites(t *testing.T) {
	buf := NewBytesRingBuffer(5)

	// 多次写入
	buf.Write([]byte("ab"))
	buf.Write([]byte("cd"))
	buf.Write([]byte("e"))

	if string(buf.Datas()) != "abcde" {
		t.Errorf("多次写入结果不匹配，期望: abcde, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != 5 {
		t.Errorf("总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}

	// 继续写入，应该覆盖
	// writePos 会变成 1，所以返回的是 buffer[1:] + buffer[:1] = "bcde" + "f" = "bcdef"
	buf.Write([]byte("f"))
	if string(buf.Datas()) != "bcdef" {
		t.Errorf("覆盖后结果不匹配，期望: bcdef, 实际: %s", string(buf.Datas()))
	}
}

// TestBytesRingBuffer_EmptyWrite 测试写入空数据
func TestBytesRingBuffer_EmptyWrite(t *testing.T) {
	buf := NewBytesRingBuffer(5)

	n, err := buf.Write(nil)
	if err != nil {
		t.Errorf("写入空数据不应该返回错误，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}

	n, err = buf.Write([]byte{})
	if err != nil {
		t.Errorf("写入空数据不应该返回错误，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}
}

// TestBytesRingBuffer_ZeroSize 测试零大小缓冲区
func TestBytesRingBuffer_ZeroSize(t *testing.T) {
	buf := NewBytesRingBuffer(0)

	data := []byte("hello")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("零大小缓冲区写入不应该返回错误，实际: %v", err)
	}
	if n != len(data) {
		t.Errorf("写入字节数不匹配，期望: %d, 实际: %d", len(data), n)
	}
	if buf.Datas() != nil {
		t.Errorf("零大小缓冲区应该返回 nil，实际: %v", buf.Datas())
	}
	// 零大小缓冲区 size 会被设置为 0，所以 totalWritten 不会增加
	// 根据代码逻辑，size <= 0 时直接返回 len(datas)，不会更新 totalWritten
	if buf.TotalWritten() != 0 {
		t.Errorf("零大小缓冲区总写入数应该为0，实际: %d", buf.TotalWritten())
	}
}

// TestBytesRingBuffer_Reset 测试重置功能
func TestBytesRingBuffer_Reset(t *testing.T) {
	buf := NewBytesRingBuffer(5)

	// 写入数据
	buf.Write([]byte("hello"))
	if buf.TotalWritten() != 5 {
		t.Errorf("总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}

	// 重置
	buf.Reset()
	if buf.TotalWritten() != 0 {
		t.Errorf("重置后总写入数应该为0，实际: %d", buf.TotalWritten())
	}
	if len(buf.Datas()) != 0 {
		t.Errorf("重置后数据应该为空，实际: %v", buf.Datas())
	}
	if buf.IsTruncated() {
		t.Error("重置后不应该被标记为截断")
	}

	// 重置后可以继续使用
	buf.Write([]byte("world"))
	if string(buf.Datas()) != "world" {
		t.Errorf("重置后写入结果不匹配，期望: world, 实际: %s", string(buf.Datas()))
	}
}

// TestBytesRingBuffer_CircularBehavior 测试循环行为
func TestBytesRingBuffer_CircularBehavior(t *testing.T) {
	buf := NewBytesRingBuffer(3)

	// 写入数据，使其循环
	buf.Write([]byte("abc"))
	// writePos = 0 (3 % 3 = 0)
	buf.Write([]byte("d"))
	// writePos = 1，返回 buffer[1:] + buffer[:1] = "bc" + "d" = "bcd"
	if string(buf.Datas()) != "bcd" {
		t.Errorf("循环写入结果不匹配，期望: bcd, 实际: %s", string(buf.Datas()))
	}

	buf.Write([]byte("ef"))
	// writePos = 0 ((1 + 2) % 3 = 0)，返回 buffer[0:] + buffer[:0] = "def"
	if string(buf.Datas()) != "def" {
		t.Errorf("循环写入结果不匹配，期望: def, 实际: %s", string(buf.Datas()))
	}
}

// TestBytesTruncatedBuffer_Basic 测试基本功能
func TestBytesTruncatedBuffer_Basic(t *testing.T) {
	buf := NewBytesTruncatedBuffer(10)

	// 初始状态
	if buf.IsFull() {
		t.Error("初始状态应该未满")
	}
	if len(buf.Datas()) != 0 {
		t.Errorf("初始数据应该为空，实际: %v", buf.Datas())
	}
	if buf.TotalWritten() != 0 {
		t.Errorf("初始总写入数应该为0，实际: %d", buf.TotalWritten())
	}
	if buf.IsTruncated() {
		t.Error("初始状态不应该被截断")
	}

	// 写入数据
	data := []byte("hello")
	n, err := buf.Write(data)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != int64(len(data)) {
		t.Errorf("写入字节数不匹配，期望: %d, 实际: %d", len(data), n)
	}
	if string(buf.Datas()) != "hello" {
		t.Errorf("数据不匹配，期望: hello, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != int64(len(data)) {
		t.Errorf("总写入数不匹配，期望: %d, 实际: %d", len(data), buf.TotalWritten())
	}
	if buf.IsTruncated() {
		t.Error("写入5字节到10字节缓冲区不应该被截断")
	}
	if buf.IsFull() {
		t.Error("写入5字节后，10字节缓冲区应该未满")
	}
}

// TestBytesTruncatedBuffer_Full 测试缓冲区满的情况
func TestBytesTruncatedBuffer_Full(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 写满缓冲区
	data1 := []byte("hello")
	n, err := buf.Write(data1)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != 5 {
		t.Errorf("写入字节数不匹配，期望: 5, 实际: %d", n)
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}
	if buf.TotalWritten() != 5 {
		t.Errorf("总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}
	if buf.IsTruncated() {
		t.Error("写入5字节到5字节缓冲区不应该被截断")
	}

	// 继续写入，应该被截断
	data2 := []byte("world")
	n, err = buf.Write(data2)
	if err != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err)
	}
	if n != 5 {
		t.Errorf("应该记录写入5字节，实际: %d", n)
	}
	if string(buf.Datas()) != "hello" {
		t.Errorf("缓冲区数据不应该改变，期望: hello, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != 10 {
		t.Errorf("总写入数不匹配，期望: 10, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}
}

// TestBytesTruncatedBuffer_Truncated 测试截断功能
func TestBytesTruncatedBuffer_Truncated(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 先写入部分数据
	data1 := []byte("hi")
	n, err := buf.Write(data1)
	if err != nil {
		t.Errorf("写入失败: %v", err)
	}
	if n != 2 {
		t.Errorf("写入字节数不匹配，期望: 2, 实际: %d", n)
	}

	// 写入超过剩余空间的数据，应该被截断
	data2 := []byte("world")
	n, err = buf.Write(data2)
	if err != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err)
	}
	if n != 5 {
		t.Errorf("应该记录写入5字节，实际: %d", n)
	}
	// 缓冲区应该只包含前3字节（hi + wor）
	if string(buf.Datas()) != "hiwor" {
		t.Errorf("数据不匹配，期望: hiwor, 实际: %s", string(buf.Datas()))
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}
	if buf.TotalWritten() != 7 {
		t.Errorf("总写入数不匹配，期望: 7, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}
}

// TestBytesTruncatedBuffer_LargeWrite 测试写入大于缓冲区大小的数据
func TestBytesTruncatedBuffer_LargeWrite(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 写入大于缓冲区大小的数据
	data := []byte("helloworld")
	n, err := buf.Write(data)
	if err != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err)
	}
	if n != int64(len(data)) {
		t.Errorf("写入字节数不匹配，期望: %d, 实际: %d", len(data), n)
	}
	// 缓冲区应该只包含前5字节
	if string(buf.Datas()) != "hello" {
		t.Errorf("应该保留前5字节，期望: hello, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != int64(len(data)) {
		t.Errorf("总写入数不匹配，期望: %d, 实际: %d", len(data), buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}
}

// TestBytesTruncatedBuffer_MultipleWrites 测试多次写入
func TestBytesTruncatedBuffer_MultipleWrites(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 多次写入
	n1, err1 := buf.Write([]byte("ab"))
	if err1 != nil {
		t.Errorf("第一次写入失败: %v", err1)
	}
	if n1 != 2 {
		t.Errorf("第一次写入字节数不匹配，期望: 2, 实际: %d", n1)
	}

	n2, err2 := buf.Write([]byte("cd"))
	if err2 != nil {
		t.Errorf("第二次写入失败: %v", err2)
	}
	if n2 != 2 {
		t.Errorf("第二次写入字节数不匹配，期望: 2, 实际: %d", n2)
	}

	n3, err3 := buf.Write([]byte("e"))
	if err3 != nil {
		t.Errorf("第三次写入失败: %v", err3)
	}
	if n3 != 1 {
		t.Errorf("第三次写入字节数不匹配，期望: 1, 实际: %d", n3)
	}

	if string(buf.Datas()) != "abcde" {
		t.Errorf("多次写入结果不匹配，期望: abcde, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != 5 {
		t.Errorf("总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}
	if buf.IsTruncated() {
		t.Error("写入5字节到5字节缓冲区不应该被截断")
	}

	// 继续写入，应该被截断
	n4, err4 := buf.Write([]byte("fgh"))
	if err4 != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err4)
	}
	if n4 != 3 {
		t.Errorf("应该记录写入3字节，实际: %d", n4)
	}
	if string(buf.Datas()) != "abcde" {
		t.Errorf("缓冲区数据不应该改变，期望: abcde, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != 8 {
		t.Errorf("总写入数不匹配，期望: 8, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}
}

// TestBytesTruncatedBuffer_EmptyWrite 测试写入空数据
func TestBytesTruncatedBuffer_EmptyWrite(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	n, err := buf.Write(nil)
	if err != nil {
		t.Errorf("写入空数据不应该返回错误，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}
	if buf.TotalWritten() != 0 {
		t.Errorf("总写入数应该为0，实际: %d", buf.TotalWritten())
	}

	n, err = buf.Write([]byte{})
	if err != nil {
		t.Errorf("写入空数据不应该返回错误，实际: %v", err)
	}
	if n != 0 {
		t.Errorf("应该写入0字节，实际: %d", n)
	}
	if buf.TotalWritten() != 0 {
		t.Errorf("总写入数应该为0，实际: %d", buf.TotalWritten())
	}
}

// TestBytesTruncatedBuffer_ZeroSize 测试零大小缓冲区
func TestBytesTruncatedBuffer_ZeroSize(t *testing.T) {
	buf := NewBytesTruncatedBuffer(0)

	if !buf.IsFull() {
		t.Error("零大小缓冲区应该始终为满")
	}

	data := []byte("hello")
	n, err := buf.Write(data)
	if err != bufio.ErrBufferFull {
		t.Errorf("应该返回 ErrBufferFull，实际: %v", err)
	}
	if n != int64(len(data)) {
		t.Errorf("应该记录写入所有字节，期望: %d, 实际: %d", len(data), n)
	}
	if buf.Datas() != nil {
		t.Errorf("零大小缓冲区应该返回 nil，实际: %v", buf.Datas())
	}
	if buf.TotalWritten() != int64(len(data)) {
		t.Errorf("总写入数不匹配，期望: %d, 实际: %d", len(data), buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("零大小缓冲区写入应该被标记为截断")
	}
}

// TestBytesTruncatedBuffer_Reset 测试重置功能
func TestBytesTruncatedBuffer_Reset(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 写入数据
	buf.Write([]byte("hello"))
	if buf.TotalWritten() != 5 {
		t.Errorf("总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsFull() {
		t.Error("缓冲区应该已满")
	}

	// 继续写入使其截断
	buf.Write([]byte("world"))
	if buf.TotalWritten() != 10 {
		t.Errorf("总写入数不匹配，期望: 10, 实际: %d", buf.TotalWritten())
	}
	if !buf.IsTruncated() {
		t.Error("应该被标记为截断")
	}

	// 重置
	buf.Reset()
	if buf.TotalWritten() != 0 {
		t.Errorf("重置后总写入数应该为0，实际: %d", buf.TotalWritten())
	}
	if len(buf.Datas()) != 0 {
		t.Errorf("重置后数据应该为空，实际: %v", buf.Datas())
	}
	if buf.IsTruncated() {
		t.Error("重置后不应该被标记为截断")
	}
	if buf.IsFull() {
		t.Error("重置后缓冲区应该未满")
	}

	// 重置后可以继续使用
	buf.Write([]byte("world"))
	if string(buf.Datas()) != "world" {
		t.Errorf("重置后写入结果不匹配，期望: world, 实际: %s", string(buf.Datas()))
	}
	if buf.TotalWritten() != 5 {
		t.Errorf("重置后总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}
}

// TestBytesTruncatedBuffer_IsTruncated 测试截断状态判断
func TestBytesTruncatedBuffer_IsTruncated(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 初始状态不应该被截断
	if buf.IsTruncated() {
		t.Error("初始状态不应该被截断")
	}

	// 写入少于缓冲区大小的数据，不应该被截断
	buf.Write([]byte("hi"))
	if buf.IsTruncated() {
		t.Error("写入2字节到5字节缓冲区不应该被截断")
	}

	// 写满缓冲区，不应该被截断
	buf.Write([]byte("llo"))
	if buf.IsTruncated() {
		t.Error("写入5字节到5字节缓冲区不应该被截断")
	}

	// 继续写入，应该被截断
	buf.Write([]byte("w"))
	if !buf.IsTruncated() {
		t.Error("写入超过缓冲区大小应该被标记为截断")
	}

	// 重置后不应该被截断
	buf.Reset()
	if buf.IsTruncated() {
		t.Error("重置后不应该被标记为截断")
	}
}

// TestBytesTruncatedBuffer_TotalWritten 测试总写入数统计
func TestBytesTruncatedBuffer_TotalWritten(t *testing.T) {
	buf := NewBytesTruncatedBuffer(5)

	// 初始总写入数应该为0
	if buf.TotalWritten() != 0 {
		t.Errorf("初始总写入数应该为0，实际: %d", buf.TotalWritten())
	}

	// 写入数据
	buf.Write([]byte("hi"))
	if buf.TotalWritten() != 2 {
		t.Errorf("总写入数不匹配，期望: 2, 实际: %d", buf.TotalWritten())
	}

	// 继续写入
	buf.Write([]byte("llo"))
	if buf.TotalWritten() != 5 {
		t.Errorf("总写入数不匹配，期望: 5, 实际: %d", buf.TotalWritten())
	}

	// 写入超过缓冲区大小的数据，总写入数应该包括所有尝试写入的字节
	buf.Write([]byte("world"))
	if buf.TotalWritten() != 10 {
		t.Errorf("总写入数不匹配，期望: 10, 实际: %d", buf.TotalWritten())
	}

	// 重置后总写入数应该为0
	buf.Reset()
	if buf.TotalWritten() != 0 {
		t.Errorf("重置后总写入数应该为0，实际: %d", buf.TotalWritten())
	}
}
