package ilink

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomWechatUin(t *testing.T) {
	seen := map[string]struct{}{}
	for i := 0; i < 20; i++ {
		uin := RandomWechatUin()
		require.NotEmpty(t, uin)

		decoded, err := base64.StdEncoding.DecodeString(uin)
		require.NoError(t, err)

		value, err := strconv.ParseUint(string(decoded), 10, 32)
		require.NoError(t, err)
		assert.LessOrEqual(t, value, uint64(^uint32(0)))
		seen[uin] = struct{}{}
	}
	assert.Greater(t, len(seen), 1)
}

func TestCommonResponseBodyErrMsgAliases(t *testing.T) {
	var body CommonResponseBody
	require.NoError(t, json.Unmarshal([]byte(`{"ret":1,"err_msg":"legacy"}`), &body))
	assert.Equal(t, "legacy", body.GetErrMsg())
	assert.EqualError(t, body.GetError(), "ret: 1, err_msg: legacy")

	body = CommonResponseBody{}
	require.NoError(t, json.Unmarshal([]byte(`{"ret":2,"errmsg":"official","errcode":-14}`), &body))
	assert.Equal(t, "official", body.GetErrMsg())
	assert.EqualError(t, body.GetError(), "ret: 2, errcode: -14, err_msg: official")
}

func TestWeixinMessageJSON(t *testing.T) {
	raw := []byte(`{
		"from_user_id": "user@im.wechat",
		"to_user_id": "bot@im.bot",
		"message_type": 1,
		"message_state": 2,
		"context_token": "token",
		"item_list": [
			{"type": 1, "text_item": {"text": "你好"}},
			{"type": 2, "image_item": {"aeskey": "abc", "media": {"full_url": "https://cdn.example/1"}}},
			{"type": 3, "voice_item": {"encode_type": 6, "text": "转写"}},
			{"type": 4, "file_item": {"file_name": "a.txt", "len": "12"}},
			{"type": 5, "video_item": {"play_length": 3}}
		]
	}`)

	var msg WeixinMessage
	require.NoError(t, json.Unmarshal(raw, &msg))
	assert.Equal(t, "user@im.wechat", msg.FromUserId)
	assert.Equal(t, MessageTypeUser, msg.MessageType)
	require.Len(t, msg.ItemList, 5)
	assert.Equal(t, "你好", msg.ItemList[0].TextItem.Text)
	assert.Equal(t, "abc", msg.ItemList[1].ImageItem.Aeskey)
	assert.Equal(t, VoiceEncodeTypeSilk, msg.ItemList[2].VoiceItem.EncodeType)
	assert.Equal(t, "12", msg.ItemList[3].FileItem.Len)
	assert.Equal(t, 3, msg.ItemList[4].VideoItem.PlayLength)
}

type capturedRequest struct {
	method string
	path   string
	header http.Header
	body   []byte
}

func newBotClientTestServer(t *testing.T, response string, captured *capturedRequest) *BotClient {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		captured.method = r.Method
		captured.path = r.URL.Path
		captured.header = r.Header.Clone()
		captured.body = body
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	client, err := NewBotClient(server.URL, "test-token")
	require.NoError(t, err)
	return client
}

func assertBotAuthHeaders(t *testing.T, header http.Header) {
	t.Helper()
	assert.Equal(t, "Bearer test-token", header.Get("Authorization"))
	assert.Equal(t, "ilink_bot_token", header.Get("AuthorizationType"))

	uin := header.Get("X-WECHAT-UIN")
	require.NotEmpty(t, uin)
	decoded, err := base64.StdEncoding.DecodeString(uin)
	require.NoError(t, err)
	_, err = strconv.ParseUint(string(decoded), 10, 32)
	require.NoError(t, err)
}

func TestBotClientGetUpdates(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":0,"msgs":[{"from_user_id":"u1","context_token":"ct","item_list":[{"type":1,"text_item":{"text":"hi"}}]}],"get_updates_buf":"next","longpolling_timeout_ms":35000}`, &captured)

	result, err := client.GetUpdates(context.Background(), "prev")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "next", result.GetUpdatesBuf)
	assert.Equal(t, 35000, result.LongpollingTimeoutMs)
	require.Len(t, result.Msgs, 1)
	assert.Equal(t, "hi", result.Msgs[0].ItemList[0].TextItem.Text)

	assert.Equal(t, http.MethodPost, captured.method)
	assert.Equal(t, GetUpdatesPath, captured.path)
	assertBotAuthHeaders(t, captured.header)

	var req GetUpdatesRequest
	require.NoError(t, json.Unmarshal(captured.body, &req))
	assert.Equal(t, "prev", req.GetUpdatesBuf)
	require.NotNil(t, req.BaseInfo)
	assert.Equal(t, DefaultChannelVersion, req.BaseInfo.ChannelVersion)
	assert.Equal(t, DefaultBotAgent, req.BaseInfo.BotAgent)
}

func TestBotClientSendMessage(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":0}`, &captured)

	err := client.SendMessage(context.Background(), &WeixinMessage{
		ToUserId:     "user@im.wechat",
		MessageType:  MessageTypeBot,
		MessageState: MessageStateFinish,
		ContextToken: "ct",
		ItemList: []*MessageItem{
			{Type: MessageItemTypeText, TextItem: &TextItem{Text: "回复"}},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, SendMessagePath, captured.path)
	assertBotAuthHeaders(t, captured.header)

	var req SendMessageRequest
	require.NoError(t, json.Unmarshal(captured.body, &req))
	require.NotNil(t, req.Msg)
	assert.Equal(t, "user@im.wechat", req.Msg.ToUserId)
	assert.Equal(t, "ct", req.Msg.ContextToken)
	assert.Equal(t, "回复", req.Msg.ItemList[0].TextItem.Text)
	require.NotNil(t, req.BaseInfo)
}

func TestBotClientGetUploadUrl(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":0,"upload_param":"p1","upload_full_url":"https://cdn.example/put"}`, &captured)

	result, err := client.GetUploadUrl(context.Background(), &GetUploadUrlRequest{
		Filekey:     "filekey",
		MediaType:   UploadMediaTypeImage,
		ToUserId:    "user@im.wechat",
		Rawsize:     12,
		Rawfilemd5:  "md5",
		Filesize:    16,
		NoNeedThumb: true,
		Aeskey:      "aes",
	})
	require.NoError(t, err)
	assert.Equal(t, "p1", result.UploadParam)
	assert.Equal(t, "https://cdn.example/put", result.UploadFullUrl)
	assert.Equal(t, GetUploadUrlPath, captured.path)

	var req GetUploadUrlRequest
	require.NoError(t, json.Unmarshal(captured.body, &req))
	assert.Equal(t, "filekey", req.Filekey)
	assert.Equal(t, UploadMediaTypeImage, req.MediaType)
	assert.True(t, req.NoNeedThumb)
	require.NotNil(t, req.BaseInfo)
}

func TestBotClientGetConfig(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":0,"typing_ticket":"ticket"}`, &captured)

	result, err := client.GetConfig(context.Background(), "user@im.wechat", "ct")
	require.NoError(t, err)
	assert.Equal(t, "ticket", result.TypingTicket)
	assert.Equal(t, GetConfigPath, captured.path)

	var req GetConfigRequest
	require.NoError(t, json.Unmarshal(captured.body, &req))
	assert.Equal(t, "user@im.wechat", req.IlinkUserId)
	assert.Equal(t, "ct", req.ContextToken)
}

func TestBotClientSendTyping(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":0}`, &captured)

	err := client.SendTyping(context.Background(), &SendTypingRequest{
		IlinkUserId:  "user@im.wechat",
		TypingTicket: "ticket",
		Status:       TypingStatusTyping,
	})
	require.NoError(t, err)
	assert.Equal(t, SendTypingPath, captured.path)

	var req SendTypingRequest
	require.NoError(t, json.Unmarshal(captured.body, &req))
	assert.Equal(t, "user@im.wechat", req.IlinkUserId)
	assert.Equal(t, TypingStatusTyping, req.Status)
}

func TestBotClientNotify(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":0}`, &captured)

	require.NoError(t, client.NotifyStart(context.Background()))
	assert.Equal(t, NotifyStartPath, captured.path)
	assertBotAuthHeaders(t, captured.header)

	require.NoError(t, client.NotifyStop(context.Background()))
	assert.Equal(t, NotifyStopPath, captured.path)
}

func TestBotClientApiError(t *testing.T) {
	var captured capturedRequest
	client := newBotClientTestServer(t, `{"ret":1,"errmsg":"session timeout","errcode":-14}`, &captured)

	_, err := client.GetUpdates(context.Background(), "")
	require.Error(t, err)
	assert.EqualError(t, err, "ret: 1, errcode: -14, err_msg: session timeout")
}
