package ilink

import (
	"context"
)

type GetConfigRequest struct {
	IlinkUserId  string    `json:"ilink_user_id,omitempty"`
	ContextToken string    `json:"context_token,omitempty"`
	BaseInfo     *BaseInfo `json:"base_info,omitempty"`
}

type GetConfigResult struct {
	TypingTicket string `json:"typing_ticket,omitempty"`
}

type GetConfigResponseBody struct {
	CommonResponseBody `json:",inline"`
	GetConfigResult    `json:",inline"`
}

func (body *GetConfigResponseBody) GetData() *GetConfigResult {
	if body == nil {
		return nil
	}
	return &body.GetConfigResult
}

func (client *BotClient) GetConfig(ctx context.Context, ilinkUserId, contextToken string) (*GetConfigResult, error) {
	return RequestForData[*GetConfigResponseBody](client, ctx, "POST", GetConfigPath, nil, &GetConfigRequest{
		IlinkUserId:  ilinkUserId,
		ContextToken: contextToken,
		BaseInfo:     client.BaseInfo(),
	})
}
