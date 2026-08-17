package ilink

type TextItem struct {
	Text string `json:"text,omitempty"`
}

type CDNMedia struct {
	EncryptQueryParam string `json:"encrypt_query_param,omitempty"`
	AesKey            string `json:"aes_key,omitempty"`
	EncryptType       int    `json:"encrypt_type,omitempty"`
	FullUrl           string `json:"full_url,omitempty"`
}

type ImageItem struct {
	Media       *CDNMedia `json:"media,omitempty"`
	ThumbMedia  *CDNMedia `json:"thumb_media,omitempty"`
	Aeskey      string    `json:"aeskey,omitempty"`
	Url         string    `json:"url,omitempty"`
	MidSize     int64     `json:"mid_size,omitempty"`
	ThumbSize   int64     `json:"thumb_size,omitempty"`
	ThumbHeight int       `json:"thumb_height,omitempty"`
	ThumbWidth  int       `json:"thumb_width,omitempty"`
	HdSize      int64     `json:"hd_size,omitempty"`
}

type VoiceItem struct {
	Media         *CDNMedia `json:"media,omitempty"`
	EncodeType    int       `json:"encode_type,omitempty"`
	BitsPerSample int       `json:"bits_per_sample,omitempty"`
	SampleRate    int       `json:"sample_rate,omitempty"`
	Playtime      int       `json:"playtime,omitempty"`
	Text          string    `json:"text,omitempty"`
}

type FileItem struct {
	Media    *CDNMedia `json:"media,omitempty"`
	FileName string    `json:"file_name,omitempty"`
	Md5      string    `json:"md5,omitempty"`
	Len      string    `json:"len,omitempty"`
}

type VideoItem struct {
	Media       *CDNMedia `json:"media,omitempty"`
	VideoSize   int64     `json:"video_size,omitempty"`
	PlayLength  int       `json:"play_length,omitempty"`
	VideoMd5    string    `json:"video_md5,omitempty"`
	ThumbMedia  *CDNMedia `json:"thumb_media,omitempty"`
	ThumbSize   int64     `json:"thumb_size,omitempty"`
	ThumbHeight int       `json:"thumb_height,omitempty"`
	ThumbWidth  int       `json:"thumb_width,omitempty"`
}

type RefMessage struct {
	MessageItem *MessageItem `json:"message_item,omitempty"`
	Title       string       `json:"title,omitempty"`
}

type ToolCallStartItem struct {
	ToolName   string `json:"tool_name,omitempty"`
	ToolCallId string `json:"tool_call_id,omitempty"`
}

type ToolCallResultItem struct {
	ToolName   string `json:"tool_name,omitempty"`
	ToolCallId string `json:"tool_call_id,omitempty"`
	Status     string `json:"status,omitempty"`
}

type MessageItem struct {
	Type               int                 `json:"type,omitempty"`
	CreateTimeMs       int64               `json:"create_time_ms,omitempty"`
	UpdateTimeMs       int64               `json:"update_time_ms,omitempty"`
	IsCompleted        bool                `json:"is_completed,omitempty"`
	MsgId              string              `json:"msg_id,omitempty"`
	RefMsg             *RefMessage         `json:"ref_msg,omitempty"`
	TextItem           *TextItem           `json:"text_item,omitempty"`
	ImageItem          *ImageItem          `json:"image_item,omitempty"`
	VoiceItem          *VoiceItem          `json:"voice_item,omitempty"`
	FileItem           *FileItem           `json:"file_item,omitempty"`
	VideoItem          *VideoItem          `json:"video_item,omitempty"`
	ToolCallStartItem  *ToolCallStartItem  `json:"tool_call_start_item,omitempty"`
	ToolCallResultItem *ToolCallResultItem `json:"tool_call_result_item,omitempty"`
}

type WeixinMessage struct {
	Seq          int64          `json:"seq,omitempty"`
	MessageId    int64          `json:"message_id,omitempty"`
	FromUserId   string         `json:"from_user_id,omitempty"`
	ToUserId     string         `json:"to_user_id,omitempty"`
	ClientId     string         `json:"client_id,omitempty"`
	CreateTimeMs int64          `json:"create_time_ms,omitempty"`
	UpdateTimeMs int64          `json:"update_time_ms,omitempty"`
	DeleteTimeMs int64          `json:"delete_time_ms,omitempty"`
	SessionId    string         `json:"session_id,omitempty"`
	GroupId      string         `json:"group_id,omitempty"`
	MessageType  int            `json:"message_type,omitempty"`
	MessageState int            `json:"message_state,omitempty"`
	ItemList     []*MessageItem `json:"item_list,omitempty"`
	ContextToken string         `json:"context_token,omitempty"`
	RunId        string         `json:"run_id,omitempty"`
}
