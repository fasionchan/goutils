/*
 * Author: fasion
 * Created time: 2022-11-02 21:47:37
 * Last Modified by: fasion
 * Last Modified time: 2023-11-08 11:11:17
 */

package baseutils

import (
	"fmt"
	"strings"

	"github.com/fasionchan/goutils/stl"
)

type GenericNotFoundError struct {
	t    string
	name string
}

func NewGenericNotFoundError(t, name string) GenericNotFoundError {
	return GenericNotFoundError{
		t:    t,
		name: name,
	}
}

func (e GenericNotFoundError) Name() string {
	return e.name
}

func (e GenericNotFoundError) Type() string {
	return e.t
}

func (e GenericNotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Type(), e.Name())
}

type EnvironmentVariableNotFoundError struct {
	name string
}

func NewEnvironmentVariableNotFoundError(name string) EnvironmentVariableNotFoundError {
	return EnvironmentVariableNotFoundError{
		name: name,
	}
}

func (e EnvironmentVariableNotFoundError) Name() string {
	return e.name
}

func (e EnvironmentVariableNotFoundError) Error() string {
	return fmt.Sprintf("environment variable not found: %s", e.name)
}

type CookieNotFoundError struct {
	name string
}

func NewCookieNotFoundError(name string) CookieNotFoundError {
	return CookieNotFoundError{
		name: name,
	}
}

func (e CookieNotFoundError) Name() string {
	return e.name
}

func (e CookieNotFoundError) Error() string {
	return fmt.Sprintf("cookie not found: %s", e.name)
}

type NotImplementedError struct {
	hint string
}

func NewNotImplementedError(hint string) NotImplementedError {
	return NotImplementedError{
		hint: hint,
	}
}

func (e NotImplementedError) Error() string {
	return fmt.Sprintf("not implemented: %s", e.hint)
}

type Errors []error

func (errors Errors) Len() int {
	return len(errors)
}

func (errors Errors) Empty() bool {
	return errors.Len() == 0
}

func (errors Errors) PurgeZero() Errors {
	return stl.Filter(errors, func(err error) bool {
		return err != nil
	})
}

func (errors Errors) Error() string {
	chips := stl.MapPro(errors, func(i int, err error, errors Errors) string {
		return fmt.Sprintf("#%d %s", i, err)
	})
	return strings.Join(chips, "\n")
}
