/*
 * Author: fasion
 * Created time: 2023-03-15 11:35:06
 * Last Modified by: fasion
 * Last Modified time: 2023-03-15 14:59:01
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
