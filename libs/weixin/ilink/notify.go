package ilink

import (
	"context"
)

type NotifyRequest struct {
	BaseInfo *BaseInfo `json:"base_info,omitempty"`
}

type NotifyResponseBody struct {
	CommonResponseBody `json:",inline"`
}

func (body *NotifyResponseBody) GetData() *CommonResponseBody {
	if body == nil {
		return nil
	}
	return &body.CommonResponseBody
}

func (client *BotClient) NotifyStart(ctx context.Context) error {
	_, err := RequestForData[*NotifyResponseBody](client, ctx, "POST", NotifyStartPath, nil, &NotifyRequest{
		BaseInfo: client.BaseInfo(),
	})
	return err
}

func (client *BotClient) NotifyStop(ctx context.Context) error {
	_, err := RequestForData[*NotifyResponseBody](client, ctx, "POST", NotifyStopPath, nil, &NotifyRequest{
		BaseInfo: client.BaseInfo(),
	})
	return err
}
