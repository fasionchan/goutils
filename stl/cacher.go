/*
 * Author: fasion
 * Created time: 2023-11-24 14:46:12
 * Last Modified by: fasion
 * Last Modified time: 2023-12-15 09:22:20
 */

package stl

import (
	"context"
	"sync"
	"time"
)

type CachedDataFetcher[Data any] struct {
	fetcher         func(context.Context) (Data, time.Time, error)
	expiresDuration time.Duration

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

func (fetcher *CachedDataFetcher[Data]) WithFetcher(fetcherFunc func(context.Context) (Data, time.Time, error)) *CachedDataFetcher[Data] {
	fetcher.fetcher = fetcherFunc
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) WithExpiresDuration(duration time.Duration) *CachedDataFetcher[Data] {
	fetcher.expiresDuration = duration
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) Get(ctx context.Context) (Data, bool) {
	return fetcher.GetPro(ctx, 0)
}

func (fetcher *CachedDataFetcher[Data]) GetPro(ctx context.Context, expiresDuration time.Duration) (data Data, ok bool) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.get(expiresDuration)
}

func (fetcher *CachedDataFetcher[Data]) Fetch(ctx context.Context) (Data, bool, error) {
	return fetcher.FetchPro(ctx, 0)
}

func (fetcher *CachedDataFetcher[Data]) FetchPro(ctx context.Context, expiresDuration time.Duration) (data Data, ok bool, err error) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	data, ok = fetcher.get(expiresDuration)
	if ok {
		return
	}

	return fetcher.refresh(ctx)
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

	fetcher.data = data
	fetcher.lastFetchTime = t

	return
}

func (fetcher *CachedDataFetcher[Data]) get(expiresDuration time.Duration) (Data, bool) {
	if fetcher.lastFetchTime.IsZero() {
		return fetcher.data, false
	}

	if expiresDuration <= 0 {
		expiresDuration = fetcher.expiresDuration
	}

	if expiresDuration <= 0 {
		return fetcher.data, true
	}

	return fetcher.data, time.Since(fetcher.lastFetchTime) <= expiresDuration
}
