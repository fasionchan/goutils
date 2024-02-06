/*
 * Author: fasion
 * Created time: 2024-02-06 10:34:54
 * Last Modified by: fasion
 * Last Modified time: 2024-02-06 10:36:05
 */

package stl

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
](handler func(*Data, Arg) *Data) func(*Data, Arg) {
	return func(data *Data, arg Arg) {
		handler(data, arg)
	}
}
