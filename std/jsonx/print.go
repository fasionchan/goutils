/*
 * Author: fasion
 * Created time: 2026-03-07 12:43:04
 * Last Modified by: fasion
 * Last Modified time: 2026-03-07 15:24:37
 */

package jsonx

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func FprintPro(w io.Writer, indent, prefix, prefixFmt, head, headFmt, tail, tailFmt string, args ...any) {
	encoder := json.NewEncoder(w)

	for i, arg := range args {
		_prefix := prefix
		if prefixFmt != "" {
			_prefix = fmt.Sprintf(prefixFmt, i)
		}

		_head := head
		if headFmt != "" {
			_head = fmt.Sprintf(headFmt, i)
		}

		_tail := tail
		if tailFmt != "" {
			_tail = fmt.Sprintf(tailFmt, i)
		}

		if _head != "" {
			w.Write([]byte(_head))
		}

		if _prefix != "" {
			w.Write([]byte(_prefix))
		}

		if indent != "" {
			encoder.SetIndent(_prefix, indent)
		}

		if err := encoder.Encode(arg); err != nil {
			fmt.Fprintf(w, "\n%s%s", _prefix, err)
			return
		}

		if _tail != "" {
			w.Write([]byte(_tail))
		}
	}
}

func PrintPro(indent, prefix, prefixFmt, head, headFmt, tail, tailFmt string, args ...any) {
	FprintPro(os.Stdout, indent, prefix, prefixFmt, head, headFmt, tail, tailFmt, args...)
}

func Print(args ...any) {
	FprintPro(os.Stdout, "  ", "", "[%d] ", "", "", "\n", "", args...)
}

func GetPrintFunc(indent, prefix, prefixFmt, head, headFmt, tail, tailFmt string) func(args ...any) {
	return func(args ...any) {
		FprintPro(os.Stdout, indent, prefix, prefixFmt, head, headFmt, tail, tailFmt, args...)
	}
}
