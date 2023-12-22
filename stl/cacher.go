/*
 * Author: fasion
 * Created time: 2023-11-24 14:46:12
 * Last Modified by: fasion
 * Last Modified time: 2023-12-22 14:15:27
 */

package stl

import (
	"context"
	"sync"
	"time"
)

type CachedDataFetcherCallbackFunc[Data any] func(context.Context, Data, time.Time)

type CachedDataFetcherCallback[Data any] struct {
	callback CachedDataFetcherCallbackFunc[Data]
	timeout  time.Duration
}

func NewCachedDataFetcherCallback[Data any](callback func(context.Context, Data, time.Time), timeout time.Duration) *CachedDataFetcherCallback[Data] {
	return &CachedDataFetcherCallback[Data]{
		callback: callback,
		timeout:  timeout,
	}
}

func (callback *CachedDataFetcherCallback[Data]) Call(data Data, lastFetchTime time.Time) {
	go callback.call(data, lastFetchTime)
}

func (callback *CachedDataFetcherCallback[Data]) call(data Data, lastFetchTime time.Time) {
	var cancel context.CancelFunc

	ctx := context.Background()
	if callback.timeout <= 0 {
		goto doCall
	}

	ctx, cancel = context.WithTimeout(ctx, callback.timeout)
	defer cancel()

doCall:
	callback.callback(ctx, data, lastFetchTime)
	return
}

type CachedDataFetcherCallbacks[Data any] []*CachedDataFetcherCallback[Data]

func (callbacks CachedDataFetcherCallbacks[Data]) Append(others ...*CachedDataFetcherCallback[Data]) CachedDataFetcherCallbacks[Data] {
	return append(callbacks, others...)
}

func (callbacks CachedDataFetcherCallbacks[Data]) Call(data Data, lastFetchTime time.Time) {
	for _, callback := range callbacks {
		callback.Call(data, lastFetchTime)
	}
}

type CachedDataFetcher[Data any] struct {
	fetcher         func(context.Context) (Data, time.Time, error)
	expiresDuration time.Duration

	callbacks CachedDataFetcherCallbacks[Data]

	data          Data
	lastFetchTime time.Time

	mutex sync.Mutex
}

func NewCachedDataFetcher[Data any](fetcher func(context.Context) (Data, time.Time, error)) *CachedDataFetcher[Data] {
	return &CachedDataFetcher[Data]{
		fetcher: fetcher,
	}
}

func NewCachedDataFetcherLite[Data any](fetcher func(context.Context) (Data, error)) *CachedDataFetcher[Data] {
	return NewCachedDataFetcher(func(ctx context.Context) (data Data, t time.Time, err error) {
		data, err = fetcher(ctx)
		if err == nil {
			t = time.Now()
		}
		return
	})
}

func (fetcher *CachedDataFetcher[Data]) Dup() *CachedDataFetcher[Data] {
	return &CachedDataFetcher[Data]{
		fetcher:         fetcher.fetcher,
		expiresDuration: fetcher.expiresDuration,
	}
}

func (fetcher *CachedDataFetcher[Data]) EnsureExpiresDuration(expiresDuration time.Duration) *CachedDataFetcher[Data] {
	if expiresDuration <= 0 {
		return fetcher
	}

	if fetcher.expiresDuration <= 0 {
		fetcher.expiresDuration = expiresDuration
		return fetcher
	}

	if expiresDuration < fetcher.expiresDuration {
		fetcher.expiresDuration = expiresDuration
	}

	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) WithFetcher(fetcherFunc func(context.Context) (Data, time.Time, error)) *CachedDataFetcher[Data] {
	fetcher.fetcher = fetcherFunc
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) WithExpiresDuration(duration time.Duration) *CachedDataFetcher[Data] {
	fetcher.expiresDuration = duration
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) NewCallbackLite(callback func(context.Context), timeout time.Duration) *CachedDataFetcherCallback[Data] {
	return fetcher.RegisterCallback(NewCachedDataFetcherCallback(func(ctx context.Context, data Data, t time.Time) {
		callback(ctx)
	}, timeout))
}

func (fetcher *CachedDataFetcher[Data]) NewCallback(callback CachedDataFetcherCallbackFunc[Data], timeout time.Duration) *CachedDataFetcherCallback[Data] {
	return fetcher.RegisterCallback(NewCachedDataFetcherCallback(callback, timeout))
}

func (fetcher *CachedDataFetcher[Data]) RegisterCallback(callback *CachedDataFetcherCallback[Data]) *CachedDataFetcherCallback[Data] {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.registerCallback(callback)
}

func (fetcher *CachedDataFetcher[Data]) registerCallback(callback *CachedDataFetcherCallback[Data]) *CachedDataFetcherCallback[Data] {
	fetcher.callbacks = fetcher.callbacks.Append(callback)
	return callback
}

func (fetcher *CachedDataFetcher[Data]) GetCached() (Data, time.Time) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.getCached()
}

func (fetcher *CachedDataFetcher[Data]) Get() (Data, bool) {
	return fetcher.GetWithExpires(0)
}

func (fetcher *CachedDataFetcher[Data]) GetWithExpires(expiresDuration time.Duration) (data Data, ok bool) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.getWithExpires(expiresDuration)
}

func (fetcher *CachedDataFetcher[Data]) GetWithSince(since time.Time) (data Data, ok bool) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.getWithSince(since)
}

func (fetcher *CachedDataFetcher[Data]) Fetch(ctx context.Context) (Data, bool, error) {
	return fetcher.FetchWithExpires(ctx, 0)
}

func (fetcher *CachedDataFetcher[Data]) FetchWithExpires(ctx context.Context, expiresDuration time.Duration) (data Data, ok bool, err error) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	data, ok = fetcher.getWithExpires(expiresDuration)
	if ok {
		return
	}

	return fetcher.refresh(ctx)
}

func (fetcher *CachedDataFetcher[Data]) FetchWithSince(ctx context.Context, since time.Time) (data Data, ok bool, err error) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	data, ok = fetcher.getWithSince(since)
	if ok {
		return
	}

	return fetcher.refresh(ctx)
}

func (fetcher *CachedDataFetcher[Data]) TriggerRefresh(ctx context.Context) {
	fetcher.Refresh(ctx)
}

func (fetcher *CachedDataFetcher[Data]) Refresh(ctx context.Context) (data Data, ok bool, err error) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.refresh(ctx)
}

func (fetcher *CachedDataFetcher[Data]) refresh(ctx context.Context) (data Data, ok bool, err error) {
	data, _, err = fetcher.refetch(ctx)
	ok = err == nil
	return
}

func (fetcher *CachedDataFetcher[Data]) refetch(ctx context.Context) (data Data, t time.Time, err error) {
	data, t, err = fetcher.fetcher(ctx)
	if err != nil {
		return
	}

	if t.IsZero() {
		t = time.Now()
	}

	fetcher.data = data
	fetcher.lastFetchTime = t

	// call it asynchronously to avoid dead lock
	// in case that callbacks may call fetcher method again
	go fetcher.callbacks.Call(data, t)

	return
}

func (fetcher *CachedDataFetcher[Data]) getCached() (Data, time.Time) {
	return fetcher.data, fetcher.lastFetchTime
}

func (fetcher *CachedDataFetcher[Data]) getWithExpires(expiresDuration time.Duration) (Data, bool) {
	if expiresDuration <= 0 {
		expiresDuration = fetcher.expiresDuration
	}

	var since time.Time
	if expiresDuration > 0 {
		since = time.Now().Add(-expiresDuration)
	}

	return fetcher.getWithSince(since)
}

func (fetcher *CachedDataFetcher[Data]) getWithSince(since time.Time) (Data, bool) {
	if fetcher.lastFetchTime.IsZero() {
		return fetcher.data, false
	}

	return fetcher.data, fetcher.lastFetchTime.After(since)
}
