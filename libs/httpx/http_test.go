/*
 * Author: fasion
 * Created time: 2025-03-28 15:55:58
 * Last Modified by: fasion
 * Last Modified time: 2025-07-24 14:05:16
 */

package httpx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/fasionchan/goutils/stl"
	"github.com/stretchr/testify/assert"
)

func TestClientRequestBodyClose(t *testing.T) {
	var hasClose bool

	buffer := bytes.NewBufferString(`{"message":"hello"}`)
	body := struct {
		*bytes.Buffer
		stl.CloseFunc
	}{
		Buffer: buffer,
		CloseFunc: func() error {
			fmt.Println("close called!!!")
			hasClose = true
			return nil
		},
	}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://httpbin.org/post", body)
	if err != nil {
		t.Fatal(err)
		return
	}

	var client http.Client
	reponse, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
		return
	}

	if _, err := io.Copy(io.Discard, reponse.Body); err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, hasClose, true)
}

func TestBodyReadTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/delay/10", nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	assert.Error(t, err)
}

func TestBaseAuthByUrl(t *testing.T) {
	client, err := NewClient("https://user:xxxx@httpbin.org/get")
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, client.baseUrl.User, (*url.Userinfo)(nil))
	response, err := client.JsonGet(context.Background(), ".", nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	assert.Equal(t, response.StatusCode, 200)
	fmt.Println(response.StatusCode)
}
