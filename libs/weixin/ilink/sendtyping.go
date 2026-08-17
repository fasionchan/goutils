package ilink

import (
	"context"
)

type SendTypingRequest struct {
	IlinkUserId  string    `json:"ilink_user_id,omitempty"`
	TypingTicket string    `json:"typing_ticket,omitempty"`
	Status       int       `json:"status,omitempty"`
	BaseInfo     *BaseInfo `json:"base_info,omitempty"`
}

type SendTypingResponseBody struct {
	CommonResponseBody `json:",inline"`
}

func (body *SendTypingResponseBody) GetData() *CommonResponseBody {
	if body == nil {
		return nil
	}
	return &body.CommonResponseBody
}

func (client *BotClient) SendTyping(ctx context.Context, req *SendTypingRequest) error {
	if req == nil {
		req = &SendTypingRequest{}
	} else {
		copied := *req
		req = &copied
	}
	req.BaseInfo = client.BaseInfo()
	_, err := RequestForData[*SendTypingResponseBody](client, ctx, "POST", SendTypingPath, nil, req)
	return err
}
