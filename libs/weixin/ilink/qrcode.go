package ilink

type GetQrCodeResult struct {
	Qrcode           string `json:"qrcode"`
	QrcodeImgContent string `json:"qrcode_img_content"`
}

type GetQrcodeResponseBody struct {
	CommonResponseBody `json:",inline"`
	GetQrCodeResult    `json:",inline"`
}

func (body *GetQrcodeResponseBody) GetData() *GetQrCodeResult {
	if body == nil {
		return nil
	}
	return &body.GetQrCodeResult
}

type GetQrcodeStatusResult struct {
	Status      string `json:"status"`
	BaseUrl     string `json:"baseurl"`
	BotToken    string `json:"bot_token"`
	IlinkBotId  string `json:"ilink_bot_id"`
	IlinkUserId string `json:"ilink_user_id"`
}

type GetQrcodeStatusResponseBody struct {
	CommonResponseBody    `json:",inline"`
	GetQrcodeStatusResult `json:",inline"`
}

func (body *GetQrcodeStatusResponseBody) GetData() *GetQrcodeStatusResult {
	if body == nil {
		return nil
	}
	return &body.GetQrcodeStatusResult
}
