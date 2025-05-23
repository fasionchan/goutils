/*
 * Author: fasion
 * Created time: 2023-11-24 14:46:12
 * Last Modified by: fasion
 * Last Modified time: 2025-01-02 11:01:03
 */

package stl

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	NopLogger = zap.NewNop()
)

type CachedDataFetcherFetchFunc[Data any] func(ctx context.Context, sinceTime time.Time) (Data, time.Time, error)
type CachedDataFetcherFetchFuncLite[Data any] func(ctx context.Context) (Data, error)
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
	if callback.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, callback.timeout)
		defer cancel()
	}

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

type TimedValue[Value any] struct {
	value Value
	t     time.Time
}

func NewTimedValue[Value any](value Value, t time.Time) *TimedValue[Value] {
	return &TimedValue[Value]{
		value: value,
		t:     t,
	}
}

func (timed *TimedValue[Value]) Value() (value Value) {
	if timed == nil {
		return
	}

	return timed.value
}

func (timed *TimedValue[Value]) Time() (t time.Time) {
	if timed == nil {
		return
	}

	return timed.t
}

func (timed *TimedValue[Value]) ValueAndTime() (value Value, t time.Time) {
	if timed == nil {
		return
	}

	return timed.value, timed.t
}

type CachedDataFetcher[Data any] struct {
	*zap.Logger

	fetcher         CachedDataFetcherFetchFunc[Data]
	expiresDuration time.Duration

	callbacks CachedDataFetcherCallbacks[Data]

	data *TimedValue[Data]

	mutex sync.Mutex
}

func NewCachedDataFetcher[Data any](fetcher CachedDataFetcherFetchFunc[Data]) *CachedDataFetcher[Data] {
	return &CachedDataFetcher[Data]{
		Logger:  zap.NewNop(),
		fetcher: fetcher,
	}
}

func NewCachedDataFetcherLite[Data any](fetcher CachedDataFetcherFetchFuncLite[Data]) *CachedDataFetcher[Data] {
	return NewCachedDataFetcher(func(ctx context.Context, sinceTime time.Time) (data Data, t time.Time, err error) {
		fetchingTime := time.Now()
		data, err = fetcher(ctx)
		if err == nil {
			t = fetchingTime
		}
		return
	})
}

func (fetcher *CachedDataFetcher[Data]) BuildAccessor() *CachedDataFetcherAccessor[Data] {
	return &CachedDataFetcherAccessor[Data]{
		fetcher: fetcher,
	}
}

func (fetcher *CachedDataFetcher[Data]) Dup() *CachedDataFetcher[Data] {
	return &CachedDataFetcher[Data]{
		fetcher:         fetcher.fetcher,
		expiresDuration: fetcher.expiresDuration,
	}
}

func (fetcher *CachedDataFetcher[Data]) WithCachedDataPurged() *CachedDataFetcher[Data] {
	fetcher.data = nil
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) WithLogger(logger *zap.Logger) *CachedDataFetcher[Data] {
	fetcher.Logger = logger
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) PurgeCachedData() {
	fetcher.WithCachedDataPurged()
}

func (fetcher *CachedDataFetcher[Data]) SinceTimeFromExpiresDuration(expiresDuration time.Duration) (since time.Time) {
	if expiresDuration <= 0 {
		expiresDuration = fetcher.expiresDuration
	}

	if expiresDuration > 0 {
		since = time.Now().Add(-expiresDuration)
	}

	return
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

func (fetcher *CachedDataFetcher[Data]) WithFetcher(fetcherFunc CachedDataFetcherFetchFunc[Data]) *CachedDataFetcher[Data] {
	fetcher.fetcher = fetcherFunc
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) WithExpiresDuration(duration time.Duration) *CachedDataFetcher[Data] {
	fetcher.expiresDuration = duration
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) WithOthersSubscribed(timeout time.Duration, others ...interface {
	RegisterCallbackFuncLite(func(context.Context), time.Duration)
}) *CachedDataFetcher[Data] {
	for _, other := range others {
		other.RegisterCallbackFuncLite(fetcher.TriggerRefreshLowerCache, timeout)
	}
	return fetcher
}

func (fetcher *CachedDataFetcher[Data]) RegisterCallbackFuncLite(callback func(context.Context), timeout time.Duration) {
	fetcher.NewCallbackLite(callback, timeout)
}

func (fetcher *CachedDataFetcher[Data]) NewCallbackLite(callback func(context.Context), timeout time.Duration) *CachedDataFetcherCallback[Data] {
	return fetcher.RegisterCallback(NewCachedDataFetcherCallback(func(ctx context.Context, data Data, t time.Time) {
		callback(ctx)
	}, timeout))
}

func (fetcher *CachedDataFetcher[Data]) RegisterCallbackFunc(callback CachedDataFetcherCallbackFunc[Data], timeout time.Duration) {
	fetcher.NewCallback(callback, timeout)
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
	return fetcher.getCached()
}

func (fetcher *CachedDataFetcher[Data]) Get() (Data, time.Time, bool) {
	return fetcher.GetWithExpires(0)
}

func (fetcher *CachedDataFetcher[Data]) GetWithExpires(expiresDuration time.Duration) (Data, time.Time, bool) {
	return fetcher.getWithExpires(expiresDuration)
}

func (fetcher *CachedDataFetcher[Data]) GetWithExpiresWarn(expiresDuration time.Duration) (Data, bool) {
	return fetcher.BuildAccessor().WithExpiresDuration(expiresDuration).WithLogExpired(true).Get()
}

func (fetcher *CachedDataFetcher[Data]) GetWithSince(since time.Time) (Data, time.Time, bool) {
	return fetcher.getWithSince(since)
}

func (fetcher *CachedDataFetcher[Data]) GetWithSinceWarn(since time.Time) (data Data, ok bool) {
	data, fetchingTime, ok := fetcher.getWithSince(since)
	if !ok {
		fetcher.Warn("GetWithSinceWarn",
			zap.Time("FetchingTime", fetchingTime),
			zap.Time("Since", since),
		)
	}
	return
}

func (fetcher *CachedDataFetcher[Data]) GetFetchLite() func(ctx context.Context) (Data, error) {
	return fetcher.FetchLite
}

func (fetcher *CachedDataFetcher[Data]) FetchLite(ctx context.Context) (data Data, err error) {
	data, _, err = fetcher.Fetch(ctx)
	return
}

func (fetcher *CachedDataFetcher[Data]) GetFetch() func(ctx context.Context) (Data, time.Time, error) {
	return fetcher.Fetch
}

func (fetcher *CachedDataFetcher[Data]) Fetch(ctx context.Context) (Data, time.Time, error) {
	return fetcher.FetchWithExpires(ctx, 0)
}

func (fetcher *CachedDataFetcher[Data]) FetchWithExpires(ctx context.Context, expiresDuration time.Duration) (data Data, t time.Time, err error) {
	return fetcher.FetchWithSince(ctx, fetcher.SinceTimeFromExpiresDuration(expiresDuration))
}

func (fetcher *CachedDataFetcher[Data]) FetchWithExpiresPro(ctx context.Context, expiresDuration, fallbackDuration time.Duration, logger *zap.Logger) (data Data, ok bool) {
	data, t, _ := fetcher.FetchWithExpires(ctx, expiresDuration)

	if logger == nil {
		logger = NopLogger
	}

	if fallbackDuration < 0 {
		fallbackDuration = expiresDuration
	}

	if time.Since(t) > fallbackDuration {
		logger.Warn("FetchDataExpired",
			zap.Duration("FallbackDuration", fallbackDuration),
			zap.Time("FetchingTime", t),
		)

		return
	}

	return data, true
}

func (fetcher *CachedDataFetcher[Data]) FetchWithSince(ctx context.Context, since time.Time) (data Data, t time.Time, err error) {
	data, t, ok := fetcher.getWithSince(since)
	if ok {
		return
	}

	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	data, t, ok = fetcher.getWithSince(since)
	if ok {
		return
	}

	data, t, err = fetcher.refetch(ctx, since)
	if err == nil {
		return
	}

	data, t = fetcher.getCached()
	return
}

func (fetcher *CachedDataFetcher[Data]) GetRefresh() func(ctx context.Context) (Data, time.Time, error) {
	return fetcher.Refresh
}

func (fetcher *CachedDataFetcher[Data]) Refresh(ctx context.Context) (data Data, t time.Time, err error) {
	return fetcher.RefreshWithSinceTime(ctx, time.Now())
}

func (fetcher *CachedDataFetcher[Data]) RefreshLowerCache(ctx context.Context) (data Data, t time.Time, err error) {
	return fetcher.RefreshWithSinceTime(ctx, time.Time{})
}

func (fetcher *CachedDataFetcher[Data]) RefreshWithSinceTime(ctx context.Context, since time.Time) (data Data, t time.Time, err error) {
	fetcher.mutex.Lock()
	defer fetcher.mutex.Unlock()

	return fetcher.refetch(ctx, since)
}

func (fetcher *CachedDataFetcher[Data]) TriggerRefreshLowerCache(ctx context.Context) {
	fetcher.RefreshLowerCache(ctx)
}

func (fetcher *CachedDataFetcher[Data]) refetch(ctx context.Context, sinceTime time.Time) (data Data, t time.Time, err error) {
	logger := fetcher.With(
		zap.Time("SinceTime", sinceTime),
		zap.Time("CurrentDataTime", fetcher.data.Time()),
	)

	logger.Info("Refetching")

	data, t, err = fetcher.fetcher(ctx, sinceTime)
	if err != nil {
		logger.Error("RefetchFailed",
			zap.Error(err),
		)
		return
	}

	logger.Info("Refetched",
		zap.Time("FetchedTime", t),
	)

	if t.IsZero() {
		t = time.Now()
	}

	fetcher.data = NewTimedValue(data, t)

	// call it asynchronously to avoid dead lock
	// in case that callbacks may call fetcher method again
	go fetcher.callbacks.Call(data, t)

	return
}

func (fetcher *CachedDataFetcher[Data]) getWithExpires(expiresDuration time.Duration) (Data, time.Time, bool) {
	return fetcher.getWithSince(fetcher.SinceTimeFromExpiresDuration(expiresDuration))
}

func (fetcher *CachedDataFetcher[Data]) getWithSince(since time.Time) (data Data, t time.Time, ok bool) {
	data, t = fetcher.getCached()
	ok = t.After(since)
	return
}

func (fetcher *CachedDataFetcher[Data]) getCached() (data Data, t time.Time) {
	return fetcher.data.ValueAndTime()
}

type CachedDataFetcherAccessor[Data any] struct {
	fetcher          *CachedDataFetcher[Data]
	expiresDuration  time.Duration
	fallbackDuration time.Duration
	logger           *zap.Logger
	logExpired       bool
	ctx              context.Context
}

func (getter *CachedDataFetcherAccessor[Data]) Dup() *CachedDataFetcherAccessor[Data] {
	return Dup(getter)
}

func (getter *CachedDataFetcherAccessor[Data]) Get() (data Data, ok bool) {
	if getter == nil {
		return
	}

	data, fetchingTime, ok := getter.fetcher.GetWithExpires(getter.expiresDuration)
	if !ok {
		if getter.logExpired {
			getter.GetLogger().Warn("GetWithExpires",
				zap.Time("FetchingTime", fetchingTime),
				zap.Duration("ExpiresDuration", getter.expiresDuration),
			)
		}

		if getter.fallbackDuration > 0 {
			if time.Since(fetchingTime) < getter.expiresDuration {
				ok = true
			}
		}

		return
	}

	return
}

func (getter *CachedDataFetcherAccessor[Data]) GetLogger() *zap.Logger {
	if getter == nil {
		return zap.NewNop()
	}

	logger := getter.logger
	if logger != nil {
		return logger
	}

	return getter.fetcher.Logger
}

func (getter *CachedDataFetcherAccessor[Data]) WithExpiresDuration(expiresDuration time.Duration) *CachedDataFetcherAccessor[Data] {
	if getter == nil {
		return nil
	}

	getter.expiresDuration = expiresDuration
	return getter
}

func (getter *CachedDataFetcherAccessor[Data]) WithFallbackDuration(fallbackDuration time.Duration) *CachedDataFetcherAccessor[Data] {
	if getter == nil {
		return nil
	}

	getter.fallbackDuration = fallbackDuration
	return getter
}

func (getter *CachedDataFetcherAccessor[Data]) WithLogExpired(logExpired bool) *CachedDataFetcherAccessor[Data] {
	if getter == nil {
		return nil
	}

	getter.logExpired = logExpired
	return getter
}

func (getter *CachedDataFetcherAccessor[Data]) WithLogger(logger *zap.Logger) *CachedDataFetcherAccessor[Data] {
	if getter == nil {
		return nil
	}

	getter.logger = logger
	return getter
}

func NewCachedDataFetcherFromAnother[Data any, BasedData any](basedFetcher *CachedDataFetcher[BasedData], newDataFromBased func(BasedData) Data) *CachedDataFetcher[Data] {
	return NewCachedDataFetcher(func(ctx context.Context, since time.Time) (data Data, t time.Time, err error) {
		var basedData BasedData
		basedData, t, err = basedFetcher.FetchWithSince(ctx, since)
		if err != nil {
			return
		}

		data = newDataFromBased(basedData)

		return
	}).WithOthersSubscribed(0, basedFetcher)
}
