package ilink

import (
	"context"
)

type GetUploadUrlRequest struct {
	Filekey         string    `json:"filekey,omitempty"`
	MediaType       int       `json:"media_type,omitempty"`
	ToUserId        string    `json:"to_user_id,omitempty"`
	Rawsize         int64     `json:"rawsize,omitempty"`
	Rawfilemd5      string    `json:"rawfilemd5,omitempty"`
	Filesize        int64     `json:"filesize,omitempty"`
	ThumbRawsize    int64     `json:"thumb_rawsize,omitempty"`
	ThumbRawfilemd5 string    `json:"thumb_rawfilemd5,omitempty"`
	ThumbFilesize   int64     `json:"thumb_filesize,omitempty"`
	NoNeedThumb     bool      `json:"no_need_thumb,omitempty"`
	Aeskey          string    `json:"aeskey,omitempty"`
	BaseInfo        *BaseInfo `json:"base_info,omitempty"`
}

type GetUploadUrlResult struct {
	UploadParam      string `json:"upload_param,omitempty"`
	ThumbUploadParam string `json:"thumb_upload_param,omitempty"`
	UploadFullUrl    string `json:"upload_full_url,omitempty"`
}

type GetUploadUrlResponseBody struct {
	CommonResponseBody `json:",inline"`
	GetUploadUrlResult `json:",inline"`
}

func (body *GetUploadUrlResponseBody) GetData() *GetUploadUrlResult {
	if body == nil {
		return nil
	}
	return &body.GetUploadUrlResult
}

func (client *BotClient) GetUploadUrl(ctx context.Context, req *GetUploadUrlRequest) (*GetUploadUrlResult, error) {
	if req == nil {
		req = &GetUploadUrlRequest{}
	} else {
		copied := *req
		req = &copied
	}
	req.BaseInfo = client.BaseInfo()
	return RequestForData[*GetUploadUrlResponseBody](client, ctx, "POST", GetUploadUrlPath, nil, req)
}
