/*
 * Author: fasion
 * Created time: 2024-10-25 14:24:15
 * Last Modified by: fasion
 * Last Modified time: 2026-07-05 13:03:39
 */

package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fasionchan/goutils/stl"
	"github.com/fasionchan/goutils/types"
)

const (
	DefaultServerSentMessageEndOfLine = "\r\n"

	ServerSentEventMessageTypeData  = "data"
	ServerSentEventMessageTypeEvent = "event"
	ServerSentEventMessageTypeId    = "id"
	ServerSentEventMessageTypeRetry = "retry"
)

type ServerSentEventField = stl.KeyValuePair[string, []byte]
type ServerSentEventFieldPtr = *ServerSentEventField

func ParseServerSentEventField(field []byte) (*ServerSentEventField, error) {
	index := bytes.IndexByte(field, ':')
	if index == -1 {
		return nil, fmt.Errorf("bad server sent event field: %s", string(field))
	}

	key := string(field[:index])
	key = strings.TrimSpace(key)
	key = strings.ToLower(key)

	return &ServerSentEventField{
		Key:   key,
		Value: field[index+1:],
	}, nil
}

type ServerSentEventMessageWriter = stl.Writer[ServerSentEventMessages, *ServerSentEventMessage]
type ServerSentEventMessageWriters = stl.Writers[ServerSentEventMessages, *ServerSentEventMessage]

type ServerSentEventMessagePtr = *ServerSentEventMessage

type ServerSentEventMessage struct {
	Id        string
	Event     string
	Retry     string
	Data      []byte
	DataLines [][]byte
	endOfLine string
}

func NewServerSentEventMessage(data []byte, event, id, retry, endOfLine string) *ServerSentEventMessage {
	return &ServerSentEventMessage{
		Data:      data,
		Event:     event,
		Id:        id,
		endOfLine: endOfLine,
	}
}

func NewServerSentEventMessageFromJsonData(data any, event, id, retry, endOfLine string) (*ServerSentEventMessage, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return NewServerSentEventMessage(jsonData, event, id, retry, endOfLine), nil
}

func ParseServerSentEventMessage(message []byte, endOfLine string) (*ServerSentEventMessage, error) {
	var result ServerSentEventMessage

	lines := bytes.Split(message, []byte(endOfLine))
	for _, line := range lines {
		field, err := ParseServerSentEventField(line)
		if err != nil {
			return nil, err
		}

		switch key := field.Key; key {
		case ServerSentEventMessageTypeId:
			result.Id = string(field.Value)
		case ServerSentEventMessageTypeEvent:
			result.Event = string(field.Value)
		case ServerSentEventMessageTypeData:
			result.DataLines = append(result.DataLines, field.Value)
		case ServerSentEventMessageTypeRetry:
		case "":
		default:
			return nil, fmt.Errorf("bad server sent event field: key=%s", key)
		}
	}

	result.Data = bytes.Join(result.DataLines, []byte(endOfLine))

	return &result, nil
}

func (message *ServerSentEventMessage) DataEmpty() bool {
	return len(message.Data) == 0
}

func (message *ServerSentEventMessage) GetData() []byte {
	if message == nil {
		return nil
	}
	return message.Data
}

func (message *ServerSentEventMessage) WriteTo(writer io.Writer) (int, error) {
	if message == nil {
		return 0, nil
	}

	chips := make([][]byte, 0, 13)
	chips = append(chips, []byte(ServerSentEventMessageTypeData), message.Data, []byte(message.endOfLine))

	if message.Id != "" {
		chips = append(chips, []byte(ServerSentEventMessageTypeId), []byte(message.Id), []byte(message.endOfLine))
	}

	if message.Event != "" {
		chips = append(chips, []byte(ServerSentEventMessageTypeEvent), []byte(message.Event), []byte(message.endOfLine))
	}

	chips = append(chips, []byte(message.endOfLine))

	return BatchWrite(writer, chips...)
}

func (message *ServerSentEventMessage) MustPack() []byte {
	data, err := message.Pack()
	if err != nil {
		panic(err)
	}
	return data
}

func (message *ServerSentEventMessage) Pack() ([]byte, error) {
	var buffer bytes.Buffer
	if _, err := message.WriteTo(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (message *ServerSentEventMessage) String() string {
	data := message.MustPack()
	return string(data)
}

// todo move to goutilsx
func BatchWrite(writer io.Writer, datas ...[]byte) (writtenBytes int, err error) {
	for _, data := range datas {
		var n int
		n, err = writer.Write(data)
		if n > 0 {
			writtenBytes += n
		}
		if err != nil {
			return
		}
	}

	return
}

type ServerSentEventMessages []*ServerSentEventMessage

func (messages ServerSentEventMessages) Append(others ...*ServerSentEventMessage) ServerSentEventMessages {
	return append(messages, others...)
}

func (messages ServerSentEventMessages) EventDatas(name string) [][]byte {
	return stl.Map(messages, func(message *ServerSentEventMessage) []byte {
		return message.Data
	})
}

func (messages ServerSentEventMessages) PurgeNil() ServerSentEventMessages {
	return stl.PurgeZero(messages)
}

type ServerSentEventReader struct {
	io.ReadCloser
	message   []byte
	endOfLine string
}

func NewServerSentEventReader(reader io.Reader, endOfLine string) *ServerSentEventReader {
	return NewServerSentEventReadCloser(io.NopCloser(reader), endOfLine)
}

func NewServerSentEventReadCloser(reader io.ReadCloser, endOfLine string) *ServerSentEventReader {
	if endOfLine == "" {
		endOfLine = DefaultServerSentMessageEndOfLine
	}

	return &ServerSentEventReader{
		ReadCloser: reader,
		endOfLine:  endOfLine,
	}
}

func NewServerSentEventReaderFromHttpResponse(response *http.Response, endOfLine string) (*ServerSentEventReader, error) {
	contentType := strings.TrimSpace(response.Header.Get("Content-Type"))
	if contentType != "" {
		if !strings.Contains(strings.ToLower(contentType), "text/event-stream") {
			return nil, fmt.Errorf("content-type is not sse: %s", contentType)
		}
	}

	return NewServerSentEventReader(response.Body, endOfLine), nil
}

func (reader *ServerSentEventReader) WithEndOfLine(endOfLine string) *ServerSentEventReader {
	if reader == nil {
		return nil
	}

	if endOfLine == "" {
		return reader
	}

	reader.endOfLine = endOfLine
	return reader
}

func (reader *ServerSentEventReader) WithReadCloser(readCloser io.ReadCloser) *ServerSentEventReader {
	if reader == nil {
		return nil
	}

	reader.ReadCloser = readCloser
	return reader
}

func (reader *ServerSentEventReader) Read() (*ServerSentEventMessage, error) {
	endOfLine := reader.endOfLine
	if endOfLine == "" {
		endOfLine = DefaultServerSentMessageEndOfLine
	}

	endOfMsg := endOfLine + endOfLine
	endOfMsgLen := len(endOfMsg)

ReadNextStart:
	index := bytes.Index(reader.message, []byte(endOfMsg))
	if index != -1 {
		message := reader.message[:index]
		reader.message = reader.message[index+endOfMsgLen:]

		return ParseServerSentEventMessage(message, endOfLine)
	}

	var buffer [102400]byte
	n, err := reader.ReadCloser.Read(buffer[:])
	if err != nil {
		if err != io.EOF {
			return nil, err
		}

		if n == 0 {
			return nil, err
		}
	}

	reader.message = append(reader.message, buffer[:n]...)
	goto ReadNextStart
}

func (reader *ServerSentEventReader) ReadAll() (ServerSentEventMessages, error) {
	messages, err := stl.ReadAll[ServerSentEventMessages](reader.Read)
	return messages.PurgeNil(), err
}

func ParseJsonEventDatas[Datas ~[]Data, Data any](messages ServerSentEventMessages, name string) (Datas, error) {
	return BatchUnmarshalJsonForDatas[Datas](messages.EventDatas(name)...)
}

func BatchUnmarshalJsonForDatas[Datas ~[]Data, Data any](datas ...[]byte) (Datas, error) {
	return stl.MapUntilError(datas, UnmarshalJsonForData[Data])
}

func UnmarshalJsonForData[Data any](data []byte) (result Data, err error) {
	if len(data) > 0 {
		err = json.Unmarshal(data, &result)
	}
	return
}

type ServerSentEventManagerWriter = stl.Writer[ServerSentEventMessages, *ServerSentEventMessage]

type ServerSentEventParser struct {
	message   []byte
	endOfLine string
	output    ServerSentEventManagerWriter
}

func NewServerSentEventParser(output ServerSentEventManagerWriter) *ServerSentEventParser {
	return &ServerSentEventParser{
		output: output,
	}
}

func NewServerSentEventParserWithWriteFunc(writeFunc stl.NilNopWriteFunc[ServerSentEventMessages, *ServerSentEventMessage]) *ServerSentEventParser {
	return NewServerSentEventParser(writeFunc)
}

func (parser *ServerSentEventParser) WithEndOfLine(endOfLine string) *ServerSentEventParser {
	if parser == nil {
		return nil
	}

	parser.endOfLine = endOfLine

	return parser
}

func (parser *ServerSentEventParser) Write(data []byte) (int, error) {
	parser.message = append(parser.message, data...)
	messages, err := parser.parse()

	if err != nil {
		return 0, err
	}

	if _, err := parser.output.Write(messages); err != nil {
		return 0, err
	}

	return len(data), nil
}

func (parser *ServerSentEventParser) parse() (messages ServerSentEventMessages, err error) {
	for {
		var message *ServerSentEventMessage
		message, err = parser.parseNext()
		if err != nil {
			return
		}

		if message == nil {
			return
		}

		messages = messages.Append(message)
	}
}

func (reader *ServerSentEventParser) parseNext() (*ServerSentEventMessage, error) {
	endOfLine := reader.endOfLine
	if endOfLine == "" {
		endOfLine = DefaultServerSentMessageEndOfLine
	}

	endOfMsg := endOfLine + endOfLine
	endOfMsgLen := len(endOfMsg)

	index := bytes.Index(reader.message, []byte(endOfMsg))
	if index != -1 {
		message := reader.message[:index]
		reader.message = reader.message[index+endOfMsgLen:]

		return ParseServerSentEventMessage(message, endOfLine)
	}

	return nil, nil
}

type ServerSentEventWriter struct {
	endOfLine    []byte
	endOfMessage []byte

	writer  io.Writer
	flusher http.Flusher
}

func NewServerSentEventWriter(writer io.Writer) *ServerSentEventWriter {
	endOfLine := "\n"
	endOfMessage := endOfLine + endOfLine

	flusher, _ := writer.(http.Flusher)

	return &ServerSentEventWriter{
		endOfLine:    []byte(endOfLine),
		endOfMessage: []byte(endOfMessage),
		writer:       writer,
		flusher:      flusher,
	}

}

func ResponseServerSentEvent(response http.ResponseWriter, httpStatus int) *ServerSentEventWriter {
	if httpStatus == 0 {
		httpStatus = http.StatusOK
	}

	response.Header().Set("Content-Type", "text/event-stream")
	response.WriteHeader(httpStatus)

	return NewServerSentEventWriter(response)
}

func (writer *ServerSentEventWriter) WithEndOfLine(endOfLine string) *ServerSentEventWriter {
	writer.endOfLine = []byte(endOfLine)
	writer.endOfMessage = []byte(endOfLine + endOfLine)
	return writer
}

func (writer *ServerSentEventWriter) WriteCommonMessage(common string) error {
	return writer.WriteMessageLite("", []byte(common))
}

func (writer *ServerSentEventWriter) WriteDataMessage(data []byte) error {
	return writer.WriteMessageLite("data", data)
}

func (writer *ServerSentEventWriter) WriteErrorMessage(err error) error {
	return writer.WriteErrorMessageRaw([]byte(err.Error()))
}

func (writer *ServerSentEventWriter) WriteErrorMessageRaw(data []byte) error {
	return writer.WriteMessageLite("error", data)
}

func (writer *ServerSentEventWriter) WriteJsonDataMessage(data interface{}) error {
	return writer.WriteJsonMessage("data", data)
}

func (writer *ServerSentEventWriter) WriteJsonMessage(field string, data any) error {
	if err := writer.writeX(
		[]byte(field),
		[]byte(": "),
	); err != nil {
		return err
	}

	if err := writer.writeJson(data); err != nil {
		return err
	}

	return writer.flushX(writer.endOfMessage)
}

func (writer *ServerSentEventWriter) write(data []byte) error {
	n, err := writer.writer.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("write: wrote %d bytes, expected %d", n, len(data))
	}

	return nil
}

func (writer *ServerSentEventWriter) writeJson(data any) error {
	return json.NewEncoder(writer.writer).Encode(data)
}

func (writer *ServerSentEventWriter) writeX(datas ...[]byte) error {
	if len(datas) == 0 {
		return nil
	}

	return stl.BatchProcessUntilFirstError(datas, writer.write)
}

func (writer *ServerSentEventWriter) flushX(datas ...[]byte) error {
	if err := writer.writeX(datas...); err != nil {
		return err
	}

	if flusher := writer.flusher; flusher != nil {
		flusher.Flush()
	}

	return nil
}

func (writer *ServerSentEventWriter) WriteMessageLite(field string, data []byte) error {
	return writer.flushX(
		[]byte(field),
		[]byte(": "),
		data,
		writer.endOfMessage,
	)
}

func (writer *ServerSentEventWriter) WriteMessage(data []byte, event, id, retry string) error {
	message := NewServerSentEventMessage(data, event, id, retry, string(writer.endOfLine))
	return writer.WriteMessagePack(message)
}

func (writer *ServerSentEventWriter) WriteMessageWithJsonData(data any, event, id, retry string) error {
	message, err := NewServerSentEventMessageFromJsonData(data, event, id, retry, string(writer.endOfLine))
	if err != nil {
		return err
	}
	return writer.WriteMessagePack(message)
}

func (writer *ServerSentEventWriter) WriteMessagePack(message *ServerSentEventMessage) error {
	if _, err := message.WriteTo(writer.writer); err != nil {
		return err
	}

	return writer.flushX(writer.endOfMessage)
}

func (writer *ServerSentEventWriter) WriteApiErrorMessage(code int, err error, message string) error {
	return writer.WriteApiErrorMessageRaw(code, err.Error(), message)
}

func (writer *ServerSentEventWriter) WriteApiErrorMessageRaw(code int, err, message string) error {
	return writer.WriteApiResultMessage(false, code, err, message, nil, nil)
}

func (writer *ServerSentEventWriter) WriteApiResultMessage(success bool, code int, err, message string, data interface{}, pagination *types.Pagination) error {
	return writer.WriteApiResultMessagePack(&types.ResponseResult{
		Success:    success,
		Code:       code,
		Error:      err,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func (writer *ServerSentEventWriter) WriteApiResultMessagePack(result *types.ResponseResult) error {
	return writer.WriteMessageWithJsonData(result, "ApiResult", "", "")
}
