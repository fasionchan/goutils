package ilink

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"strconv"
	"time"

	"github.com/fasionchan/goutils/libs/httpx"
)

type BotClient struct {
	*httpx.Client
	baseInfo *BaseInfo
}

func NewBotClient(baseUrl, botToken string) (*BotClient, error) {
	if baseUrl == "" {
		baseUrl = DefaultBaseUrl
	}

	httpClient, err := httpx.NewClient(baseUrl)
	if err != nil {
		return nil, err
	}

	return &BotClient{
		Client: httpClient.
			WithBearerToken(botToken).
			SetHeader("AuthorizationType", "ilink_bot_token").
			WithRequestAuthenticatorFunc(func(r *http.Request) error {
				r.Header.Set("X-WECHAT-UIN", RandomWechatUin())
				return nil
			}),
		baseInfo: DefaultBaseInfo(),
	}, nil
}

func (client *BotClient) BaseInfo() *BaseInfo {
	if client == nil || client.baseInfo == nil {
		return DefaultBaseInfo()
	}
	return client.baseInfo
}

func (client *BotClient) WithBaseInfo(baseInfo *BaseInfo) *BotClient {
	if client == nil {
		return nil
	}
	client.baseInfo = baseInfo
	return client
}

func RandomWechatUin() string {
	var buf [4]byte
	if _, err := rand.Read(buf[:]); err != nil {
		binary.BigEndian.PutUint32(buf[:], uint32(time.Now().UnixNano()))
	}
	uin := binary.BigEndian.Uint32(buf[:])
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(uin), 10)))
}
