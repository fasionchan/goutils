package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	RemoteProtocolVersion = 1

	RemoteTypeSessionPing    = "session.ping"
	RemoteTypeSessionPong    = "session.pong"
	RemoteTypeSessionReady   = "session.ready"
	RemoteTypeScreencastMeta = "screencast.meta"

	RemoteTypeMouseMove  = "mouse.move"
	RemoteTypeMouseDown  = "mouse.down"
	RemoteTypeMouseUp    = "mouse.up"
	RemoteTypeMouseWheel = "mouse.wheel"

	RemoteTypeKeyDown  = "key.down"
	RemoteTypeKeyUp    = "key.up"
	RemoteTypeKeyPress = "key.press"

	RemoteTypeAck   = "ack"
	RemoteTypeError = "error"

	// Reserved namespaces (not implemented in this change).
	RemoteTypeNavGo          = "nav.go"
	RemoteTypeNavBack        = "nav.back"
	RemoteTypeNavForward     = "nav.forward"
	RemoteTypeNavReload      = "nav.reload"
	RemoteTypeEventNavigated = "event.navigated"
	RemoteTypeEventConsole   = "event.console"
	RemoteTypeEventDialog    = "event.dialog"
	RemoteTypeEventTabClosed = "event.tab_closed"

	RemoteErrorCodeInvalidJSON    = "invalid_json"
	RemoteErrorCodeInvalidPayload = "invalid_payload"
	RemoteErrorCodeUnsupported    = "unsupported"
	RemoteErrorCodeInternal       = "internal"
)

type ScreencastOptionsQuery struct {
	Format        *string `query:"format,omitempty"`
	Quality       *int    `query:"quality,omitempty"`
	MaxWidth      *int    `query:"max_width,omitempty"`
	MaxHeight     *int    `query:"max_height,omitempty"`
	EventNthFrame *int    `query:"event_nth_frame,omitempty"`
}

func (query *ScreencastOptionsQuery) ToScreencastOptions() *ScreencastOptions {
	if query == nil {
		return nil
	}

	return &ScreencastOptions{
		Format:        query.Format,
		Quality:       query.Quality,
		MaxWidth:      query.MaxWidth,
		MaxHeight:     query.MaxHeight,
		EventNthFrame: query.EventNthFrame,
	}
}

// RemoteEnvelope is the versioned JSON control-plane message.
type RemoteEnvelope struct {
	V       int             `json:"v"`
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Ts      int64           `json:"ts,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type RemoteMousePayload struct {
	X          *float64 `json:"x"`
	Y          *float64 `json:"y"`
	Button     string   `json:"button,omitempty"`
	ClickCount int      `json:"click_count,omitempty"`
	DeltaX     float64  `json:"delta_x,omitempty"`
	DeltaY     float64  `json:"delta_y,omitempty"`
	Modifiers  []string `json:"modifiers,omitempty"`
}

type RemoteKeyPayload struct {
	Key        string   `json:"key,omitempty"`
	Code       string   `json:"code,omitempty"`
	Text       string   `json:"text,omitempty"`
	Modifiers  []string `json:"modifiers,omitempty"`
	AutoRepeat bool     `json:"auto_repeat,omitempty"`
}

type RemoteSessionReadyPayload struct {
	TabID           string `json:"tab_id"`
	ProtocolVersion int    `json:"protocol_version"`
}

type RemoteScreencastMetaPayload struct {
	Format            string  `json:"format,omitempty"`
	ViewportWidth     int     `json:"viewport_width"`
	ViewportHeight    int     `json:"viewport_height"`
	FrameWidth        int     `json:"frame_width"`
	FrameHeight       int     `json:"frame_height"`
	DeviceScaleFactor float64 `json:"device_scale_factor,omitempty"`
}

type RemoteAckPayload struct {
	RefType string `json:"ref_type,omitempty"`
}

type RemoteErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	RefType string `json:"ref_type,omitempty"`
}

type remoteHandlerError struct {
	Code    string
	Message string
}

func (e *remoteHandlerError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func newRemoteHandlerError(code, message string) *remoteHandlerError {
	return &remoteHandlerError{Code: code, Message: message}
}

var DefaultWebSocketUpgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RemoteController struct {
	browser  Browser
	id       string
	upgrader *websocket.Upgrader

	writeMu sync.Mutex
}

func NewRemoteController(browser Browser, id string, upgrader *websocket.Upgrader) *RemoteController {
	if upgrader == nil {
		upgrader = DefaultWebSocketUpgrader
	}
	return &RemoteController{
		browser:  browser,
		id:       id,
		upgrader: upgrader,
	}
}

func (controller *RemoteController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opts, err := NewScreencastOptionsFromUrlValues(r.URL.Query())
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := controller.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	frames, err := controller.browser.StartScreencast(controller.id, opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer frames.Close()

	if err := controller.writeSessionBootstrap(conn, opts); err != nil {
		log.Println(err)
		return
	}

	go func() {
		for frame := range frames.BytesChan {
			if err := controller.writeBinary(conn, frame); err != nil {
				log.Println(err)
				return
			}
		}
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if messageType != websocket.TextMessage {
			continue
		}
		if err := controller.handleTextMessage(conn, data); err != nil {
			log.Println(err)
			return
		}
	}
}

func (controller *RemoteController) writeSessionBootstrap(conn *websocket.Conn, opts *ScreencastOptions) error {
	if err := controller.writeEnvelope(conn, NewRemoteEnvelope(RemoteTypeSessionReady, &RemoteSessionReadyPayload{
		TabID:           controller.id,
		ProtocolVersion: RemoteProtocolVersion,
	})); err != nil {
		return err
	}

	meta, err := controller.browser.GetScreencastSessionMeta(controller.id, opts)
	if err != nil {
		return err
	}
	return controller.writeEnvelope(conn, NewRemoteEnvelope(RemoteTypeScreencastMeta, &RemoteScreencastMetaPayload{
		Format:            meta.Format,
		ViewportWidth:     meta.ViewportWidth,
		ViewportHeight:    meta.ViewportHeight,
		FrameWidth:        meta.FrameWidth,
		FrameHeight:       meta.FrameHeight,
		DeviceScaleFactor: meta.DeviceScaleFactor,
	}))
}

func (controller *RemoteController) handleTextMessage(conn *websocket.Conn, data []byte) error {
	responses := ProcessRemoteControlMessage(controller.browser, controller.id, data)
	for _, response := range responses {
		if response == nil {
			continue
		}
		if err := controller.writeEnvelope(conn, response); err != nil {
			return err
		}
	}
	return nil
}

func (controller *RemoteController) writeBinary(conn *websocket.Conn, frame []byte) error {
	controller.writeMu.Lock()
	defer controller.writeMu.Unlock()
	return conn.WriteMessage(websocket.BinaryMessage, frame)
}

func (controller *RemoteController) writeEnvelope(conn *websocket.Conn, envelope *RemoteEnvelope) error {
	controller.writeMu.Lock()
	defer controller.writeMu.Unlock()
	return conn.WriteJSON(envelope)
}

func NewRemoteEnvelope(msgType string, payload any) *RemoteEnvelope {
	envelope := &RemoteEnvelope{
		V:    RemoteProtocolVersion,
		Type: msgType,
		Ts:   time.Now().UnixMilli(),
	}
	if payload == nil {
		return envelope
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return envelope
	}
	envelope.Payload = raw
	return envelope
}

func ParseRemoteEnvelope(data []byte) (*RemoteEnvelope, error) {
	var envelope RemoteEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, newRemoteHandlerError(RemoteErrorCodeInvalidJSON, err.Error())
	}
	if envelope.Type == "" {
		return nil, newRemoteHandlerError(RemoteErrorCodeInvalidPayload, "type is required")
	}
	if envelope.V == 0 {
		envelope.V = RemoteProtocolVersion
	}
	return &envelope, nil
}

func DecodeRemotePayload[T any](payload json.RawMessage) (*T, error) {
	var value T
	if len(payload) == 0 || string(payload) == "null" {
		return &value, nil
	}
	if err := json.Unmarshal(payload, &value); err != nil {
		return nil, newRemoteHandlerError(RemoteErrorCodeInvalidPayload, err.Error())
	}
	return &value, nil
}

// ProcessRemoteControlMessage parses and dispatches one inbound control message.
// It returns zero or more outbound envelopes (pong/ack/error).
func ProcessRemoteControlMessage(b Browser, tabID string, data []byte) []*RemoteEnvelope {
	envelope, err := ParseRemoteEnvelope(data)
	if err != nil {
		return []*RemoteEnvelope{newRemoteErrorEnvelope("", "", err)}
	}

	if err := DispatchRemoteEnvelope(b, tabID, envelope); err != nil {
		return []*RemoteEnvelope{newRemoteErrorEnvelope(envelope.ID, envelope.Type, err)}
	}

	switch envelope.Type {
	case RemoteTypeSessionPing:
		pong := NewRemoteEnvelope(RemoteTypeSessionPong, map[string]any{})
		pong.ID = envelope.ID
		return []*RemoteEnvelope{pong}
	default:
		if envelope.ID == "" {
			return nil
		}
		ack := NewRemoteEnvelope(RemoteTypeAck, &RemoteAckPayload{RefType: envelope.Type})
		ack.ID = envelope.ID
		return []*RemoteEnvelope{ack}
	}
}

func DispatchRemoteEnvelope(b Browser, tabID string, envelope *RemoteEnvelope) error {
	if envelope == nil {
		return newRemoteHandlerError(RemoteErrorCodeInvalidPayload, "envelope is nil")
	}

	switch envelope.Type {
	case RemoteTypeSessionPing:
		return nil
	case RemoteTypeMouseMove, RemoteTypeMouseDown, RemoteTypeMouseUp, RemoteTypeMouseWheel:
		return dispatchRemoteMouse(b, tabID, envelope)
	case RemoteTypeKeyDown, RemoteTypeKeyUp, RemoteTypeKeyPress:
		return dispatchRemoteKey(b, tabID, envelope)
	case RemoteTypeNavGo, RemoteTypeNavBack, RemoteTypeNavForward, RemoteTypeNavReload:
		return newRemoteHandlerError(RemoteErrorCodeUnsupported, fmt.Sprintf("type %q is not implemented", envelope.Type))
	default:
		if strings.HasPrefix(envelope.Type, "nav.") || strings.HasPrefix(envelope.Type, "event.") {
			return newRemoteHandlerError(RemoteErrorCodeUnsupported, fmt.Sprintf("type %q is not implemented", envelope.Type))
		}
		return newRemoteHandlerError(RemoteErrorCodeUnsupported, fmt.Sprintf("unsupported type %q", envelope.Type))
	}
}

func dispatchRemoteMouse(b Browser, tabID string, envelope *RemoteEnvelope) error {
	payload, err := DecodeRemotePayload[RemoteMousePayload](envelope.Payload)
	if err != nil {
		return err
	}
	if payload.X == nil || payload.Y == nil {
		return newRemoteHandlerError(RemoteErrorCodeInvalidPayload, "x and y are required")
	}

	eventType := ""
	switch envelope.Type {
	case RemoteTypeMouseMove:
		eventType = MouseEventTypeMove
	case RemoteTypeMouseDown:
		eventType = MouseEventTypeDown
	case RemoteTypeMouseUp:
		eventType = MouseEventTypeUp
	case RemoteTypeMouseWheel:
		eventType = MouseEventTypeWheel
	}

	button := payload.Button
	if button == "" && (eventType == MouseEventTypeDown || eventType == MouseEventTypeUp) {
		button = MouseButtonLeft
	}

	if err := b.DispatchMouseEvent(tabID, &MouseEvent{
		Type:       eventType,
		X:          *payload.X,
		Y:          *payload.Y,
		Button:     button,
		ClickCount: payload.ClickCount,
		DeltaX:     payload.DeltaX,
		DeltaY:     payload.DeltaY,
		Modifiers:  payload.Modifiers,
	}); err != nil {
		return newRemoteHandlerError(RemoteErrorCodeInternal, err.Error())
	}
	return nil
}

func dispatchRemoteKey(b Browser, tabID string, envelope *RemoteEnvelope) error {
	payload, err := DecodeRemotePayload[RemoteKeyPayload](envelope.Payload)
	if err != nil {
		return err
	}
	if payload.Key == "" && payload.Code == "" {
		return newRemoteHandlerError(RemoteErrorCodeInvalidPayload, "key or code is required")
	}

	switch envelope.Type {
	case RemoteTypeKeyDown:
		return dispatchOneKey(b, tabID, KeyEventTypeDown, payload)
	case RemoteTypeKeyUp:
		return dispatchOneKey(b, tabID, KeyEventTypeUp, payload)
	case RemoteTypeKeyPress:
		if err := dispatchOneKey(b, tabID, KeyEventTypeDown, payload); err != nil {
			return err
		}
		return dispatchOneKey(b, tabID, KeyEventTypeUp, payload)
	default:
		return newRemoteHandlerError(RemoteErrorCodeUnsupported, envelope.Type)
	}
}

func dispatchOneKey(b Browser, tabID, eventType string, payload *RemoteKeyPayload) error {
	if err := b.DispatchKeyEvent(tabID, &KeyEvent{
		Type:       eventType,
		Key:        payload.Key,
		Code:       payload.Code,
		Text:       payload.Text,
		Modifiers:  payload.Modifiers,
		AutoRepeat: payload.AutoRepeat,
	}); err != nil {
		return newRemoteHandlerError(RemoteErrorCodeInternal, err.Error())
	}
	return nil
}

func newRemoteErrorEnvelope(id, refType string, err error) *RemoteEnvelope {
	code := RemoteErrorCodeInternal
	message := "internal error"
	var handlerErr *remoteHandlerError
	if errors.As(err, &handlerErr) {
		code = handlerErr.Code
		message = handlerErr.Message
	} else if err != nil {
		message = err.Error()
	}

	envelope := NewRemoteEnvelope(RemoteTypeError, &RemoteErrorPayload{
		Code:    code,
		Message: message,
		RefType: refType,
	})
	envelope.ID = id
	return envelope
}
