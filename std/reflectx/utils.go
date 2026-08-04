package reflectx

import (
	"fmt"
	"iter"
	"reflect"

	"github.com/fasionchan/goutils/stl"
)

func Len(data any) (int, bool) {
	return ValueLen(reflect.ValueOf(data))
}

func ValueLen(value reflect.Value) (int, bool) {
	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		return value.Len(), true
	case reflect.Pointer:
		if value.IsNil() {
			return 0, false
		}
		return ValueLen(value.Elem())
	default:
		return 0, false
	}
}

func Seq(datas any) (iter.Seq[reflect.Value], error) {
	datasValue := reflect.ValueOf(datas)

	for {
		if datasValue.Type().CanSeq() {
			return datasValue.Seq(), nil
		}

		if datasValue.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("datas is sequenceable")
		}

		if datasValue.IsNil() {
			return stl.EmptySeq[reflect.Value], nil
		}

		datasValue = datasValue.Elem()
	}
}

func SeqAsAny(datas any) (iter.Seq[any], error) {
	seq, err := Seq(datas)
	if err != nil {
		return nil, err
	}

	return stl.MapSeq(seq, reflect.Value.Interface), nil
}

func Seq2(datas any) (iter.Seq2[reflect.Value, reflect.Value], error) {
	if datas == nil {
		return stl.EmptySeq2[reflect.Value, reflect.Value], nil
	}

	datasValue := reflect.ValueOf(datas)

	for datasValue.IsValid() {
		if datasValue.Type().CanSeq2() {
			return datasValue.Seq2(), nil
		}

		if datasValue.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("datas is sequenceable")
		}

		if datasValue.IsNil() {
			return stl.EmptySeq2[reflect.Value, reflect.Value], nil
		}

		datasValue = datasValue.Elem()
	}

	return stl.EmptySeq2[reflect.Value, reflect.Value], nil
}

func Seq2AsAnyPairs(datas any) (iter.Seq[stl.KeyValuePair[any, any]], error) {
	if datas == nil {
		return stl.EmptySeq[stl.KeyValuePair[any, any]], nil
	}

	seq2, err := Seq2(datas)
	if err != nil {
		return nil, err
	}

	return stl.Seq2ToSeq(seq2, func(key, value reflect.Value) stl.KeyValuePair[any, any] {
		return stl.KeyValuePair[any, any]{
			Key:   key.Interface(),
			Value: value.Interface(),
		}
	}), nil
}