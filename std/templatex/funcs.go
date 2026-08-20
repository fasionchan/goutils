package templatex

import (
	"time"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/std/_time"
	"github.com/fasionchan/goutils/std/jsonx"
)

var TemplateFuncs = TemplateFuncMap{
	"now":       time.Now,
	"yesterday": baseutils.Yesterday,
	"today":     baseutils.Today,
	"tomorrow":  baseutils.Tomorrow,

	"jsonDecode": jsonx.JsonDecodeAnyToAny,
	"jsonEncode": jsonx.JsonEncodeAnyToString,

	"durationFromNative":      _time.DurationFromNative,
	"durationLocaleString":    _time.DurationLocaleString,
	"durationLocaleStringPro": _time.DurationLocaleStringPro,

	"nativeDurationLocaleString":    _time.NativeDurationLocaleString,
	"nativeDurationLocaleStringPro": _time.NativeDurationLocaleStringPro,
}
