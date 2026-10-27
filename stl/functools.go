/*
 * Author: fasion
 * Created time: 2024-02-06 10:34:54
 * Last Modified by: fasion
 * Last Modified time: 2024-07-16 09:37:33
 */

package stl

type UnaryFuncs[Arg any] []func(Arg)

func NewUnaryFuncs[Arg any](funcs ...func(Arg)) UnaryFuncs[Arg] {
	return funcs
}

func (funcs UnaryFuncs[Arg]) Call(arg Arg) {
	for _, fn := range funcs {
		fn(arg)
	}
}

type UnaryWithErrorFuncs[Arg any] []func(Arg) error

func NewUnaryWithErrorFuncs[Arg any](funcs ...func(Arg) error) UnaryWithErrorFuncs[Arg] {
	return funcs
}

func (funcs UnaryWithErrorFuncs[Arg]) Call(arg Arg) (errs Errors) {
	errs = make(Errors, len(funcs))
	for i, fn := range funcs {
		errs[i] = fn(arg)
	}
	return
}

type UnaryWithResultFuncs[Arg any, Result any] []func(Arg) Result

func NewUnaryWithResultFuncs[Arg any, Result any](funcs ...func(Arg) Result) UnaryWithResultFuncs[Arg, Result] {
	return funcs
}

func (funcs UnaryWithResultFuncs[Arg, Result]) Call(arg Arg) (results []Result) {
	results = make([]Result, len(funcs))
	for i, fn := range funcs {
		results[i] = fn(arg)
	}
	return
}

type UnaryWithResultAndErrorFuncs[Arg any, Result any] []func(Arg) (Result, error)

func NewUnaryWithResultAndErrorFuncs[Arg any, Result any](funcs ...func(Arg) (Result, error)) UnaryWithResultAndErrorFuncs[Arg, Result] {
	return funcs
}

func (funcs UnaryWithResultAndErrorFuncs[Arg, Result]) Call(arg Arg) (results []Result, errs Errors) {
	results = make([]Result, len(funcs))
	errs = make(Errors, len(funcs))
	for i, fn := range funcs {
		results[i], errs[i] = fn(arg)
	}
	return
}

type BinaryWithErrorFuncs[Arg1 any, Arg2 any] []func(Arg1, Arg2) error

func NewBinaryWithErrorFuncs[Arg1 any, Arg2 any](funcs ...func(Arg1, Arg2) error) BinaryWithErrorFuncs[Arg1, Arg2] {
	return funcs
}

func (funcs BinaryWithErrorFuncs[Arg1, Arg2]) Call(arg1 Arg1, arg2 Arg2) (errs Errors) {
	errs = make(Errors, len(funcs))
	for i, fn := range funcs {
		errs[i] = fn(arg1, arg2)
	}
	return
}

func ChainUnaryHandler[
	Data any,
	Arg any,
](handler func(*Data, Arg)) func(*Data, Arg) *Data {
	return func(data *Data, arg Arg) *Data {
		handler(data, arg)
		return data
	}
}

func UnchainUnaryHandler[
	Data any,
	Arg any,
](handler func(Data, Arg) Data) func(Data, Arg) {
	return func(data Data, arg Arg) {
		handler(data, arg)
	}
}
