package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/baseutils/netutils"
	"github.com/fasionchan/goutils/libs/logging"
	"github.com/fasionchan/goutils/std/templatex"
	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
	"go.uber.org/zap"
)

const (
	HttpHeaderContentType                    = "Content-Type"
	ContentTypeApplicationJson               = "application/json"
	ContentTypeMultipartFormData             = "multipart/form-data"
	ContentTypeApplicationXwwwFormUrlencoded = "application/x-www-form-urlencoded"

	QueryParameterPrefixQuery             = "_query_"
	QueryParameterPrefixCookie            = "_cookie_"
	QueryParameterPrefixHeader            = "_header_"
	QueryParameterPrefixLengthCookie      = len(QueryParameterPrefixCookie)
	QueryParameterNameBearerToken         = "_bearerToken"
	QueryParameterNameHmacAuthAccessKey   = "_hmacAuthAccessKey"
	QueryParameterNameHmacAuthSecretKey   = "_hmacAuthSecretKey"
	QueryParameterNameHmacAuthHashMethod  = "_hmacAuthHashMethod"
	QueryParameterNameHmacAuthPaths       = "_hmacAuthPaths"
	QueryParameterNameJwtPathTokenSecret  = "_jwtPathTokenSecret"
	QueryParameterNameJwtPathTokenSeconds = "_jwtPathTokenSeconds"

	DefaultClientName                 = "HttpxClient"
	DefaultJwtPathTokenExpireDuration = 60 * time.Second
)

type ResponseValidator = func(*http.Response) error

type Client = HttpClient

type HttpClient struct {
	*zap.Logger
	*http.Client

	baseUrl *url.URL
	query   url.Values
	headers http.Header
	cookies []*http.Cookie

	urlPrefixesReplacer PrefixesReplacer

	bearerToken string

	requestAuthenticator RequestAuthenticator
	responseValidator    ResponseValidator

	requestDesensitizer  RequestDesensitizer
	responseDesensitizer ResponseDesensitizer

	sensitiveCookieNames types.Strings
	sensitiveHeaderNames types.Strings

	serverSentEventEndOfLine string
}

func NewClient(baseUrl string) (*Client, error) {
	return NewBlankHttpClient().WithRawBaseUrl(baseUrl)
}

func NewClientWithBaseUrl(baseUrl *url.URL) *Client {
	return NewBlankHttpClient().WithBaseUrl(baseUrl)
}

func NewBlankHttpClient() *Client {
	return &Client{
		Client: http.DefaultClient,
		Logger: logging.GetLogger().Named("HttpxClient"),
	}
}

func (client *Client) LoggerFromContext(ctx context.Context) (*logging.LoggerRef, context.Context) {
	return logging.LoggerRefFromContextPro(ctx, true, true, DefaultClientName, client.Logger, logging.GetNopLogger())
}

func (client *Client) BearerTokenAuthHeader() http.Header {
	if client == nil {
		return nil
	}

	if client.bearerToken == "" {
		return nil
	}

	return http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", client.bearerToken)},
	}
}

func (client *Client) BuiltinHeaders() http.Header {
	return stl.ConcatMaps(client.BearerTokenAuthHeader())
}

func (client *Client) DeepDup() *Client {
	if client == nil {
		return nil
	}

	// todo deepdup other fields
	return stl.Dup(client).
		WithCookies(stl.DupSlice(client.cookies)).
		WithHeaders(stl.DupMap(client.headers))
}

func (client *Client) DesensitizeHeader(header http.Header) http.Header {
	sensitiveNames := client.GetSensitiveHeaderNames()
	if sensitiveNames.Empty() {
		return header
	}

	return stl.PurgeMapKeys(stl.DupMap(header), sensitiveNames...)
}

func (client *Client) Dup() *Client {
	return stl.Dup(client)
}

// GetXxxx 获取内部属性

func (client *Client) GetBaseUrl() *url.URL {
	if client == nil {
		return nil
	}
	return client.baseUrl
}

func (client *Client) GetBaseUrlString() string {
	baseUrl := client.GetBaseUrl()
	if baseUrl == nil {
		return ""
	}
	return baseUrl.String()
}

func (client *Client) GetSchemeHostUrl() string {
	baseUrl := client.GetBaseUrl()
	if baseUrl == nil {
		return ""
	}

	return fmt.Sprintf("%s://%s", baseUrl.Scheme, baseUrl.Host)
}

func (client *Client) GetBearerToken() string {
	if client == nil {
		return ""
	}
	return client.bearerToken
}

func (client *Client) GetHeaders() http.Header {
	if client == nil {
		return nil
	}

	return client.headers
}

func (client *Client) GetCookies() []*http.Cookie {
	if client == nil {
		return nil
	}

	return client.cookies
}

func (client *Client) GetQuery() url.Values {
	if client == nil {
		return nil
	}
	return client.query
}

func (client *Client) GetSensitiveHeaderNames() types.Strings {
	if client == nil {
		return nil
	}
	return client.sensitiveHeaderNames
}

func (client *Client) GetSensitiveCookieNames() types.Strings {
	if client == nil {
		return nil
	}
	return client.sensitiveCookieNames
}

func (client *Client) GetUrlPrefixesReplacer() PrefixesReplacer {
	if client == nil {
		return nil
	}
	return client.urlPrefixesReplacer
}

// WithXxxx 设置内部属性

func (client *Client) WithRawBaseUrl(rawBaseUrl string) (*Client, error) {
	baseUrl, err := url.Parse(rawBaseUrl)
	if err != nil {
		return nil, err
	}

	client.WithBaseUrl(baseUrl)

	if err := client.ParseJwtPathTokenAuthenticatorFromBaseUrl(); err != nil {
		return nil, err
	}

	return client, nil
}

func (client *Client) WithUrlPrefixesReplacer(urlPrefixesReplacer PrefixesReplacer) *Client {
	if client == nil {
		return nil
	}

	client.urlPrefixesReplacer = urlPrefixesReplacer
	return client
}

func (client *Client) WithBaseUrl(baseUrl *url.URL) *Client {
	if client == nil {
		return nil
	}

	client.baseUrl = baseUrl

	if baseUrl != nil {
		client.ParseBasicAuthFromBaseUrl()
		client.ParseBearerTokenFromBaseUrl()
		client.ParseCookiesFromBaseUrl()
		client.ParseHmacAuthenticatorFromBaseUrl()
		client.ParseHeadersFromBaseUrl()
		client.ParseQueryFromBaseUrl()
	}

	return client
}

func (client *Client) WithServerSentEventEndOfLine(endOfLine string) *Client {
	if client == nil {
		return nil
	}

	client.serverSentEventEndOfLine = endOfLine

	return client
}

func (client *Client) ParseBasicAuthFromBaseUrl() *Client {
	if client == nil {
		return nil
	}

	baseUrl := client.baseUrl
	if baseUrl == nil {
		return client
	}

	user := baseUrl.User
	if user == nil {
		return client
	}

	baseUrl.User = nil
	client.WithRequestAuthenticator(NewBasicAuthenticatorFromUrlUserinfo(user))

	return client
}

func (client *Client) ParseBearerTokenFromBaseUrl() *Client {
	if client == nil {
		return nil
	}

	baseUrl := client.baseUrl
	if baseUrl == nil {
		return client
	}

	query := baseUrl.Query()
	if len(query) == 0 {
		return client
	}

	token, ok := query[QueryParameterNameBearerToken]
	if !ok {
		return client
	}

	var bearerToken string
	if len(token) > 0 {
		bearerToken = token[0]
	}

	client.bearerToken = bearerToken

	query.Del(QueryParameterNameBearerToken)
	baseUrl.RawQuery = query.Encode()

	return client
}

func (client *Client) ParseCookiesFromBaseUrl() *Client {
	if client == nil {
		return nil
	}

	client.cookies = NewCookiesFromUrlValues(client.RemovePrefixedQueriesFromBaseUrl(QueryParameterPrefixCookie))

	return client
}

func (client *Client) ParseHeadersFromBaseUrl() *Client {
	if client == nil {
		return nil
	}

	client.headers = http.Header(client.RemovePrefixedQueriesFromBaseUrl(QueryParameterPrefixHeader))
	return client
}

func (client *Client) ParseQueryFromBaseUrl() *Client {
	if client == nil {
		return nil
	}

	client.query = client.RemovePrefixedQueriesFromBaseUrl(QueryParameterPrefixQuery)
	return client
}

func (client *Client) RemovePrefixedQueriesFromBaseUrl(prefix string) url.Values {
	if client == nil {
		return nil
	}

	baseUrl := client.baseUrl
	if baseUrl == nil {
		return nil
	}

	query := baseUrl.Query()
	if len(query) == 0 {
		return nil
	}

	prefixLen := len(prefix)
	result := url.Values{}

	for name, values := range stl.DupMap(query) {
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		unprefixedName := name[prefixLen:]
		if unprefixedName == "" {
			continue
		}

		for _, value := range values {
			result.Add(unprefixedName, value)
		}

		query.Del(name)
	}

	baseUrl.RawQuery = query.Encode()

	return result
}

func (client *Client) ParseHmacAuthenticatorFromBaseUrl() *Client {
	if client == nil {
		return nil
	}

	baseUrl := client.baseUrl
	if baseUrl == nil {
		return client
	}

	query := baseUrl.Query()
	authenticator := NewHmacAuthenticatorsFromUrlValues(query).RequestAuthenticator()
	if authenticator == nil {
		return client
	}

	query.Del(QueryParameterNameHmacAuthAccessKey)
	query.Del(QueryParameterNameHmacAuthSecretKey)
	query.Del(QueryParameterNameHmacAuthHashMethod)
	query.Del(QueryParameterNameHmacAuthPaths)

	baseUrl.RawQuery = query.Encode()

	return client.WithRequestAuthenticator(authenticator)
}

func (client *Client) ParseJwtPathTokenAuthenticatorFromBaseUrl() error {
	if client == nil {
		return nil
	}

	baseUrl := client.baseUrl
	if baseUrl == nil {
		return nil
	}

	query := baseUrl.Query()
	authenticator, err := NewJwtPathTokenAuthenticatorFromUrlValues(query)
	if err != nil {
		return err
	}
	if authenticator == nil {
		return nil
	}

	query.Del(QueryParameterNameJwtPathTokenSecret)
	query.Del(QueryParameterNameJwtPathTokenSeconds)
	baseUrl.RawQuery = query.Encode()

	client.WithRequestAuthenticator(authenticator)
	return nil
}

func (client *Client) WithBearerToken(bearerToken string) *Client {
	if client == nil {
		return nil
	}

	client.bearerToken = bearerToken

	return client
}

func (client *Client) WithClient(httpClient *http.Client) *Client {
	if client == nil {
		return nil
	}
	client.Client = httpClient
	return client
}

func (client *Client) WithCookiesFromMap(cookies map[string]string) *Client {
	return client.WithCookies(netutils.NewCookiesFromMap(cookies))
}

func (client *Client) WithCookies(cookies []*http.Cookie) *Client {
	if client == nil {
		return nil
	}
	client.cookies = cookies
	return client
}

func (client *Client) WithCookiesX(cookies ...*http.Cookie) *Client {
	return client.WithCookies(cookies)
}

func (client *Client) WithHeader(name, value string) *Client {
	if client == nil {
		return nil
	}

	if client.headers == nil {
		client.headers = http.Header{}
	}

	client.headers.Add(name, value)

	return client
}

func (client *Client) SetHeaders(headers http.Header) *Client {
	if client == nil {
		return nil
	}

	if client.headers == nil {
		client.headers = http.Header{}
	}

	stl.ConcatMapsTo(client.headers, headers)

	return client
}

func (client *Client) SetHeader(name, value string) *Client {
	if client == nil {
		return nil
	}

	if client.headers == nil {
		client.headers = http.Header{}
	}

	client.headers.Set(name, value)

	return client
}

func (client *Client) WithHeadersFromMap(headers map[string]string) *Client {
	return client.WithHeaders(netutils.NewHeaderFromMap(headers))
}

func (client *Client) WithHeaders(headers http.Header) *Client {
	if client == nil {
		return nil
	}
	client.headers = headers
	return client
}

func (client *Client) WithLogger(logger *zap.Logger) *Client {
	if client == nil {
		return nil
	}
	client.Logger = logger
	return client
}

func (client *Client) WithLoggerNamed(loggerName string) *Client {
	return client.WithLogger(logging.GetLogger().Named(loggerName))
}

func (client *Client) WithQueryFromMap(query map[string]string) *Client {
	return client.WithQuery(netutils.NewUrlValuesFromMap(query))
}

func (client *Client) WithQuery(query url.Values) *Client {
	if client == nil {
		return nil
	}
	client.query = query
	return client
}

func (client *Client) WithRequestAuthenticator(authenticator RequestAuthenticator) *Client {
	if client == nil {
		return nil
	}

	client.requestAuthenticator = authenticator
	return client
}

func (client *Client) WithRequestAuthenticatorFunc(fn RequestAuthenticatorFunc) *Client {
	return client.WithRequestAuthenticator(fn)
}

func (client *Client) WithResponseValidatorDefault() *Client {
	return client.WithResponseValidatorFor200()
}

func (client *Client) WithResponseValidatorFor200() *Client {
	return client.WithResponseValidator(ValidateResponseFor200)
}

func (client *Client) WithResponseValidatorFor2xx() *Client {
	return client.WithResponseValidator(ValidateResponseFor2xx)
}

func (client *Client) WithResponseValidator(validator ResponseValidator) *Client {
	if client == nil {
		return nil
	}
	client.responseValidator = validator
	return client
}

func (client *Client) WithRequestDesensitizer(desensitizer RequestDesensitizer) *Client {
	if client == nil {
		return nil
	}
	client.requestDesensitizer = desensitizer
	return client
}

func (client *Client) WithRequestDesensitizerFunc(fn RequestDesensitizerFunc) *Client {
	return client.WithRequestDesensitizer(fn)
}

func (client *Client) WithResponseDesensitizer(desensitizer ResponseDesensitizer) *Client {
	if client == nil {
		return nil
	}
	client.responseDesensitizer = desensitizer
	return client
}

func (client *Client) WithResponseDesensitizerFunc(fn ResponseDesensitizerFunc) *Client {
	return client.WithResponseDesensitizer(fn)
}

func (client *Client) WithRequestResponseDesensitizer(desensitizer RequestResponseDesensitizer) *Client {
	if client == nil {
		return nil
	}

	client.requestDesensitizer = desensitizer
	client.responseDesensitizer = desensitizer
	return client
}

func (client *Client) WithSensitiveCookieNames(names ...string) *Client {
	if client == nil {
		return nil
	}

	client.sensitiveCookieNames = client.sensitiveCookieNames.Append(names...).Unique()
	return client
}

func (client *Client) WithSensitiveHeaderNames(names ...string) *Client {
	if client == nil {
		return nil
	}

	client.sensitiveHeaderNames = client.sensitiveHeaderNames.Append(names...).Unique()
	return client
}

func (client *Client) JsonGet(ctx context.Context, url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any) (*http.Response, error) {
	return client.JsonRequest(ctx, http.MethodGet, url, query, headers, cookies, nil, data)
}

func (client *Client) JsonPost(ctx context.Context, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body, data any) (*http.Response, error) {
	return client.JsonRequest(ctx, http.MethodPost, url, query, headers, cookies, body, data)
}

// JsonRequest 是 BuildAndDoRequestXxxxForJson 的快捷方式

func (client *Client) JsonRequest(ctx context.Context, method, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body, data any) (response *http.Response, err error) {
	if _, specified := RetryIfContextSpecified(ctx, func() error {
		response, err = client.JsonRequest(ctx, method, url, query, headers, cookies, body, data)
		return err
	}); specified {
		return
	}

	return client.BuildAndDoRequestWithJsonForJson(ctx, method, url, query, headers, cookies, body, data)
}

func (client *Client) JsonRequestWithForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm, data any) (*http.Response, error) {
	return client.BuildAndDoRequestWithFormForJson(ctx, method, _url, query, headers, cookies, form, data)
}

func (client *Client) JsonRequestWithStreamingForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm, data any) (*http.Response, error) {
	return client.BuildAndDoRequestWithStreamingFormForJson(ctx, method, _url, query, headers, cookies, form, data)
}

func (client *Client) JsonRequestWithUrlEncodedForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form url.Values, data any) (*http.Response, error) {
	return client.BuildAndDoRequestWithUrlEncodedFormForJson(ctx, method, _url, query, headers, cookies, form, data)
}

func (client *Client) JsonRequestWithBodyBytes(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body []byte, contentType string, data any) (*http.Response, error) {
	return client.JsonRequestWithBody(ctx, method, _url, query, headers, cookies, bytes.NewReader(body), contentType, data)
}

func (client *Client) JsonRequestWithBody(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, contentType string, data any) (*http.Response, error) {
	return client.BuildAndDoRequestWithBodyForJson(ctx, method, _url, query, headers, cookies, body, contentType, data)
}

// BuildAndDoRequestWithXxxxForYyyy
// 构建请求、执行请求、处理响应

func (client *Client) BuildAndDoRequestWithJsonForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithDataForResult(ctx, method, _url, query, headers, cookies, data, ContentTypeApplicationJson, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithJsonForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, result any, resultType string) (*http.Response, error) {
	return client.BuildAndDoRequestWithDataForResult(ctx, method, _url, query, headers, cookies, data, ContentTypeApplicationJson, result, resultType)
}

func (client *Client) BuildAndDoRequestWithJsonForServerSentEvents(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any) (ServerSentEventMessages, error) {
	return client.BuildAndDoRequestWithDataForServerSentEvents(ctx, method, _url, query, headers, cookies, data, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithJsonForServerSentEventReader(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any) (*ServerSentEventReader, error) {
	return client.BuildAndDoRequestWithDataForServerSentEventReader(ctx, method, _url, query, headers, cookies, data, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithDataForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, dataType string, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithDataForResult(ctx, method, _url, query, headers, cookies, data, dataType, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithDataForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, dataType string, result any, resultType string) (*http.Response, error) {
	request, err := client.BuildRequestWithData(ctx, method, _url, query, headers, cookies, data, dataType)
	if err != nil {
		return nil, err
	}

	return client.DoRequestForResult(request, result, resultType)
}

func (client *Client) BuildAndDoRequestWithDataForServerSentEvents(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, dataType string) (ServerSentEventMessages, error) {
	reader, err := client.BuildAndDoRequestWithDataForServerSentEventReader(ctx, method, _url, query, headers, cookies, data, dataType)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return reader.ReadAll()
}

func (client *Client) BuildAndDoRequestWithDataForServerSentEventReader(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, dataType string) (*ServerSentEventReader, error) {
	response, err := client.BuildAndDoRequestWithData(ctx, method, _url, query, headers, cookies, data, dataType)
	if err != nil {
		return nil, err
	}

	return NewServerSentEventReaderFromHttpResponse(response, client.serverSentEventEndOfLine)
}

func (client *Client) BuildAndDoRequestWithFormForServerSentEventReader(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm) (*ServerSentEventReader, error) {
	response, err := client.BuildAndDoRequestWithForm(ctx, method, _url, query, headers, cookies, form)
	if err != nil {
		return nil, err
	}

	return NewServerSentEventReaderFromHttpResponse(response, client.serverSentEventEndOfLine)
}

func (client *Client) BuildAndDoRequestWithFormForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithFormForResult(ctx, method, _url, query, headers, cookies, form, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithStreamingFormForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithStreamingFormForResult(ctx, method, _url, query, headers, cookies, form, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithUrlEncodedFormForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form url.Values, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithUrlEncodedFormForResult(ctx, method, _url, query, headers, cookies, form, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithFormForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm, result any, resultType string) (*http.Response, error) {
	request, err := client.BuildRequestWithForm(ctx, method, _url, query, headers, cookies, form)
	if err != nil {
		return nil, err
	}

	return client.DoRequestForResult(request, result, resultType)
}

func (client *Client) BuildAndDoRequestWithStreamingFormForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm, result any, resultType string) (*http.Response, error) {
	request, err := client.BuildRequestWithStreamingForm(ctx, method, _url, query, headers, cookies, form)
	if err != nil {
		return nil, err
	}

	return client.DoRequestForResult(request, result, resultType)
}

func (client *Client) BuildAndDoRequestWithUrlEncodedFormForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form url.Values, result any, resultType string) (*http.Response, error) {
	request, err := client.BuildRequestWithUrlEncodedForm(ctx, method, _url, query, headers, cookies, form)
	if err != nil {
		return nil, err
	}
	return client.DoRequestForResult(request, result, resultType)
}

func (client *Client) BuildAndDoRequestWithBodyBytesForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body []byte, bodyType string, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithBodyBytesForResult(ctx, method, _url, query, headers, cookies, body, bodyType, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithBodyBytesForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body []byte, bodyType string, result any, resultType string) (*http.Response, error) {
	return client.BuildAndDoRequestWithBodyForResult(ctx, method, _url, query, headers, cookies, bytes.NewReader(body), bodyType, result, resultType)
}

func (client *Client) BuildAndDoRequestWithBodyForJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, bodyType string, result any) (*http.Response, error) {
	return client.BuildAndDoRequestWithBodyForResult(ctx, method, _url, query, headers, cookies, body, bodyType, result, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithBodyForResult(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, bodyType string, result any, resultType string) (*http.Response, error) {
	request, err := client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, body, bodyType)
	if err != nil {
		return nil, err
	}
	return client.DoRequestForResult(request, result, resultType)
}

// BuildAndDoRequestWithXxxx
// 构建请求、执行请求（但不处理响应）

func (client *Client) BuildAndDoRequestWithJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any) (*http.Response, error) {
	return client.BuildAndDoRequestWithData(ctx, method, _url, query, headers, cookies, data, ContentTypeApplicationJson)
}

func (client *Client) BuildAndDoRequestWithData(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, contentType string) (*http.Response, error) {
	request, err := client.BuildRequestWithData(ctx, method, _url, query, headers, cookies, data, contentType)
	if err != nil {
		return nil, err
	}
	return client.DoRequest(request)
}

func (client *Client) BuildAndDoRequestWithForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm) (*http.Response, error) {
	request, err := client.BuildRequestWithForm(ctx, method, _url, query, headers, cookies, form)
	if err != nil {
		return nil, err
	}
	return client.DoRequest(request)
}

func (client *Client) BuildAndDoRequestWithBodyBytes(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body []byte, contentType string) (*http.Response, error) {
	return client.BuildAndDoRequestWithBody(ctx, method, _url, query, headers, cookies, bytes.NewReader(body), contentType)
}

func (client *Client) BuildAndDoRequestWithBody(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, contentType string) (*http.Response, error) {
	request, err := client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, body, contentType)
	if err != nil {
		return nil, err
	}
	return client.DoRequest(request)
}

func (client *Client) BuildAndDoRangeRequestWithBody(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, contentType string, start, end int64) (dataStart, dataEnd, total int64, hasMore bool, response *http.Response, err error) {
	request, err := client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, body, contentType)
	if err != nil {
		return 0, 0, 0, false, nil, err
	}
	return client.DoRangeRequest(ctx, request, start, end)
}

func (client *Client) BuildAndDoRangeRequest(ctx context.Context, buildRequest func(ctx context.Context) (*http.Request, error), start, end int64) (dataStart, dataEnd int64, total int64, hasMore bool, response *http.Response, err error) {
	request, err := buildRequest(ctx)
	if err != nil {
		return 0, 0, 0, false, nil, err
	}

	return client.DoRangeRequest(ctx, request, start, end)
}

// BuildRequestWithXxxx
// 构建请求（但不执行请求）

func (client *Client) BuildRequestWithJson(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any) (*http.Request, error) {
	return client.BuildRequestWithData(ctx, method, _url, query, headers, cookies, data, ContentTypeApplicationJson)
}

func (client *Client) BuildRequestWithData(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, data any, contentType string) (*http.Request, error) {
	contentType = strings.ToLower(contentType)

	// todo refactor to call WithBodyBytes
	var body io.Reader
	if data != nil {
		var err error
		if body, err = client.BuildDataReaderAsType(data, contentType); err != nil {
			return nil, err
		}
	}

	return client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, body, contentType)
}

func (client *Client) BuildRequestWithForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm) (*http.Request, error) {
	bodyReader, contentType, err := form.MarshalToReader()
	if err != nil {
		return nil, err
	}

	return client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, bodyReader, contentType)
}

func (client *Client) BuildRequestWithStreamingForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form MultipartForm) (*http.Request, error) {
	bodyReader, contentType, err := form.MarshalToStreamingReadCloser()
	if err != nil {
		return nil, err
	}

	return client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, bodyReader, contentType)
}

func (client *Client) BuildRequestWithUrlEncodedForm(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, form url.Values) (*http.Request, error) {
	return client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, strings.NewReader(form.Encode()), ContentTypeApplicationXwwwFormUrlencoded)
}

func (client *Client) BuildRequestWithBodyBytes(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body []byte, contentType string) (*http.Request, error) {
	return client.BuildRequestWithBody(ctx, method, _url, query, headers, cookies, bytes.NewReader(body), contentType)
}

func (client *Client) BuildRequestWithBody(ctx context.Context, method, _url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, contentType string) (*http.Request, error) {
	loggerRef, ctx := client.LoggerFromContext(ctx)
	defer loggerRef.Reset()()

	loggerRef.With()
	method = strings.ToUpper(method)

	loggerRef.With(
		zap.String("Method", method),
		zap.String("Url", _url),
		zap.Any("Query", query),
		zap.Any("Headers", client.DesensitizeHeader(headers)),
	)

	loggerRef.Info("HttpxClient.RequestWithBodyBuilding")

	finalUrl, err := client.ResolveRawUrl(_url)
	if err != nil {
		loggerRef.Warn("HttpxClient.RequestWithBodyResolveRawUrlFailed",
			zap.Error(err),
		)

		return nil, err
	}

	if finalUrl == nil {
		loggerRef.Error("HttpxClient.RequestWithBodyResolveRawUrlReturnedNil")
		return nil, fmt.Errorf("ResolveRawUrlFailed: %w", err)
	}

	if len(client.query) > 0 || len(query) > 0 {
		finalUrl.RawQuery = stl.ConcatMaps(client.query, finalUrl.Query(), query).Encode()
	}

	headers = stl.ConcatMaps(client.BuiltinHeaders(), client.headers, headers)

	if body != nil {
		if contentType != "" {
			headers.Set("Content-Type", contentType)
		}
	}

	rawFinalUrl := finalUrl.String()
	_finalUrl := client.urlPrefixesReplacer.Replace(rawFinalUrl)
	if _finalUrl != rawFinalUrl {
		loggerRef.Info("HttpxClient.UrlPrefixReplaced",
			zap.String("RawUrl", rawFinalUrl),
			zap.String("FinalUrl", _finalUrl),
		)
	}

	var request *http.Request
	if ctx == nil {
		request, err = http.NewRequest(method, _finalUrl, body)
	} else {
		request, err = http.NewRequestWithContext(ctx, method, _finalUrl, body)
	}

	if err != nil {
		loggerRef.Warn("HttpxClient.RequestWithBodyNewRequestFailed",
			zap.Error(err),
		)
		return nil, err
	}

	host := headers.Get("Host")
	if host != "" {
		request.Host = host
	}

	request.Header = headers
	stl.ForEach(client.cookies, request.AddCookie)
	stl.ForEach(cookies, request.AddCookie)

	loggerRef.With(
		zap.String("Host", request.URL.Host),
		zap.String("Path", request.URL.Path),
		zap.Any("Headers", client.DesensitizeHeader(request.Header)),
	)

	loggerRef.Info("HttpxClient.RequestNewed")

	// 请求认证器
	if authenticator := client.requestAuthenticator; authenticator != nil {
		loggerRef.Info("HttpxClient.RequestWithBodyAuthenticating")
		if err = authenticator.Authenticate(request); err != nil {
			loggerRef.Warn("HttpxClient.RequestWithBodyAuthenticateFailed",
				zap.Error(err),
			)
			return nil, err
		}
		loggerRef.Info("HttpxClient.RequestWithBodyAuthenticated")
	}

	loggerRef.Info("HttpxClient.RequestWithBodyBuilt")

	return request, nil
}

func (client *Client) RenderAndBuildRequestWithTemplate(ctx context.Context, tpl *RequestTemplate, data any) (*http.Request, error) {
	request, err := tpl.Render(data)
	if err != nil {
		return nil, err
	}

	return client.BuildRequestWithBody(ctx, request.Method, request.Url, nil, request.Header.Native(), nil, bytes.NewReader(request.BodyBytes), "")
}

func (client *Client) ParseAndBuildRequestWithTemplate(ctx context.Context, tpl *RequestTemplate, funcMap templatex.TemplateFuncMap, data any) (*http.Request, error) {
	request, err := tpl.ParseAndRender(funcMap, data)
	if err != nil {
		return nil, err
	}

	return client.BuildRequestWithBody(ctx, request.Method, request.Url, nil, request.Header.Native(), nil, bytes.NewReader(request.BodyBytes), "")
}

// DoRequestForXxxx
// 执行请求、并处理响应

func (client *Client) DoRequestForResult(request *http.Request, result any, resultType string) (response *http.Response, err error) {
	response, err = client.DoRequest(request)
	if err != nil {
		return
	}

	err = ReadResponseBodyAsType(response, result, resultType)
	return
}

// DoRequest
// 执行请求（但不处理响应）

func (client *Client) DoRequest(request *http.Request) (response *http.Response, err error) {
	response, err = client.DoRequestForResponse(request)
	if err != nil {
		return
	}

	if responseDesensitizer := client.responseDesensitizer; responseDesensitizer != nil {
		response, err = responseDesensitizer.DesensitizeResponse(response)
		if err != nil {
			return
		}
	}

	if validator := client.responseValidator; validator != nil {
		err = validator(response)
	}

	return
}

func (client *Client) DoRangeRequest(ctx context.Context, request *http.Request, start, end int64) (dataStart, dataEnd int64, total int64, hasMore bool, response *http.Response, err error) {
	if start < 0 {
		return 0, 0, 0, false, nil, fmt.Errorf("invalid start: %d", start)
	}

	if end <= start {
		return 0, 0, 0, false, nil, fmt.Errorf("end is less than or equal to start: %d <= %d", end, start)
	}

	end -= 1

	request.Header.Set("Range", FormatRange("bytes", start, end))

	response, err = client.DoRequest(request)
	if err != nil {
		return 0, 0, 0, false, nil, err
	}

	unit, dataStart, dataEnd, total, err := Header(response.Header).ParseContentRange()
	if err != nil {
		return 0, 0, 0, false, nil, err
	}

	if unit != "bytes" {
		return 0, 0, 0, false, nil, fmt.Errorf("invalid content range unit: %s", unit)
	}

	if dataEnd >= 0 {
		dataEnd += 1
	}

	return dataStart, dataEnd, total, dataEnd >= 0 && dataEnd < total, response, nil
}

func (client *Client) DoRequestForResponse(request *http.Request) (response *http.Response, err error) {
	loggerRef, _ := client.LoggerFromContext(request.Context())
	defer loggerRef.Reset()()

	startTime := time.Now()
	loggerRef.With(
		zap.String("Method", request.Method),
		zap.String("Host", request.URL.Host),
		zap.String("Path", request.URL.Path),
	)

	loggerRef.Info("DoRequestForResponseStarting")

	response, err = client.Client.Do(request)
	if err != nil {
		loggerRef.Info("DoRequestForResponseFailed",
			zap.Error(err),
			zap.Duration("ExpiredDuration", time.Since(startTime)),
		)
		return
	}

	loggerRef.Info("DoRequestForResponseFinished",
		zap.Int("StatusCode", response.StatusCode),
		zap.Duration("ExpiredDuration", time.Since(startTime)),
	)

	return
}

func (client *Client) ResolveRawUrl(_url string) (*url.URL, error) {
	if _url == "" {
		return client.ResolveUrl(nil), nil
	}

	parsedUrl, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}

	return client.ResolveUrl(parsedUrl), nil
}

func (client *Client) ResolveUrl(_url *url.URL) *url.URL {
	baseUrl := client.GetBaseUrl()
	if baseUrl == nil {
		return _url
	}

	if _url == nil {
		return stl.Dup(baseUrl)
	}

	return netutils.JoinUrl(baseUrl, _url)
}

func (client *Client) BuildDataReaderAsType(data any, contentType string) (io.Reader, error) {
	return BuildDataReaderAsType(data, contentType)
}

func (client *Client) BuildDataReaderAsJson(data any) (io.Reader, error) {
	return BuildDataReaderAsJson(data)
}

func (client *Client) ReadResponseBody(respose *http.Response, data any) error {
	return ReadResponseBody(respose, data)
}

func (client *Client) ReadResponseBodyAsType(respose *http.Response, data any, contentType string) error {
	return ReadResponseBodyAsType(respose, data, contentType)
}

func (client *Client) ReadResponseBodyAsJson(respose *http.Response, data any) error {
	return ReadJson(respose.Body, data)
}

func (client *Client) TestConnectTcp() error {
	baseUrl := client.GetBaseUrl()
	if baseUrl == nil {
		return baseutils.NewNilError("httpx.Client.baseUrl")
	}

	defaultPort := 80
	if baseUrl.Scheme == "https" {
		defaultPort = 443
	}

	host, port, err := baseutils.ParseNetloc(baseUrl.Host, "", defaultPort)
	if err != nil {
		return err
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), time.Second*10)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func BuildDataReaderAsType(data any, contentType string) (io.Reader, error) {
	contentType = strings.ToLower(contentType)

	switch contentType {
	case ContentTypeApplicationJson:
		return BuildDataReaderAsJson(data)
	default:
		return nil, baseutils.NewNotImplementedError(contentType)
	}
}

func BuildDataReaderAsJson(data any) (io.Reader, error) {
	var bodyBuffer bytes.Buffer
	if err := json.NewEncoder(&bodyBuffer).Encode(data); err != nil {
		return nil, err
	}
	return &bodyBuffer, nil
}

func DetectResponseContentType(response *http.Response) string {
	contentType := strings.ToLower(response.Header.Get(HttpHeaderContentType))
	switch {
	case strings.Contains(contentType, ContentTypeApplicationJson):
		return ContentTypeApplicationJson
	default:
		return contentType
	}
}

func ReadResponseBody(response *http.Response, data any) error {
	return ReadResponseBodyAsType(response, data, "")
}

func ReadResponseBodyAsType(response *http.Response, data any, contentType string) (err error) {
	defer response.Body.Close()

	if data == nil {
		return nil
	}

	if contentType == "" {
		contentType = DetectResponseContentType(response)
	}

	// todo: performance?
	body, bodyData := TeeReaderToBuffer(response.Body)

	contentType = strings.ToLower(contentType)
	switch contentType {
	case ContentTypeApplicationJson:
		err = ReadJson(body, data)
	default:
		return baseutils.NewNotImplementedError(contentType)
	}

	if err != nil {
		return fmt.Errorf("ReadResponseBodyAsTypeFailed: statusCode=%d || expectedContentType=%s || contentType=%s || err=%w || body=%s", response.StatusCode, contentType, response.Header.Get("content-type"), err, bodyData.String())
	}

	return
}

func ReadResponseBodyAsJson(response *http.Response, data any) error {
	return ReadJson(response.Body, data)
}

func ReadJson(body io.Reader, data any) (err error) {
	buffer := types.NewBytesTruncatedBuffer(10 * 1024)
	if err := json.NewDecoder(io.TeeReader(body, buffer)).Decode(data); err != nil {
		return fmt.Errorf("ReadJsonFailed: err=%w || body=%s", err, buffer.Datas())
	}
	return nil
}

func TeeReaderToBuffer(reader io.Reader) (io.Reader, *bytes.Buffer) {
	buffer := bytes.NewBuffer(nil)
	return io.TeeReader(reader, buffer), buffer
}

func JsonRequestWithBody[Data any](ctx context.Context, client *Client, method, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body io.Reader, bodyType string) (data Data, response *http.Response, err error) {
	response, err = client.JsonRequestWithBody(ctx, method, url, query, headers, cookies, body, bodyType, &data)
	return
}

func JsonRequestWithForm[Data any](ctx context.Context, client *Client, method, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body MultipartForm) (data Data, response *http.Response, err error) {
	response, err = client.JsonRequestWithForm(ctx, method, url, query, headers, cookies, body, &data)
	return
}

func JsonRequestWithStreamingForm[Data any](ctx context.Context, client *Client, method, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body MultipartForm) (data Data, response *http.Response, err error) {
	response, err = client.JsonRequestWithStreamingForm(ctx, method, url, query, headers, cookies, body, &data)
	return
}

func JsonRequest[Data any](ctx context.Context, client *Client, method, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body any) (data Data, response *http.Response, err error) {
	response, err = client.JsonRequest(ctx, method, url, query, headers, cookies, body, &data)
	return
}

func JsonGet[Data any](ctx context.Context, client *Client, url string, query url.Values, headers http.Header, cookies []*http.Cookie) (data Data, response *http.Response, err error) {
	return JsonRequest[Data](ctx, client, http.MethodGet, url, query, headers, cookies, nil)
}

func JsonPost[Data any](ctx context.Context, client *Client, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body any) (data Data, response *http.Response, err error) {
	return JsonRequest[Data](ctx, client, http.MethodPost, url, query, headers, cookies, body)
}

func ValidateResponseForExpectedStatusCodes(response *http.Response, expecteds ...int) error {
	if response == nil {
		return baseutils.NewNilError("http.Response")
	}

	if statusCode := response.StatusCode; !stl.Contain(expecteds, statusCode) {
		defer response.Body.Close()
		bodyBytes, err := io.ReadAll(response.Body)
		return fmt.Errorf("httpx.ValidateResponseFailed: expecteds=%s || statusCode=%d || body=%s || err=%w || method=%s || url=%s",
			types.Strings(stl.Map(expecteds, strconv.Itoa)).Join(","),
			statusCode,
			string(bodyBytes),
			err,
			response.Request.Method,
			response.Request.URL,
		)
	}

	return nil
}

func NewResponseValidatorForExpectedStatusCodes(expecteds ...int) ResponseValidator {
	return func(response *http.Response) error {
		return ValidateResponseForExpectedStatusCodes(response, expecteds...)
	}
}

func ValidateResponseFor200(response *http.Response) error {
	return ValidateResponseForExpectedStatusCodes(response, 200)
}

func ValidateResponseFor2xx(response *http.Response) error {
	if response == nil {
		return baseutils.NewNilError("http.Response")
	}

	return ValidateResponseStatusCodeFor2xx(response.StatusCode)
}

func ValidateResponseStatusCodeFor2xx(statusCode int) error {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("httpx.ValidateResponseStatusCodeFor2xxFailed: statusCode=%d", statusCode)
	}

	return nil
}
