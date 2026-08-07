package text

import (
	"fmt"

	"github.com/fasionchan/goutils/stl"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	MessagePrinterLoaderEnglish = stl.NewLazyDataLoader(NewPrinterLoadFunc(language.English))
	MessagePrinterLoaderChinese = stl.NewLazyDataLoader(NewPrinterLoadFunc(language.Chinese))
)

func NewLanguageTagLoadFunc(lang string) func() (language.Tag, error) {
	return func() (language.Tag, error) {
		return language.Parse(lang)
	}
}

func NewPrinterLoadFunc(tag language.Tag) func() (*message.Printer, error) {
	return func() (*message.Printer, error) {
		return message.NewPrinter(tag), nil
	}
}

func NewPrinterLoadFuncFromTagLoader(langTagLoader *stl.LazyDataLoader[language.Tag]) func() (*message.Printer, error) {
	return func() (*message.Printer, error) {
		langTag, err := langTagLoader.Load()
		if err != nil {
			return nil, err
		}

		return message.NewPrinter(langTag), nil
	}
}

func NewI8nSprintFunc(
	laoder *stl.LazyDataLoader[*message.Printer],
	fn func(printer *message.Printer, key message.Reference, a ...any) string,
	fallback func(format string, a ...any) string,
) func(format string, a ...any) string {
	printer, err := laoder.Load()
	if err != nil {
		return fallback
	}

	return func(format string, a ...any) string {
		return fn(printer, format, a...)
	}
}

func NewI8nSprintFuncChinese() func(format string, a ...any) string {
	return NewI8nSprintFunc(
		MessagePrinterLoaderChinese,
		(*message.Printer).Sprintf,
		fmt.Sprintf,
	)
}

func NewI8nSprintFuncEnglish() func(format string, a ...any) string {
	return NewI8nSprintFunc(
		MessagePrinterLoaderEnglish,
		(*message.Printer).Sprintf,
		fmt.Sprintf,
	)
}