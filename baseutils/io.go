/*
 * Author: fasion
 * Created time: 2025-07-31 14:32:00
 * Last Modified by: fasion
 * Last Modified time: 2025-07-31 14:32:07
 */

package baseutils

import (
	"io"

	"github.com/fasionchan/goutils/stl"
)

type Closers []io.Closer

func NewClosers(closers ...io.Closer) Closers {
	return Closers(closers)
}

func (closers Closers) Append(others ...io.Closer) Closers {
	return append(closers, others...)
}

func (closers Closers) Concat(others ...Closers) Closers {
	return stl.ConcatSlicesTo(closers, others...)
}

func (closers Closers) Close() error {
	var errs stl.Errors = stl.Map(closers, io.Closer.Close)
	return errs.Simplify()
}

func (closers Closers) PurgeNil() Closers {
	return stl.PurgeZero(closers)
}

func (closers Closers) Simplify() io.Closer {
	if len(closers) == 0 {
		return nil
	} else if len(closers) == 1 {
		return closers[0]
	} else {
		return closers
	}
}
