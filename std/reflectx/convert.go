package reflectx

import (
	"fmt"
	"reflect"

	"github.com/fasionchan/goutils/basic"
)

func ConvertToNumber(data any, typ string) (any, error) {
	value := reflect.ValueOf(data)

	switch value.Kind() {
	case reflect.String:
		str, ok := data.(string)
		if !ok {
			str = value.Convert(StringType).Interface().(string)
		}

		return basic.ParseNumber(str, typ)
	}

	target, ok := TypeMapping[typ]
	if !ok {
		return nil, fmt.Errorf("unsupported type: %s", typ)
	}

	if !value.CanConvert(target) {
		return nil, fmt.Errorf("cannot convert %s to %s", value.Type().Name(), target.Name())
	}

	return value.Convert(target).Interface(), nil
}
