/*
 * Author: fasion
 * Created time: 2023-12-22 13:09:37
 * Last Modified by: fasion
 * Last Modified time: 2025-01-02 10:42:16
 */

package stl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAutoRefresh(t *testing.T) {
	fetcher := NewCachedDataFetcher(func(ctx context.Context, since time.Time) (result any, t time.Time, err error) {
		fmt.Println("fetch data")
		return nil, time.Now(), nil
	}).WithExpiresDuration(time.Second * 2)

	fetcher.Fetch(nil)
	lastFetchTime := fetcher.data.t
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))

	fetcher.Fetch(nil)
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	time.Sleep(time.Second)

	fetcher.Fetch(nil)
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	fetcher.FetchWithExpires(nil, time.Second)
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))
	assert.NotEqual(t, lastFetchTime, fetcher.data.t)
	lastFetchTime = fetcher.data.t

	time.Sleep(time.Second + time.Second/2)

	fetcher.Fetch(nil)
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	time.Sleep(time.Second)
	fetcher.Fetch(nil)
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))
	assert.NotEqual(t, lastFetchTime, fetcher.data.t)
	lastFetchTime = fetcher.data.t

	time.Sleep(time.Second)
	fetcher.FetchWithSince(nil, time.Now().Add(-2*time.Second))
	fmt.Println("LastFetchTime", fetcher.data.t.Format("15:04:05.999"))
	assert.Equal(t, lastFetchTime, fetcher.data.t)
}

func TestCachedFetcherWithExpires(t *testing.T) {
	fetcher := NewCachedDataFetcher(func(ctx context.Context, since time.Time) (result any, t time.Time, err error) {
		return nil, time.Now(), nil
	}).WithExpiresDuration(time.Second * 2)

	fetcher.Fetch(nil)
	assert.True(t, !fetcher.data.t.IsZero())
	lastFetchTime := fetcher.data.t

	fetcher.Fetch(nil)
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	fetcher.FetchWithExpires(nil, time.Second)
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	fetcher.FetchWithExpires(nil, time.Second*5)
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	time.Sleep(time.Second * 3)
	fetcher.FetchWithExpires(nil, time.Second*5)
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	fetcher.Fetch(nil)
	assert.NotEqual(t, lastFetchTime, fetcher.data.t)
	lastFetchTime = fetcher.data.t

	fetcher.Fetch(nil)
	assert.Equal(t, lastFetchTime, fetcher.data.t)

	fetcher.FetchWithExpires(nil, time.Nanosecond)
	assert.NotEqual(t, lastFetchTime, fetcher.data.t)
}

func TestCacherFetcherCallback(t *testing.T) {
	fetcher := NewCachedDataFetcher(func(ctx context.Context, since time.Time) (result any, t time.Time, err error) {
		return nil, time.Now(), nil
	})

	var one bool
	var two bool

	fetcher.NewCallbackLite(func(_ context.Context) {
		one = true
	}, time.Hour)

	fetcher.NewCallbackLite(func(_ context.Context) {
		two = true
	}, time.Hour)

	fetcher.TriggerRefreshLowerCache(nil)

	time.Sleep(time.Second)
	assert.True(t, one)
	assert.True(t, two)
}

func TestCachedFetcherSubscribeOthersCompile(t *testing.T) {
	if false {
		var a CachedDataFetcher[int]
		var b CachedDataFetcher[string]
		var c CachedDataFetcher[error]

		NewCachedDataFetcher[bool](nil).WithOthersSubscribed(0, &a, &b, &c)
	}
}

func TestCachedDataFetcherGetter(t *testing.T) {
	fetcher := NewCachedDataFetcher(func(ctx context.Context, since time.Time) (result any, t time.Time, err error) {
		fmt.Println("fetch data")
		return nil, time.Now(), nil
	}).WithExpiresDuration(time.Second * 10)

	time.Sleep(time.Second * 2)

	fmt.Println(fetcher.GetWithExpiresWarn(time.Second))

	getter := fetcher.BuildAccessor().WithLogger(zap.L()).WithExpiresDuration(time.Second * 1).WithLogExpired(true).Get
	fmt.Println(getter())
}

func TestCustomExpireDuration(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
		return
	}

	basic := NewCachedDataFetcher(func(ctx context.Context, since time.Time) (time.Time, time.Time, error) {
		now := time.Now()
		return now, now, nil
	}).
		WithExpiresDuration(time.Second * 3).
		WithLogger(logger.Named("basic"))

	cacher := NewCachedDataFetcherFromAnother(basic, Echo[time.Time]).
		WithExpiresDuration(time.Second * 3).
		WithLogger(logger.Named("cacher"))

	data1, ok := cacher.FetchWithExpiresPro(nil, time.Second, time.Second*3, logger)
	fmt.Println("data1", data1)
	assert.Equal(t, ok, true)
	assert.Less(t, time.Since(data1), time.Millisecond)

	time.Sleep(time.Millisecond * 100)
	fmt.Println("data1", data1)

	data2, ok := cacher.FetchWithExpiresPro(nil, time.Second, time.Second*3, logger)
	fmt.Println("data2", data2)
	assert.Equal(t, ok, true)
	assert.Equal(t, data1, data2)

	time.Sleep(time.Second)
	fmt.Println("before data3", time.Now())

	data3, ok := cacher.FetchWithExpiresPro(nil, time.Second, time.Second*3, logger)
	fmt.Println("data3", data3)
	assert.Equal(t, ok, true)
	assert.NotEqual(t, data3, data2)
}
