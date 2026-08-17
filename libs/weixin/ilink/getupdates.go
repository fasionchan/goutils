package ilink

import (
	"context"
)

type GetUpdatesRequest struct {
	SyncBuf       string    `json:"sync_buf,omitempty"`
	GetUpdatesBuf string    `json:"get_updates_buf"`
	BaseInfo      *BaseInfo `json:"base_info,omitempty"`
}

type GetUpdatesResult struct {
	Msgs                 []*WeixinMessage `json:"msgs"`
	SyncBuf              string           `json:"sync_buf,omitempty"`
	GetUpdatesBuf        string           `json:"get_updates_buf"`
	LongpollingTimeoutMs int              `json:"longpolling_timeout_ms,omitempty"`
}

type GetUpdatesResponseBody struct {
	CommonResponseBody `json:",inline"`
	GetUpdatesResult   `json:",inline"`
}

func (body *GetUpdatesResponseBody) GetData() *GetUpdatesResult {
	if body == nil {
		return nil
	}
	return &body.GetUpdatesResult
}

func (client *BotClient) GetUpdates(ctx context.Context, getUpdatesBuf string) (*GetUpdatesResult, error) {
	return RequestForData[*GetUpdatesResponseBody](client, ctx, "POST", GetUpdatesPath, nil, &GetUpdatesRequest{
		GetUpdatesBuf: getUpdatesBuf,
		BaseInfo:      client.BaseInfo(),
	})
}
