/*
 * Author: fasion
 * Created time: 2025-05-26 11:36:03
 * Last Modified by: fasion
 * Last Modified time: 2025-06-12 18:48:47
 */

package httpx

import (
	"github.com/avast/retry-go"
	"golang.org/x/net/context"
)

const (
	ContextKeyRetryRef = "__RetryRef__"
)

type RetryRef struct {
	opts []retry.Option
}

func NewRetryRef(opts ...retry.Option) *RetryRef {
	return &RetryRef{opts: opts}
}

func (ref *RetryRef) Push(opts ...retry.Option) *RetryRef {
	if ref == nil {
		return nil
	}

	ref.opts = opts

	return ref
}

func (ref *RetryRef) Pop() (opts RetryOptions) {
	if ref == nil {
		return nil
	}

	opts = ref.opts
	ref.opts = nil

	return
}

func (ref *RetryRef) Do(fn func() error) error {
	return ref.Pop().Do(fn)
}

func (ref *RetryRef) DoIfSpecified(fn func() error) (err error, specified bool) {
	return ref.Pop().DoIfSpecified(fn)
}

type RetryOptions []retry.Option

func (opts RetryOptions) Do(fn func() error) error {
	if opts == nil {
		return fn()
	}

	return retry.Do(fn, opts...)
}

func (opts RetryOptions) DoIfSpecified(fn func() error) (err error, specified bool) {
	if opts == nil {
		return nil, false
	}

	return retry.Do(fn, opts...), true
}

func ContextWithRetryOptions(ctx context.Context, opts ...retry.Option) context.Context {
	ref, ctx := RetryRefFromContextPro(ctx, true, true, opts...)
	ref.Push(opts...)
	return ctx
}

func ContextWithRetryRef(ctx context.Context, ref *RetryRef) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, ContextKeyRetryRef, ref)
}

func RetryRefFromContextPro(ctx context.Context, create, wrapContext bool, opts ...retry.Option) (*RetryRef, context.Context) {
	ref := RetryRefFromContext(ctx)
	if ref != nil {
		return ref, ctx
	}

	if !create {
		return nil, ctx
	}

	ref = NewRetryRef(opts...)
	if wrapContext {
		if ctx == nil {
			ctx = context.Background()
		}

		ctx = ContextWithRetryRef(ctx, ref)
	}

	return ref, ctx
}

func RetryRefFromContext(ctx context.Context) (ref *RetryRef) {
	if ctx == nil {
		return nil
	}

	ref, _ = ctx.Value(ContextKeyRetryRef).(*RetryRef)
	return
}

func RetryByContext(ctx context.Context, fn func() error) error {
	return RetryRefFromContext(ctx).Do(fn)
}

func RetryIfContextSpecified(ctx context.Context, fn func() error) (err error, specified bool) {
	return RetryRefFromContext(ctx).DoIfSpecified(fn)
}
