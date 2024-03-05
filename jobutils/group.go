/*
 * Author: fasion
 * Created time: 2023-03-15 11:35:06
 * Last Modified by: fasion
 * Last Modified time: 2024-03-05 14:53:19
 */

package jobutils

import (
	"context"
	"sync"
	"time"
)

var bgCtx = context.Background()

type JobTokens chan struct{}

func NewJobTokens(n int) JobTokens {
	if n >= 0 {
		return make(JobTokens, n)
	} else {
		return nil
	}
}

func (tokens JobTokens) Totals() int {
	if tokens == nil {
		return 0
	} else {
		return cap(tokens)
	}
}

func (tokens JobTokens) Useds() int {
	if tokens == nil {
		return 0
	} else {
		return len(tokens)
	}
}

func (tokens JobTokens) Lefts() int {
	return tokens.Totals() - tokens.Useds()
}

func (tokens JobTokens) Acquire(ctx context.Context, timeout time.Duration) bool {
	if tokens == nil {
		return true
	}

	if ctx == nil {
		ctx = bgCtx
	}

	switch {
	case timeout < 0: // 永久等待
		select {
		case tokens <- struct{}{}:
			return true
		case <-ctx.Done():
			return false
		}
	case timeout == 0: // 不等待
		select {
		case tokens <- struct{}{}:
			return true
		case <-ctx.Done():
			return false
		default:
			return false
		}
	default: // 自定义等待时间
		select {
		case tokens <- struct{}{}:
			return true
		case <-ctx.Done():
			return false
		case <-time.After(timeout):
			return false
		}
	}
}

func (tokens JobTokens) Release() {
	if tokens != nil {
		<-tokens
	}
}

type JobGroup struct {
	tokens       JobTokens
	wg           sync.WaitGroup
	enterTimeout time.Duration
}

func NewJobGroup(tokens int, enterTimeout time.Duration) *JobGroup {
	return &JobGroup{
		tokens:       NewJobTokens(tokens),
		enterTimeout: enterTimeout,
	}
}

func (group *JobGroup) Tokens() JobTokens {
	return group.tokens
}

func (group *JobGroup) EnterPro(ctx context.Context, enterTimeout time.Duration) bool {
	if group == nil {
		return true
	}

	ok := group.tokens.Acquire(ctx, enterTimeout)
	if ok {
		group.wg.Add(1)
	}

	return ok
}

func (group *JobGroup) Enter(ctx context.Context) bool {
	return group.EnterPro(ctx, group.enterTimeout)
}

func (group *JobGroup) Exit() {
	if group == nil {
		return
	}

	group.tokens.Release()
	group.wg.Done()
}

func (group *JobGroup) Wait() {
	if group == nil {
		return
	}

	group.wg.Wait()
}

type TokenGenerator struct {
	Tokens chan struct{}

	maxTokens     int
	initalTokens  int
	batchSize     int
	batchInterval time.Duration
}

func (generator *TokenGenerator) StartWith(maxTokens int, initalTokens int, batchSize int, batchInterval time.Duration) *TokenGenerator {
	return generator.
		WithMaxTokens(maxTokens).
		WithInitalTokens(initalTokens).
		WithBatchSize(batchSize).
		WithBatchInterval(batchInterval).
		Start()
}

func (generator *TokenGenerator) Start() *TokenGenerator {
	tokens := make(chan struct{}, generator.maxTokens)
	for i := 0; i < generator.initalTokens; i++ {
		tokens <- struct{}{}
	}
	generator.Tokens = tokens

	go func() {
		for {
			time.Sleep(generator.batchInterval)

			for i := 0; i < generator.batchSize; i++ {
				generator.Tokens <- struct{}{}
			}
		}
	}()

	return generator
}

func (generator *TokenGenerator) WithMaxTokens(maxTokens int) *TokenGenerator {
	generator.maxTokens = maxTokens
	return generator
}

func (generator *TokenGenerator) WithInitalTokens(initalTokens int) *TokenGenerator {
	generator.initalTokens = initalTokens
	return generator
}

func (generator *TokenGenerator) WithBatchSize(batchSize int) *TokenGenerator {
	generator.batchSize = batchSize
	return generator
}

func (generator *TokenGenerator) WithBatchInterval(batchInterval time.Duration) *TokenGenerator {
	generator.batchInterval = batchInterval
	return generator
}

func (generator *TokenGenerator) Acquire(ctx context.Context, timeout time.Duration) bool {
	tokens := generator.Tokens
	if tokens == nil {
		return true
	}

	switch {
	case timeout < 0:
		select {
		case <-tokens:
			return true
		case <-ctx.Done():
			return false
		}
	case timeout == 0:
		select {
		case <-tokens:
			return true
		case <-ctx.Done():
			return false
		default:
			return false
		}
	default:
		select {
		case <-tokens:
			return true
		case <-ctx.Done():
			return false
		case <-time.After(timeout):
			return false
		}
	}
}
