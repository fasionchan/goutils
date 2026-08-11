package httpx

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)
 
 func TestParseContentRange(t *testing.T) {
	 for _, tc := range []struct {
		 name              string
		 value             string
		 wantUnit          string
		 wantStart         int64
		 wantEnd           int64
		 wantTotal         int64
		 wantErr           bool
		 wantNotFound      bool // 语法不符合：GenericNotFoundError("Content-Range", ...)
		 wantNumParseError bool // 数字段非法：*strconv.NumError
	 }{
		 {
			 name:      "206 常见 Content-Range",
			 value:     "bytes 0-499/1234",
			 wantUnit:  "bytes",
			 wantStart: 0,
			 wantEnd:   499,
			 wantTotal: 1234,
		 },
		 {
			 name:      "单字节区间",
			 value:     "bytes 0-0/1",
			 wantUnit:  "bytes",
			 wantStart: 0,
			 wantEnd:   0,
			 wantTotal: 1,
		 },
		 {
			 name:      "总长为星号",
			 value:     "bytes 0-499/*",
			 wantUnit:  "bytes",
			 wantStart: 0,
			 wantEnd:   499,
			 wantTotal: -1,
		 },
		 {
			 name:      "416 形式 bytes */N",
			 value:     "bytes */1234",
			 wantUnit:  "bytes",
			 wantStart: -1,
			 wantEnd:   -1,
			 wantTotal: 1234,
		 },
		 {
			 name:         "缺少空格",
			 value:        "bytes0-1/10",
			 wantErr:      true,
			 wantNotFound: true,
		 },
		 {
			 name:         "缺少斜杠",
			 value:        "bytes 0-1",
			 wantErr:      true,
			 wantNotFound: true,
		 },
		 {
			 name:         "range 段缺少连字符",
			 value:        "bytes 01/10",
			 wantErr:      true,
			 wantNotFound: true,
		 },
		 {
			 name:              "总长非数字",
			 value:             "bytes 0-1/x",
			 wantErr:           true,
			 wantNumParseError: true,
		 },
		 {
			 name:              "start 非数字",
			 value:             "bytes a-1/10",
			 wantErr:           true,
			 wantNumParseError: true,
		 },
		 {
			 name:              "end 非数字",
			 value:             "bytes 0-b/10",
			 wantErr:           true,
			 wantNumParseError: true,
		 },
		 {
			 name:         "空串",
			 value:        "",
			 wantErr:      true,
			 wantNotFound: true,
		 },
	 } {
		 t.Run(tc.name, func(t *testing.T) {
			 unit, start, end, total, err := ParseContentRange(tc.value)
			 if tc.wantErr {
				 assert.True(t, err != nil)
				 if tc.wantNumParseError {
					 var ne *strconv.NumError
					 assert.True(t, errors.As(err, &ne))
				 }
				 return
			 }
			 assert.NoError(t, err)
			 assert.Equal(t, tc.wantUnit, unit)
			 assert.Equal(t, tc.wantStart, start)
			 assert.Equal(t, tc.wantEnd, end)
			 assert.Equal(t, tc.wantTotal, total)
		 })
	 }
 }