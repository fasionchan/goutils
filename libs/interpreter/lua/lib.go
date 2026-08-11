package lua

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Json(v any) (string, error) {
	content, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Unjson(content string) (v any, err error) {
	err = json.Unmarshal([]byte(content), &v)
	return v, err
}

func Http(url, method string, query url.Values, headers http.Header, cookies url.Values, body string) (statusCode int, respHeaders http.Header, respBody []byte, err error) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return 0, nil, nil, err
	}

	request.URL.RawQuery = query.Encode()
	request.Header = headers

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, nil, err
	}

	statusCode = response.StatusCode
	respHeaders = response.Header

	respBody, err = io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	return
}