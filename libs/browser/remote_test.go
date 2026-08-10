package browser

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fasionchan/goutils/types"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type stubBrowser struct {
	mouseEvents []*MouseEvent
	keyEvents   []*KeyEvent
	meta        *ScreencastSessionMeta
	mouseErr    error
	keyErr      error
}

func (b *stubBrowser) GetCdpAddress() (*net.TCPAddr, error)        { return nil, nil }
func (b *stubBrowser) NewTab(options *NewTabOptions) (*Tab, error) { return nil, nil }
func (b *stubBrowser) GetTab(id string) (*Tab, error)              { return nil, nil }
func (b *stubBrowser) ListTabs() (Tabs, error)                     { return nil, nil }
func (b *stubBrowser) SwitchToTab(id string) error                 { return nil }
func (b *stubBrowser) CloseTab(id string) error                    { return nil }
func (b *stubBrowser) Navigate(id, url string) error               { return nil }
func (b *stubBrowser) GoBack(id string) error                      { return nil }
func (b *stubBrowser) GoForward(id string) error                   { return nil }
func (b *stubBrowser) Reload(id string) error                      { return nil }
func (b *stubBrowser) Click(id, selector, selectorType, button string, count int) error {
	return nil
}
func (b *stubBrowser) Type(id, selector, selectorType, text string) error { return nil }
func (b *stubBrowser) Hover(id, selector, selectorType string) error      { return nil }
func (b *stubBrowser) SelectOption(id, target, targetType string, options []string, optionType string, selected bool) error {
	return nil
}
func (b *stubBrowser) SetInputFiles(id, selector, selectorType string, files []string) error {
	return nil
}
func (b *stubBrowser) Screenshot(id string, opts *ScreenshotOptions) ([]byte, error) {
	return nil, nil
}
func (b *stubBrowser) Snapshot(id, snapshotType string) (string, error) { return "", nil }
func (b *stubBrowser) GetTexts(id, target, targetType string) (types.Strings, error) {
	return nil, nil
}
func (b *stubBrowser) GetHtmls(id, target, targetType string) (types.Strings, error) {
	return nil, nil
}
func (b *stubBrowser) SetCookies(id string, cookies []*http.Cookie) error { return nil }
func (b *stubBrowser) GetCookies(id string) ([]*http.Cookie, error)       { return nil, nil }
func (b *stubBrowser) PrintToPdf(id string) (io.ReadCloser, error)        { return nil, nil }
func (b *stubBrowser) StartScreencast(id string, opts *ScreencastOptions) (*ScreencastStream, error) {
	return nil, nil
}
func (b *stubBrowser) Close() error { return nil }

func (b *stubBrowser) DispatchMouseEvent(id string, event *MouseEvent) error {
	if b.mouseErr != nil {
		return b.mouseErr
	}
	copied := *event
	b.mouseEvents = append(b.mouseEvents, &copied)
	return nil
}

func (b *stubBrowser) DispatchKeyEvent(id string, event *KeyEvent) error {
	if b.keyErr != nil {
		return b.keyErr
	}
	copied := *event
	b.keyEvents = append(b.keyEvents, &copied)
	return nil
}

func (b *stubBrowser) GetScreencastSessionMeta(id string, opts *ScreencastOptions) (*ScreencastSessionMeta, error) {
	if b.meta != nil {
		return b.meta, nil
	}
	return &ScreencastSessionMeta{
		Format:         "jpeg",
		ViewportWidth:  1280,
		ViewportHeight: 720,
		FrameWidth:     1280,
		FrameHeight:    720,
	}, nil
}

func TestParseRemoteEnvelopeAndIgnoreUnknownPayloadFields(t *testing.T) {
	raw := []byte(`{"v":1,"id":"1","type":"mouse.down","payload":{"x":10,"y":20,"button":"left","extra":"ok"}}`)
	envelope, err := ParseRemoteEnvelope(raw)
	require.NoError(t, err)
	require.Equal(t, RemoteTypeMouseDown, envelope.Type)

	payload, err := DecodeRemotePayload[RemoteMousePayload](envelope.Payload)
	require.NoError(t, err)
	require.NotNil(t, payload.X)
	require.NotNil(t, payload.Y)
	require.Equal(t, 10.0, *payload.X)
	require.Equal(t, 20.0, *payload.Y)
	require.Equal(t, MouseButtonLeft, payload.Button)
}

func TestProcessRemoteControlMessageMouseDownAck(t *testing.T) {
	stub := &stubBrowser{}
	responses := ProcessRemoteControlMessage(stub, "tab-1", []byte(`{
		"v":1,"id":"req-1","type":"mouse.down","payload":{"x":640,"y":360,"button":"left","click_count":1}
	}`))
	require.Len(t, responses, 1)
	require.Equal(t, RemoteTypeAck, responses[0].Type)
	require.Equal(t, "req-1", responses[0].ID)
	require.Len(t, stub.mouseEvents, 1)
	require.Equal(t, MouseEventTypeDown, stub.mouseEvents[0].Type)
	require.Equal(t, 640.0, stub.mouseEvents[0].X)
	require.Equal(t, 360.0, stub.mouseEvents[0].Y)
}

func TestProcessRemoteControlMessageMouseMoveNoAckWithoutID(t *testing.T) {
	stub := &stubBrowser{}
	responses := ProcessRemoteControlMessage(stub, "tab-1", []byte(`{
		"v":1,"type":"mouse.move","payload":{"x":1,"y":2}
	}`))
	require.Empty(t, responses)
	require.Len(t, stub.mouseEvents, 1)
}

func TestProcessRemoteControlMessageMissingCoordinates(t *testing.T) {
	stub := &stubBrowser{}
	responses := ProcessRemoteControlMessage(stub, "tab-1", []byte(`{
		"v":1,"id":"req-2","type":"mouse.down","payload":{"button":"left"}
	}`))
	require.Len(t, responses, 1)
	require.Equal(t, RemoteTypeError, responses[0].Type)
	require.Equal(t, "req-2", responses[0].ID)

	var payload RemoteErrorPayload
	require.NoError(t, json.Unmarshal(responses[0].Payload, &payload))
	require.Equal(t, RemoteErrorCodeInvalidPayload, payload.Code)
	require.Empty(t, stub.mouseEvents)
}

func TestProcessRemoteControlMessageUnsupportedNav(t *testing.T) {
	stub := &stubBrowser{}
	responses := ProcessRemoteControlMessage(stub, "tab-1", []byte(`{
		"v":1,"id":"nav-1","type":"nav.go","payload":{"url":"https://example.com"}
	}`))
	require.Len(t, responses, 1)
	require.Equal(t, RemoteTypeError, responses[0].Type)

	var payload RemoteErrorPayload
	require.NoError(t, json.Unmarshal(responses[0].Payload, &payload))
	require.Equal(t, RemoteErrorCodeUnsupported, payload.Code)
	require.Equal(t, RemoteTypeNavGo, payload.RefType)
}

func TestProcessRemoteControlMessageKeyPress(t *testing.T) {
	stub := &stubBrowser{}
	responses := ProcessRemoteControlMessage(stub, "tab-1", []byte(`{
		"v":1,"id":"k1","type":"key.press","payload":{"key":"a","code":"KeyA","text":"a"}
	}`))
	require.Len(t, responses, 1)
	require.Equal(t, RemoteTypeAck, responses[0].Type)
	require.Len(t, stub.keyEvents, 2)
	require.Equal(t, KeyEventTypeDown, stub.keyEvents[0].Type)
	require.Equal(t, KeyEventTypeUp, stub.keyEvents[1].Type)
}

func TestProcessRemoteControlMessageSessionPing(t *testing.T) {
	stub := &stubBrowser{}
	responses := ProcessRemoteControlMessage(stub, "tab-1", []byte(`{"v":1,"id":"p1","type":"session.ping","payload":{}}`))
	require.Len(t, responses, 1)
	require.Equal(t, RemoteTypeSessionPong, responses[0].Type)
	require.Equal(t, "p1", responses[0].ID)
}

func TestEncodeInputModifiersAndEstimateFrameSize(t *testing.T) {
	require.Equal(t, 1|2|8, EncodeInputModifiers([]string{"alt", "ctrl", "shift"}))

	maxW, maxH := 640, 360
	fw, fh := EstimateScreencastFrameSize(1280, 720, &ScreencastOptions{MaxWidth: &maxW, MaxHeight: &maxH})
	require.Equal(t, 640, fw)
	require.Equal(t, 360, fh)
}

func TestScreencastOptionsToProtoMapsCDPFields(t *testing.T) {
	quality, maxW, maxH, nth := 80, 1280, 720, 2
	format := "png"
	opts := &ScreencastOptions{
		Format:        &format,
		Quality:       &quality,
		MaxWidth:      &maxW,
		MaxHeight:     &maxH,
		EventNthFrame: &nth,
	}

	req := screencastOptionsToProto(opts)
	require.Equal(t, "png", string(req.Format))
	require.Equal(t, &quality, req.Quality)
	require.Equal(t, &maxW, req.MaxWidth)
	require.Equal(t, &maxH, req.MaxHeight)
	require.Equal(t, &nth, req.EveryNthFrame)
}

func TestWriteScreencastFrameSendsMarkerThenBinary(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		controller := NewRemoteController(nil, "tab-1", nil)
		require.NoError(t, controller.writeScreencastFrame(conn, 0, "jpeg", []byte{0xff, 0xd8, 0xff}))
		require.NoError(t, controller.writeScreencastFrame(conn, 1, "jpeg", []byte{0xff, 0xd8, 0xfe}))
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msgType, data, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, msgType)

	var envelope RemoteEnvelope
	require.NoError(t, json.Unmarshal(data, &envelope))
	require.Equal(t, RemoteTypeScreencastFrame, envelope.Type)

	var payload RemoteScreencastFramePayload
	require.NoError(t, json.Unmarshal(envelope.Payload, &payload))
	require.Equal(t, uint64(0), payload.Seq)
	require.Equal(t, "jpeg", payload.Format)
	require.NotZero(t, payload.Ts)

	msgType, data, err = conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.BinaryMessage, msgType)
	require.Equal(t, []byte{0xff, 0xd8, 0xff}, data)

	msgType, data, err = conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, msgType)
	require.NoError(t, json.Unmarshal(data, &envelope))
	require.NoError(t, json.Unmarshal(envelope.Payload, &payload))
	require.Equal(t, uint64(1), payload.Seq)

	msgType, data, err = conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.BinaryMessage, msgType)
	require.Equal(t, []byte{0xff, 0xd8, 0xfe}, data)
}

func TestRodInputEventTypeConverters(t *testing.T) {
	mouseType, err := rodMouseEventTypeFromStd(MouseEventTypeWheel)
	require.NoError(t, err)
	require.Equal(t, "mouseWheel", string(mouseType))

	keyType, err := rodKeyEventTypeFromStd(KeyEventTypeDown)
	require.NoError(t, err)
	require.Equal(t, "keyDown", string(keyType))

	_, err = rodMouseEventTypeFromStd("unknown")
	require.Error(t, err)
	_, err = rodKeyEventTypeFromStd("unknown")
	require.Error(t, err)
}

var (
	_ Browser = (*stubBrowser)(nil)
	_ Browser = (*RodBrowser)(nil)
)
