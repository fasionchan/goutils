package reflectx

import "reflect"

func Len(data any) (int, bool) {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		return value.Len(), true
	default:
		return 0, false
	}
}