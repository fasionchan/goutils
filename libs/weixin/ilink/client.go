package ilink

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fasionchan/goutils/libs/httpx"
	"github.com/fasionchan/goutils/stl"
)

type ClientConfig struct {
	baseUrl string
}

func (config *ClientConfig) NewClient() (*Client, error) {
	baseUrl := config.baseUrl
	if baseUrl == "" {
		baseUrl = DefaultBaseUrl
	}

	httpClient, err := httpx.NewClient(baseUrl)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: httpClient,
	}, nil
}

type NewClientOption = stl.Option[*ClientConfig]

type Client struct {
	*httpx.Client
}

func NewClient(opts ...NewClientOption) (*Client, error) {
	config := &ClientConfig{
		baseUrl: DefaultBaseUrl,
	}
	return stl.NewOptions(opts...).Apply(config).NewClient()
}

func (client *Client) GetBotQrcode(ctx context.Context, botType int) (*GetQrCodeResult, error) {
	query := url.Values{}
	query.Set("bot_type", strconv.Itoa(botType))
	return RequestForData[*GetQrcodeResponseBody](client, ctx, "POST", GetBotQrcodePath, query, nil)
}

func (client *Client) GetQrcodeStatus(ctx context.Context, qrcode string) (*GetQrcodeStatusResult, error) {
	query := url.Values{}
	query.Set("qrcode", qrcode)
	return RequestForData[*GetQrcodeStatusResponseBody](client, ctx, "POST", GetQrcodeStatusPath, query, nil)
}

func (client *Client) WaitQrcodeStatus(ctx context.Context, qrcode string) (*GetQrcodeStatusResult, error) {
	for {
		result, err := client.GetQrcodeStatus(ctx, qrcode)
		if err != nil {
			return nil, err
		}

		switch result.Status {
		case GetQrCodeStatusConfirmed:
			return result, nil
		case GetQrCodeStatusExpired:
			return nil, fmt.Errorf("qrcode expired")
		case GetQrCodeStatusWait:
			continue
		default:
			return nil, fmt.Errorf("unknown status: %s", result.Status)
		}
	}
}

func (client *Client) WaitQrcodeStatusForBotClient(ctx context.Context, qrcode string) (*BotClient, error) {
	result, err := client.WaitQrcodeStatus(ctx, qrcode)
	if err != nil {
		return nil, err
	}

	return NewBotClient(result.BaseUrl, result.BotToken)
}

type jsonRequester interface {
	JsonRequest(ctx context.Context, method, url string, query url.Values, headers http.Header, cookies []*http.Cookie, body, data any) (*http.Response, error)
}

func RequestForData[
	ResponseBody interface {
		GetData() Data
		GetError() error
	},
	Data any,
](client jsonRequester, ctx context.Context, method string, path string, query url.Values, body any) (data Data, err error) {
	var responseBody ResponseBody
	response, err := client.JsonRequest(ctx, method, path, query, nil, nil, body, &responseBody)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("request %s failed: %s", path, response.Status)
		return
	}

	if err = responseBody.GetError(); err != nil {
		return
	}

	return responseBody.GetData(), nil
}
