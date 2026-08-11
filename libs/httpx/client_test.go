/*
 * Author: fasion
 * Created time: 2024-12-11 10:14:19
 * Last Modified by: fasion
 * Last Modified time: 2025-08-14 14:36:16
 */

package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/fasionchan/goutils/baseutils/netutils"
	"github.com/fasionchan/goutils/std/_testing"
	"github.com/stretchr/testify/assert"
)

func TestHttpClientJson(t *testing.T) {
	client, err := NewClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
		return
	}

	var data any
	if _, err := client.JsonGet(context.Background(), "", nil, nil, nil, &data); err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(data)
}

func TestHttpClientJsonAsMap(t *testing.T) {
	client, err := NewClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
		return
	}

	data, _, err := JsonGet[map[string]any](context.Background(), client, "", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(data)
}

func TestHttpClientJsonAsMapPtr(t *testing.T) {
	client, err := NewClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
		return
	}

	data, _, err := JsonGet[*map[string]any](context.Background(), client, "", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(data)
}

func TestHttpClientJsonAsMapSecondaryPtr(t *testing.T) {
	client, err := NewClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
		return
	}

	data, _, err := JsonGet[**map[string]any](context.Background(), client, "", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(data)
	fmt.Println(*data)
	fmt.Println(**data)
}

func TestHttpClientJsonParseError(t *testing.T) {
	client, err := NewClient("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
		return
	}

	var data any
	if _, err := client.JsonGet(context.Background(), "/", nil, nil, nil, &data); err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(data)
}

func TestHttpClientTestConnectTcp(t *testing.T) {
	for _, testCase := range []struct {
		url string
		ok  bool
	}{
		{
			url: "http://www.baidu.com",
			ok:  true,
		},
		{
			url: "http://www.baidu.com:80",
			ok:  true,
		},
		{
			url: "https://www.baidu.com",
			ok:  true,
		},
		{
			url: "https://www.baidu.com:443",
			ok:  true,
		},
		{
			url: "http://127.0.0.1:65432",
			ok:  false,
		},
		{
			url: "http://nosuch.nosuch.nosuch",
			ok:  false,
		},
	} {
		client, err := NewClient(testCase.url)
		if err != nil {
			t.Fatal(err)
			return
		}

		err = client.TestConnectTcp()
		fmt.Println(err)

		assert.Equal(t, testCase.ok, err == nil)
	}
}

func TestHttpClientHeader(t *testing.T) {
	client, err := NewBlankHttpClient().WithHeadersFromMap(map[string]string{
		"X": "X",
		"Y": "Y",
	}).WithCookiesFromMap(map[string]string{
		"X": "X",
		"Y": "Y",
	}).WithQueryFromMap(map[string]string{
		"a": "a",
		"b": "b",
	}).WithRawBaseUrl("https://httpbin.org/")
	if err != nil {
		t.Fatal(err)
		return
	}

	var result map[string]any
	response, err := client.JsonGet(context.Background(), "/get",
		url.Values{
			"b": []string{"bb"},
			"c": []string{"c"},
		},
		netutils.NewHeaderFromMap(map[string]string{
			"Y": "YY",
			"Z": "Z",
		}),
		netutils.NewCookiesFromMap(map[string]string{
			"Y": "YY",
			"Z": "Z",
		}),
		&result,
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(response)
	fmt.Println()

	fmt.Println(result)
}

type HttpClientBuildAndDoRangeRequestTestCase struct {
	url string
	start int64
	end int64

	expectedStart int64
	expectedEnd int64
	expectedTotal int64
	expectedHasMore bool
	expectedData []byte
	expectedError bool
}

func (testCase *HttpClientBuildAndDoRangeRequestTestCase) GetName() string {
	return fmt.Sprintf("TestHttpClientBuildAndDoRangeRequest-[%d-%d]", testCase.start, testCase.end)
}

func (testCase *HttpClientBuildAndDoRangeRequestTestCase) Run(t *testing.T) {
	start, end, total, hasMore, data, err := testCase.fetch()
	fmt.Println(start, end, total, hasMore, data, err)

	assert.Equal(t, testCase.expectedStart, start, "start")
	assert.Equal(t, testCase.expectedEnd, end, "end")
	assert.Equal(t, testCase.expectedTotal, total, "total")
	assert.Equal(t, testCase.expectedHasMore, hasMore, "hasMore")
	assert.Equal(t, testCase.expectedData, data, "data")
	assert.Equal(t, testCase.expectedError, err != nil, fmt.Sprintf("error: %v", err))
}

func (testCase *HttpClientBuildAndDoRangeRequestTestCase) fetch() (start, end int64, total int64, hasMore bool, data []byte, err error) {
	client, err := NewClient("")
	if err != nil {
		return
	}

	start, end, total, hasMore, response, err := client.BuildAndDoRangeRequestWithBody(context.Background(), "GET", testCase.url, nil, nil, nil, nil, "", testCase.start, testCase.end)
	if err != nil {
		return
	}

	if response != nil {
		defer response.Body.Close()

		data, err = io.ReadAll(response.Body)
		if err != nil {
			return
		}
	}

	return
}

func TestHttpClientBuildAndDoRangeRequest(t *testing.T) {
	url := "https://httpbin.org/range/10240"
	response, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
		return
	}

	n := int64(len(data))

	_testing.RunNamedTestCasesX(t, 
		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: -2,
			end: -1,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: -1,
			end: -1,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: -1,
			end: 0,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: 0,
			end: 0,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: 0,
			end: 1,
			expectedStart: 0,
			expectedEnd: 1,
			expectedTotal: n,
			expectedHasMore: true,
			expectedData: data[:1],
			expectedError: false,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: 1,
			end: 1,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: 1,
			end: 2,
			expectedStart: 1,
			expectedEnd: 2,
			expectedTotal: n,
			expectedHasMore: true,
			expectedData: data[1:2],
			expectedError: false,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n-2,
			end: n-1,
			expectedStart: n-2,
			expectedEnd: n-1,
			expectedTotal: n,
			expectedHasMore: true,
			expectedData: data[n-2:n-1],
			expectedError: false,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n-1,
			end: n-1,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n-1,
			end: n,
			expectedStart: n-1,
			expectedEnd: n,
			expectedTotal: n,
			expectedHasMore: false,
			expectedData: data[n-1:],
			expectedError: false,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n,
			end: n,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n,
			end: n+1,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n+1,
			end: n+1,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},

		&HttpClientBuildAndDoRangeRequestTestCase{
			url: url,
			start: n+1,
			end: n+2,
			expectedStart: 0,
			expectedEnd: 0,
			expectedTotal: 0,
			expectedHasMore: false,
			expectedData: nil,
			expectedError: true,
		},
	)
}