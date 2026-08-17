package ilink

import (
	"context"
)

type SendMessageRequest struct {
	Msg      *WeixinMessage `json:"msg,omitempty"`
	BaseInfo *BaseInfo      `json:"base_info,omitempty"`
}

type SendMessageResponseBody struct {
	CommonResponseBody `json:",inline"`
}

func (body *SendMessageResponseBody) GetData() *CommonResponseBody {
	if body == nil {
		return nil
	}
	return &body.CommonResponseBody
}

func (client *BotClient) SendMessage(ctx context.Context, msg *WeixinMessage) error {
	_, err := RequestForData[*SendMessageResponseBody](client, ctx, "POST", SendMessagePath, nil, &SendMessageRequest{
		Msg:      msg,
		BaseInfo: client.BaseInfo(),
	})
	return err
}
