package datarender

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type DataParserFunc[Data any] func(text string) (Data, error)

type DataParserFuncMappingByString[Data any] map[string]DataParserFunc[Data]

func NewDataParserFuncMappingByString[Data any]() DataParserFuncMappingByString[Data] {
	return make(DataParserFuncMappingByString[Data])
}

func (mapping DataParserFuncMappingByString[Data]) GetParser(format string) DataParserFunc[Data] {
	if format == "" {
		format = DataParserTypeDefault
	}

	parser, ok := mapping[format]
	if !ok {
		// return error if key is not found?
		return func(text string) (data Data, err error) { return }
	}

	return parser
}

func (mapping DataParserFuncMappingByString[Data]) Parse(key, text string) (data Data, err error) {
	// return error if key is not found?
	return mapping.GetParser(key)(text)
}

func (mapping DataParserFuncMappingByString[Data]) With(key string, parser DataParserFunc[Data]) DataParserFuncMappingByString[Data] {
	if mapping != nil {
		mapping[key] = parser
	}
	return mapping
}

func (mapping DataParserFuncMappingByString[Data]) WithDefault(parser DataParserFunc[Data]) DataParserFuncMappingByString[Data] {
	return mapping.With(DataParserTypeDefault, parser)
}

func GenericParseDatasJSON[Datas any](text string) (Datas, error) {
	return GenericParseDatasUnmarshal[Datas](text, json.Unmarshal)
}

func GenericParseDatasYAML[Datas any](text string) (Datas, error) {
	return GenericParseDatasUnmarshal[Datas](text, yaml.Unmarshal)
}

func GenericParseDatasUnmarshal[Datas any](text string, unmarshal func(data []byte, v any) error) (Datas, error) {
	var datas Datas
	err := unmarshal([]byte(text), &datas)
	if err != nil {
		return datas, nil
	}

	return datas, nil
}
