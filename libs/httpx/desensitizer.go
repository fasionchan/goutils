/*
 * Author: fasion
 * Created time: 2025-09-16 11:47:46
 * Last Modified by: fasion
 * Last Modified time: 2025-09-16 12:35:56
 */

package httpx

import "net/http"

type RequestDesensitizerFunc func(*http.Request) (*http.Request, error)

func (fn RequestDesensitizerFunc) DesensitizeRequest(request *http.Request) (*http.Request, error) {
	if fn == nil {
		return request, nil
	}

	return fn(request)
}

type ResponseDesensitizerFunc func(*http.Response) (*http.Response, error)

func (fn ResponseDesensitizerFunc) DesensitizeResponse(response *http.Response) (*http.Response, error) {
	if fn == nil {
		return response, nil
	}

	return fn(response)
}

type RequestDesensitizer interface {
	DesensitizeRequest(request *http.Request) (*http.Request, error)
}

type ResponseDesensitizer interface {
	DesensitizeResponse(response *http.Response) (*http.Response, error)
}

type RequestResponseDesensitizer interface {
	RequestDesensitizer
	ResponseDesensitizer
}

func NewRequestResponseDesensitizer(requestDesensitizer RequestDesensitizer, responseDesensitizer ResponseDesensitizer) RequestResponseDesensitizer {
	return struct {
		RequestDesensitizer
		ResponseDesensitizer
	}{
		RequestDesensitizer:  requestDesensitizer,
		ResponseDesensitizer: responseDesensitizer,
	}
}

func NewRequestResponseDesensitizerFromFunc(requestDesensitizer RequestDesensitizerFunc, responseDesensitizer ResponseDesensitizerFunc) RequestResponseDesensitizer {
	return NewRequestResponseDesensitizer(requestDesensitizer, responseDesensitizer)
}
