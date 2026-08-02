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
		return basic.ParseNumber(data.(string), typ)
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
