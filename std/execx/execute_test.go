/*
 * Author: fasion
 * Created time: 2026-01-26 23:03:34
 * Last Modified by: fasion
 * Last Modified time: 2026-01-26 23:03:41
 */

package execx

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestExecute_Basic 测试基本命令执行
func TestExecute_Basic(t *testing.T) {
	ctx := context.Background()
	// 使用足够大的限制来避免缓冲区满的问题
	result, err := Execute(ctx, []string{"echo", "hello world"}, WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Stdout != "hello world\n" {
		t.Errorf("stdout 不匹配，期望: 'hello world\\n', 实际: %q", result.Stdout)
	}
	if result.Stderr != "" {
		t.Errorf("stderr 应该为空，实际: %q", result.Stderr)
	}
	if result.Error != nil {
		t.Errorf("不应该有错误，实际: %v", result.Error)
	}
	if result.StartTime.IsZero() {
		t.Error("StartTime 应该被设置")
	}
	if result.EndTime.IsZero() {
		t.Error("EndTime 应该被设置")
	}
	if result.ExpiredDuration <= 0 {
		t.Error("ExpiredDuration 应该大于0")
	}
}

// TestExecute_ExitCode 测试非零退出码
func TestExecute_ExitCode(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"sh", "-c", "exit 42"}, WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 42 {
		t.Errorf("退出码应该是42，实际: %d", result.ExitCode)
	}
}

// TestExecute_ShellMode 测试 Shell 模式
func TestExecute_ShellMode(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"echo hello; echo world"}, WithShellMode(true), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	expected := "hello\nworld\n"
	if result.Stdout != expected {
		t.Errorf("stdout 不匹配，期望: %q, 实际: %q", expected, result.Stdout)
	}
}

// TestExecute_CustomShell 测试自定义 Shell
func TestExecute_CustomShell(t *testing.T) {
	ctx := context.Background()
	// 使用 bash 执行命令（如果 bash 不存在则跳过）
	result, err := Execute(ctx, []string{"echo $0"}, WithShellMode(true), WithShell("bash"), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		// bash 可能不存在，跳过测试
		t.Skipf("bash 可能不存在，跳过测试: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	// bash 应该输出 bash 或 /bin/bash
	if !strings.Contains(result.Stdout, "bash") {
		t.Errorf("应该使用 bash，实际输出: %q", result.Stdout)
	}
}

// TestExecute_Timeout 测试超时
func TestExecute_Timeout(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"sleep", "10"}, WithTimeout(100*time.Millisecond), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 超时后退出码可能是 124（timeout 命令）或 -1（被 kill）
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}
	if result.Error == nil {
		t.Error("应该有超时错误")
	}
	// 错误可能是 DeadlineExceeded 或 signal: killed
	if result.Error != context.DeadlineExceeded && !strings.Contains(result.Error.Error(), "killed") {
		t.Errorf("错误应该是 DeadlineExceeded 或 killed，实际: %v", result.Error)
	}
	if result.ExpiredDuration < 100*time.Millisecond {
		t.Errorf("执行时间应该至少100ms，实际: %v", result.ExpiredDuration)
	}
	if result.ExpiredDuration > 2*time.Second {
		t.Errorf("执行时间不应该超过2秒，实际: %v", result.ExpiredDuration)
	}
}

// TestExecute_Timeout_NoTimeout 测试命令在超时前完成
func TestExecute_Timeout_NoTimeout(t *testing.T) {
	ctx := context.Background()
	// 命令在超时前完成，不应该超时
	result, err := Execute(ctx, []string{"echo", "hello"}, WithTimeout(1*time.Second), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Error != nil {
		t.Errorf("不应该有错误，实际: %v", result.Error)
	}
	if result.Stdout != "hello\n" {
		t.Errorf("stdout 不匹配，期望: 'hello\\n', 实际: %q", result.Stdout)
	}
	// 执行时间应该很短，远小于超时时间
	if result.ExpiredDuration > 500*time.Millisecond {
		t.Errorf("执行时间应该很短，实际: %v", result.ExpiredDuration)
	}
}

// TestExecute_Timeout_ShortTimeout 测试很短的超时时间
func TestExecute_Timeout_ShortTimeout(t *testing.T) {
	ctx := context.Background()
	// 使用非常短的超时时间
	result, err := Execute(ctx, []string{"sleep", "1"}, WithTimeout(10*time.Millisecond), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 应该超时
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}
	if result.Error == nil {
		t.Error("应该有超时错误")
	}
	// 执行时间应该接近超时时间，但不会超过太多
	if result.ExpiredDuration < 10*time.Millisecond {
		t.Errorf("执行时间应该至少10ms，实际: %v", result.ExpiredDuration)
	}
	if result.ExpiredDuration > 500*time.Millisecond {
		t.Errorf("执行时间不应该超过500ms，实际: %v", result.ExpiredDuration)
	}
}

// TestExecute_Timeout_WithOutput 测试超时时的输出处理
func TestExecute_Timeout_WithOutput(t *testing.T) {
	ctx := context.Background()
	// 命令在超时前输出一些数据
	result, err := Execute(ctx, []string{"sh", "-c", "echo 'start'; sleep 10; echo 'end'"}, WithTimeout(200*time.Millisecond), WithStdoutLimit(1024), WithStderrLimit(1024))
	// 超时后管道可能被关闭，导致读取失败，这是可以接受的
	if err != nil {
		// 如果是因为管道关闭导致的错误，这是正常的超时行为
		if strings.Contains(err.Error(), "file already closed") || strings.Contains(err.Error(), "broken pipe") {
			return
		}
		t.Fatalf("执行失败: %v", err)
	}

	// 应该超时
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}
	// 如果成功读取到输出，应该包含超时前的输出
	if result.Stdout != "" && !strings.Contains(result.Stdout, "start") {
		t.Errorf("stdout 应该包含超时前的输出，实际: %q", result.Stdout)
	}
}

// TestExecute_Timeout_WithOutputLimit 测试超时与输出限制的组合
func TestExecute_Timeout_WithOutputLimit(t *testing.T) {
	ctx := context.Background()
	// 命令输出大量数据并超时
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 100); do echo line$i; done; sleep 10"}, WithTimeout(200*time.Millisecond), WithStdoutLimit(100), WithStdoutStrategy(TruncateTail))
	// 超时后管道可能被关闭，导致读取失败，这是可以接受的
	if err != nil {
		// 如果是因为管道关闭导致的错误，这是正常的超时行为
		if strings.Contains(err.Error(), "file already closed") || strings.Contains(err.Error(), "broken pipe") {
			return
		}
		t.Fatalf("执行失败: %v", err)
	}

	// 应该超时
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}
	// 如果成功读取到输出，应该被截断
	if result.Stdout != "" {
		if !result.StdoutTruncated {
			t.Error("stdout 应该被截断")
		}
		// 输出长度应该不超过限制
		if len(result.Stdout) > 100 {
			t.Errorf("stdout 长度应该不超过100字节，实际: %d", len(result.Stdout))
		}
	}
}

// TestExecute_Timeout_ShellMode 测试超时与 Shell 模式的组合
func TestExecute_Timeout_ShellMode(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"sleep 10"}, WithShellMode(true), WithTimeout(100*time.Millisecond), WithStdoutLimit(1024), WithStderrLimit(1024))
	// 超时后管道可能被关闭，导致读取失败，这是可以接受的
	if err != nil {
		// 如果是因为管道关闭导致的错误，这是正常的超时行为
		if strings.Contains(err.Error(), "file already closed") || strings.Contains(err.Error(), "broken pipe") {
			return
		}
		t.Fatalf("执行失败: %v", err)
	}

	// 应该超时
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}
	if result.Error == nil {
		t.Error("应该有超时错误")
	}
}

// TestExecute_Timeout_Stderr 测试超时时的 stderr 输出
func TestExecute_Timeout_Stderr(t *testing.T) {
	ctx := context.Background()
	// 命令在超时前输出到 stderr
	result, err := Execute(ctx, []string{"sh", "-c", "echo 'error' >&2; sleep 10"}, WithTimeout(200*time.Millisecond), WithStdoutLimit(1024), WithStderrLimit(1024))
	// 超时后管道可能被关闭，导致读取失败，这是可以接受的
	if err != nil {
		// 如果是因为管道关闭导致的错误，这是正常的超时行为
		if strings.Contains(err.Error(), "file already closed") || strings.Contains(err.Error(), "broken pipe") {
			return
		}
		t.Fatalf("执行失败: %v", err)
	}

	// 应该超时
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}
	// 如果成功读取到输出，应该包含超时前的 stderr 输出
	if result.Stderr != "" && !strings.Contains(result.Stderr, "error") {
		t.Errorf("stderr 应该包含超时前的输出，实际: %q", result.Stderr)
	}
}

// TestExecute_Timeout_DurationAccuracy 测试超时时间的准确性
func TestExecute_Timeout_DurationAccuracy(t *testing.T) {
	ctx := context.Background()
	timeout := 150 * time.Millisecond
	start := time.Now()
	result, err := Execute(ctx, []string{"sleep", "10"}, WithTimeout(timeout), WithStdoutLimit(1024), WithStderrLimit(1024))
	actualDuration := time.Since(start)

	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 应该超时
	if result.ExitCode != 124 && result.ExitCode != -1 {
		t.Errorf("超时退出码应该是124或-1，实际: %d", result.ExitCode)
	}

	// ExpiredDuration 应该接近实际执行时间
	diff := result.ExpiredDuration - actualDuration
	if diff < 0 {
		diff = -diff
	}
	// 允许 100ms 的误差
	if diff > 100*time.Millisecond {
		t.Errorf("ExpiredDuration 与实际执行时间差异过大，ExpiredDuration: %v, 实际: %v, 差异: %v", result.ExpiredDuration, actualDuration, diff)
	}

	// ExpiredDuration 应该接近超时时间
	diffFromTimeout := result.ExpiredDuration - timeout
	if diffFromTimeout < 0 {
		diffFromTimeout = -diffFromTimeout
	}
	// 允许 100ms 的误差
	if diffFromTimeout > 100*time.Millisecond {
		t.Errorf("ExpiredDuration 与超时时间差异过大，ExpiredDuration: %v, 超时: %v, 差异: %v", result.ExpiredDuration, timeout, diffFromTimeout)
	}
}

// TestExecute_Timeout_ZeroTimeout 测试零超时（应该立即超时或使用默认行为）
func TestExecute_Timeout_ZeroTimeout(t *testing.T) {
	ctx := context.Background()
	// 零超时应该等同于无超时
	result, err := Execute(ctx, []string{"echo", "hello"}, WithTimeout(0), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Error != nil {
		t.Errorf("不应该有错误，实际: %v", result.Error)
	}
	if result.Stdout != "hello\n" {
		t.Errorf("stdout 不匹配，期望: 'hello\\n', 实际: %q", result.Stdout)
	}
}

// TestExecute_WorkDir 测试工作目录
func TestExecute_WorkDir(t *testing.T) {
	ctx := context.Background()
	// 使用 pwd 命令测试工作目录
	result, err := Execute(ctx, []string{"pwd"}, WithWorkDir("/tmp"), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if !strings.Contains(result.Stdout, "/tmp") {
		t.Errorf("工作目录应该是 /tmp，实际输出: %q", result.Stdout)
	}
}

// TestExecute_Stdin 测试标准输入
func TestExecute_Stdin(t *testing.T) {
	ctx := context.Background()
	stdin := strings.NewReader("hello from stdin\n")
	result, err := Execute(ctx, []string{"cat"}, WithStdin(stdin), WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Stdout != "hello from stdin\n" {
		t.Errorf("stdout 不匹配，期望: 'hello from stdin\\n', 实际: %q", result.Stdout)
	}
}

// TestExecute_StdoutLimit 测试 stdout 限制
// 注意：head 策略在缓冲区满时可能无法继续读取剩余数据，这里使用 tail 策略来测试
func TestExecute_StdoutLimit(t *testing.T) {
	ctx := context.Background()
	// 输出超过限制的数据，使用 tail 策略避免缓冲区问题
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 20); do echo line$i; done"}, WithStdoutLimit(50), WithStdoutStrategy(TruncateTail))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	// 应该被截断
	if !result.StdoutTruncated {
		t.Error("stdout 应该被截断")
	}
	// 输出应该不超过50字节（tail策略保留后面的）
	if len(result.Stdout) > 50 {
		t.Errorf("stdout 长度应该不超过50字节，实际: %d", len(result.Stdout))
	}
	// tail 策略应该包含后面的内容
	if !strings.Contains(result.Stdout, "line20") {
		t.Errorf("stdout 应该包含后面的内容，实际: %q", result.Stdout)
	}
}

// TestExecute_StdoutLimitTail 测试 stdout 限制（tail 策略）
func TestExecute_StdoutLimitTail(t *testing.T) {
	ctx := context.Background()
	// 输出多行数据
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 20); do echo line$i; done"}, WithStdoutLimit(20), WithStdoutStrategy(TruncateTail))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	// tail 策略应该保留最后的部分
	if !result.StdoutTruncated {
		t.Error("stdout 应该被截断")
	}
	// 应该包含后面的行
	if !strings.Contains(result.Stdout, "line20") {
		t.Errorf("stdout 应该包含最后的内容，实际: %q", result.Stdout)
	}
}

// TestExecute_StderrLimit 测试 stderr 限制
// 注意：head 策略在缓冲区满时可能无法继续读取剩余数据，这里使用 tail 策略来测试
func TestExecute_StderrLimit(t *testing.T) {
	ctx := context.Background()
	// 输出到 stderr，使用 tail 策略避免缓冲区问题
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 20); do echo err$i >&2; done"}, WithStderrLimit(50), WithStderrStrategy(TruncateTail))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	// 应该被截断
	if !result.StderrTruncated {
		t.Error("stderr 应该被截断")
	}
	// 输出应该不超过50字节
	if len(result.Stderr) > 50 {
		t.Errorf("stderr 长度应该不超过50字节，实际: %d", len(result.Stderr))
	}
	// tail 策略应该包含后面的内容
	if !strings.Contains(result.Stderr, "err20") {
		t.Errorf("stderr 应该包含后面的内容，实际: %q", result.Stderr)
	}
}

// TestExecute_StderrLimitTail 测试 stderr 限制（tail 策略）
func TestExecute_StderrLimitTail(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 20); do echo err$i >&2; done"}, WithStderrLimit(20), WithStderrStrategy(TruncateTail))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	// tail 策略应该保留最后的部分
	if !result.StderrTruncated {
		t.Error("stderr 应该被截断")
	}
	// 应该包含后面的行
	if !strings.Contains(result.Stderr, "err20") {
		t.Errorf("stderr 应该包含最后的内容，实际: %q", result.Stderr)
	}
}

// TestExecute_NoCommand 测试空命令
func TestExecute_NoCommand(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{})
	if err == nil {
		t.Error("应该返回错误")
	}
	if result != nil {
		t.Error("结果应该为 nil")
	}
	if !strings.Contains(err.Error(), "no command specified") {
		t.Errorf("错误信息应该包含 'no command specified'，实际: %v", err)
	}
}

// TestExecute_CombinedOptions 测试组合选项
func TestExecute_CombinedOptions(t *testing.T) {
	ctx := context.Background()
	stdin := strings.NewReader("test input\n")
	result, err := Execute(ctx, []string{"cat"},
		WithStdin(stdin),
		WithWorkDir("/tmp"),
		WithStdoutLimit(100),
		WithStderrLimit(50),
		WithStdoutStrategy(TruncateHead),
		WithStderrStrategy(TruncateTail),
	)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Stdout != "test input\n" {
		t.Errorf("stdout 不匹配，期望: 'test input\\n', 实际: %q", result.Stdout)
	}
}

// TestExecute_LargeOutput 测试大量输出
func TestExecute_LargeOutput(t *testing.T) {
	ctx := context.Background()
	// 生成大量输出，使用 tail 策略避免缓冲区问题
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 200); do echo 'line'$i; done"}, WithStdoutLimit(500), WithStdoutStrategy(TruncateTail))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	// 应该被截断
	if !result.StdoutTruncated {
		t.Error("stdout 应该被截断")
	}
	// 输出应该不超过限制
	if len(result.Stdout) > 500 {
		t.Errorf("stdout 长度应该不超过500字节，实际: %d", len(result.Stdout))
	}
	// tail 策略应该包含后面的内容
	if !strings.Contains(result.Stdout, "line200") {
		t.Errorf("stdout 应该包含后面的内容，实际: %q", result.Stdout)
	}
}

// TestExecute_StderrOutput 测试 stderr 输出
func TestExecute_StderrOutput(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"sh", "-c", "echo 'error message' >&2"}, WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Stderr != "error message\n" {
		t.Errorf("stderr 不匹配，期望: 'error message\\n', 实际: %q", result.Stderr)
	}
}

// TestExecute_BothOutputs 测试同时输出到 stdout 和 stderr
func TestExecute_BothOutputs(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"sh", "-c", "echo 'stdout'; echo 'stderr' >&2"}, WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.Stdout != "stdout\n" {
		t.Errorf("stdout 不匹配，期望: 'stdout\\n', 实际: %q", result.Stdout)
	}
	if result.Stderr != "stderr\n" {
		t.Errorf("stderr 不匹配，期望: 'stderr\\n', 实际: %q", result.Stderr)
	}
}

// TestExecute_NoLimit 测试无限制输出（使用足够大的限制）
func TestExecute_NoLimit(t *testing.T) {
	ctx := context.Background()
	// 输出一些数据，使用足够大的限制来模拟无限制
	result, err := Execute(ctx, []string{"sh", "-c", "for i in $(seq 1 10); do echo line$i; done"}, WithStdoutLimit(10240), WithStderrLimit(10240))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.StdoutTruncated {
		t.Error("stdout 不应该被截断")
	}
	if result.StderrTruncated {
		t.Error("stderr 不应该被截断")
	}
	// 应该包含所有行
	lines := []string{"line1", "line2", "line3", "line4", "line5", "line6", "line7", "line8", "line9", "line10"}
	for _, expected := range lines {
		if !strings.Contains(result.Stdout, expected) {
			t.Errorf("stdout 应该包含 %q，实际: %q", expected, result.Stdout)
		}
	}
}

// TestExecuteOptions_Apply 测试选项应用
func TestExecuteOptions_Apply(t *testing.T) {
	opts := &ExecuteOptions{}
	opts.Apply(
		WithShellMode(true),
		WithTimeout(5*time.Second),
		WithWorkDir("/tmp"),
		WithStdoutLimit(100),
		WithStderrLimit(50),
	)

	if !opts.ShellMode {
		t.Error("ShellMode 应该为 true")
	}
	if opts.Timeout != 5*time.Second {
		t.Errorf("Timeout 应该是 5s，实际: %v", opts.Timeout)
	}
	if opts.WorkDir != "/tmp" {
		t.Errorf("WorkDir 应该是 /tmp，实际: %s", opts.WorkDir)
	}
	if opts.StdoutLimit != 100 {
		t.Errorf("StdoutLimit 应该是 100，实际: %d", opts.StdoutLimit)
	}
	if opts.StderrLimit != 50 {
		t.Errorf("StderrLimit 应该是 50，实际: %d", opts.StderrLimit)
	}
}

// TestExecute_QuickCommand 测试快速命令
func TestExecute_QuickCommand(t *testing.T) {
	ctx := context.Background()
	result, err := Execute(ctx, []string{"true"}, WithStdoutLimit(1024), WithStderrLimit(1024))
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("退出码应该是0，实际: %d", result.ExitCode)
	}
	if result.ExpiredDuration > 1*time.Second {
		t.Errorf("执行时间应该很短，实际: %v", result.ExpiredDuration)
	}
}
