/*
 * Author: fasion
 * Created time: 2026-01-25 14:33:28
 * Last Modified by: fasion
 * Last Modified time: 2026-01-26 23:03:22
 */

package execx

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// ExecuteOptions 执行选项
type ExecuteOptions struct {
	ShellMode      bool
	Shell          string // shell 命令，默认为 "sh"
	Timeout        time.Duration
	WorkDir        string
	Stdin          io.Reader
	StdoutLimit    int              // stdout 输出长度限制，0 表示不限制
	StdoutStrategy TruncateStrategy // stdout 截断策略，head 或 tail
	StderrLimit    int              // stderr 输出长度限制，0 表示不限制
	StderrStrategy TruncateStrategy // stderr 截断策略，head 或 tail
}

// Apply 应用选项
func (opts *ExecuteOptions) Apply(options ...Option) {
	for _, opt := range options {
		opt(opts)
	}
}

// Option 选项函数类型
type Option func(*ExecuteOptions)

// WithShellMode 设置 shell 模式
func WithShellMode(shellMode bool) Option {
	return func(opts *ExecuteOptions) {
		opts.ShellMode = shellMode
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(opts *ExecuteOptions) {
		opts.Timeout = timeout
	}
}

// WithWorkDir 设置工作目录
func WithWorkDir(workDir string) Option {
	return func(opts *ExecuteOptions) {
		opts.WorkDir = workDir
	}
}

// WithStdin 设置标准输入
func WithStdin(stdin io.Reader) Option {
	return func(opts *ExecuteOptions) {
		opts.Stdin = stdin
	}
}

// WithStdoutLimit 设置 stdout 输出长度限制
func WithStdoutLimit(limit int) Option {
	return func(opts *ExecuteOptions) {
		opts.StdoutLimit = limit
	}
}

// WithStderrLimit 设置 stderr 输出长度限制
func WithStderrLimit(limit int) Option {
	return func(opts *ExecuteOptions) {
		opts.StderrLimit = limit
	}
}

// WithStdoutStrategy 设置 stdout 截断策略
func WithStdoutStrategy(strategy TruncateStrategy) Option {
	return func(opts *ExecuteOptions) {
		opts.StdoutStrategy = strategy
	}
}

// WithStderrStrategy 设置 stderr 截断策略
func WithStderrStrategy(strategy TruncateStrategy) Option {
	return func(opts *ExecuteOptions) {
		opts.StderrStrategy = strategy
	}
}

// WithShell 设置 shell 命令（仅在 ShellMode 为 true 时生效）
func WithShell(shell string) Option {
	return func(opts *ExecuteOptions) {
		opts.Shell = shell
	}
}

type ExecuteResult struct {
	ExitCode        int
	Stdout          string
	StdoutTruncated bool
	Stderr          string
	StderrTruncated bool
	StartTime       time.Time
	EndTime         time.Time
	ExpiredDuration time.Duration
	Error           error
}

// Execute 在当前系统执行命令
func Execute(ctx context.Context, args []string, opts ...Option) (*ExecuteResult, error) {
	// 构建默认选项
	options := &ExecuteOptions{
		ShellMode:      false,
		Shell:          "sh", // 默认使用 sh
		Timeout:        0,
		WorkDir:        "",
		Stdin:          nil,
		StdoutLimit:    0,
		StdoutStrategy: TruncateHead, // 默认保留最先输出的部分
		StderrLimit:    0,
		StderrStrategy: TruncateHead, // 默认保留最先输出的部分
	}

	// 应用选项
	options.Apply(opts...)

	startTime := time.Now()

	// 如果指定了超时，创建带超时的 context
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	var cmd *exec.Cmd

	if options.ShellMode {
		// Shell 模式：使用指定的 shell -c 执行命令字符串
		commandStr := strings.Join(args, " ")
		cmd = exec.CommandContext(ctx, options.Shell, "-c", commandStr)
	} else {
		// 参数模式：直接执行命令
		if len(args) == 0 {
			return nil, fmt.Errorf("no command specified")
		}
		cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	}

	if options.WorkDir != "" {
		cmd.Dir = options.WorkDir
	}

	if options.Stdin != nil {
		cmd.Stdin = options.Stdin
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// 并发读取输出
	stdoutChan := make(chan BytesReadResult, 1)
	stderrChan := make(chan BytesReadResult, 1)

	// 读取 stdout
	go func() {
		var result BytesReadResult
		if options.StdoutLimit > 0 {
			result = Read(stdout,
				WithLimit(options.StdoutLimit),
				WithStrategy(options.StdoutStrategy),
				WithReadRest(true), // 读取剩余数据，防止阻塞
			)
		} else {
			result = Read(stdout)
		}
		stdoutChan <- result
	}()

	// 读取 stderr
	go func() {
		var result BytesReadResult
		if options.StderrLimit > 0 {
			result = Read(stderr,
				WithLimit(options.StderrLimit),
				WithStrategy(options.StderrStrategy),
				WithReadRest(true), // 读取剩余数据，防止阻塞
			)
		} else {
			result = Read(stderr)
		}
		stderrChan <- result
	}()

	// 等待命令完成
	waitErr := cmd.Wait()
	endTime := time.Now()

	// 读取输出
	stdoutResult := <-stdoutChan
	stderrResult := <-stderrChan

	if stdoutResult.Error != nil {
		return nil, fmt.Errorf("failed to read stdout: %w", stdoutResult.Error)
	}
	if stderrResult.Error != nil {
		return nil, fmt.Errorf("failed to read stderr: %w", stderrResult.Error)
	}

	stdoutBytes := stdoutResult.Data
	stderrBytes := stderrResult.Data
	stdoutTruncated := stdoutResult.Truncated
	stderrTruncated := stderrResult.Truncated

	stdoutStr := string(stdoutBytes)
	stderrStr := string(stderrBytes)

	exitCode := 0
	var expiredDuration time.Duration
	if waitErr != nil {
		if exitError, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// Context 超时
			if ctx.Err() == context.DeadlineExceeded {
				expiredDuration = endTime.Sub(startTime)
				return &ExecuteResult{
					ExitCode:        124,
					Stdout:          stdoutStr,
					StdoutTruncated: stdoutTruncated,
					Stderr:          stderrStr,
					StderrTruncated: stderrTruncated,
					StartTime:       startTime,
					EndTime:         endTime,
					ExpiredDuration: expiredDuration,
					Error:           ctx.Err(),
				}, nil
			}
			return nil, fmt.Errorf("command execution failed: %w", waitErr)
		}
	}

	expiredDuration = endTime.Sub(startTime)

	return &ExecuteResult{
		ExitCode:        exitCode,
		Stdout:          stdoutStr,
		StdoutTruncated: stdoutTruncated,
		Stderr:          stderrStr,
		StderrTruncated: stderrTruncated,
		StartTime:       startTime,
		EndTime:         endTime,
		ExpiredDuration: expiredDuration,
		Error:           waitErr,
	}, nil
}
