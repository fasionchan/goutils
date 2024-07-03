/*
 * Author: fasion
 * Created time: 2023-11-23 14:54:24
 * Last Modified by: fasion
 * Last Modified time: 2024-06-25 17:57:23
 */

package stl

import (
	"fmt"
	"strings"
)

type Errors []error

func NewErrors(errs ...error) Errors {
	return errs
}

func (errors Errors) Append(more ...error) Errors {
	return append(errors, more...)
}

func (errors Errors) Concat(mores ...Errors) Errors {
	return ConcatSlicesTo(errors, mores...)
}

func (errors Errors) Len() int {
	return len(errors)
}

func (errors Errors) Empty() bool {
	return errors.Len() == 0
}

func (errors Errors) PurgeZero() Errors {
	return Filter(errors, func(err error) bool {
		return err != nil
	})
}

func (errors Errors) Simplify() error {
	errors = errors.PurgeZero()
	if errors.Empty() {
		return nil
	}
	return errors
}

func (errors Errors) FirstOneOrZero() error {
	return FirstOneOrZero(errors)
}

func (errors Errors) FirstError() error {
	return FindFirstOrZero(errors, func(err error) bool {
		return err != nil
	})
}

func (errors Errors) Error() string {
	chips := MapPro(errors, func(i int, err error, errors Errors) string {
		return fmt.Sprintf("#%d %s", i, err)
	})
	return strings.Join(chips, "\n")
}
