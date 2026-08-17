package ilink

import (
	"fmt"
)

const (
	DefaultBaseUrl = "https://ilinkai.weixin.qq.com"

	DefaultChannelVersion = "unknown"
	DefaultBotAgent       = "OpenClaw"

	GetBotQrcodePath    = "/ilink/bot/get_bot_qrcode"
	GetQrcodeStatusPath = "/ilink/bot/get_qrcode_status"
	GetUpdatesPath      = "/ilink/bot/getupdates"
	SendMessagePath     = "/ilink/bot/sendmessage"
	GetUploadUrlPath    = "/ilink/bot/getuploadurl"
	GetConfigPath       = "/ilink/bot/getconfig"
	SendTypingPath      = "/ilink/bot/sendtyping"
	NotifyStartPath     = "/ilink/bot/msg/notifystart"
	NotifyStopPath      = "/ilink/bot/msg/notifystop"

	GetQrCodeStatusConfirmed = "confirmed"
	GetQrCodeStatusWait      = "wait"
	GetQrCodeStatusExpired   = "expired"

	MessageTypeNone = 0
	MessageTypeUser = 1
	MessageTypeBot  = 2

	MessageItemTypeNone           = 0
	MessageItemTypeText           = 1
	MessageItemTypeImage          = 2
	MessageItemTypeVoice          = 3
	MessageItemTypeFile           = 4
	MessageItemTypeVideo          = 5
	MessageItemTypeToolCallStart  = 11
	MessageItemTypeToolCallResult = 12

	MessageStateNew        = 0
	MessageStateGenerating = 1
	MessageStateFinish     = 2

	UploadMediaTypeImage = 1
	UploadMediaTypeVideo = 2
	UploadMediaTypeFile  = 3
	UploadMediaTypeVoice = 4

	TypingStatusTyping = 1
	TypingStatusCancel = 2

	VoiceEncodeTypePcm      = 1
	VoiceEncodeTypeAdpcm    = 2
	VoiceEncodeTypeFeature  = 3
	VoiceEncodeTypeSpeex    = 4
	VoiceEncodeTypeAmr      = 5
	VoiceEncodeTypeSilk     = 6
	VoiceEncodeTypeMp3      = 7
	VoiceEncodeTypeOggSpeex = 8
)

type BaseInfo struct {
	ChannelVersion string `json:"channel_version"`
	BotAgent       string `json:"bot_agent"`
}

func DefaultBaseInfo() *BaseInfo {
	return &BaseInfo{
		ChannelVersion: DefaultChannelVersion,
		BotAgent:       DefaultBotAgent,
	}
}

type CommonResponseBody struct {
	Ret     int    `json:"ret"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	ErrMsg2 string `json:"err_msg"`
}

func (body *CommonResponseBody) Error() string {
	if body == nil {
		return "CommonResponseBody is nil"
	}
	msg := body.GetErrMsg()
	if body.ErrCode != 0 {
		return fmt.Sprintf("ret: %d, errcode: %d, err_msg: %s", body.Ret, body.ErrCode, msg)
	}
	return fmt.Sprintf("ret: %d, err_msg: %s", body.Ret, msg)
}

func (body *CommonResponseBody) GetErrMsg() string {
	if body == nil {
		return ""
	}
	if body.ErrMsg != "" {
		return body.ErrMsg
	}
	return body.ErrMsg2
}

func (body *CommonResponseBody) GetError() error {
	if body.IsSuccess() {
		return nil
	}
	return body
}

func (body *CommonResponseBody) IsSuccess() bool {
	return body.Ret == 0
}
