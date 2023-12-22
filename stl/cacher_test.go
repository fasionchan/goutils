/*
 * Author: fasion
 * Created time: 2023-12-22 13:09:37
 * Last Modified by: fasion
 * Last Modified time: 2023-12-22 16:16:44
 */

package stl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAutoRefresh(t *testing.T) {
	fetcher := NewCachedDataFetcher(func(ctx context.Context) (result any, t time.Time, err error) {
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
	fetcher := NewCachedDataFetcher(func(ctx context.Context) (result any, t time.Time, err error) {
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
	fetcher := NewCachedDataFetcher(func(ctx context.Context) (result any, t time.Time, err error) {
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

	fetcher.TriggerRefresh(nil)

	time.Sleep(time.Second)
	assert.True(t, one)
	assert.True(t, two)
}

func TestCachedFetcherSubscribeOthersCompile(t *testing.T) {
	if false {
		var a CachedDataFetcher[int]
		var b CachedDataFetcher[string]
		var c CachedDataFetcher[error]

		NewCachedDataFetcher[bool](nil).SubscribeOthers(0, &a, &b, &c)
	}
}
