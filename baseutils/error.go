/*
 * Author: fasion
 * Created time: 2022-11-02 21:47:37
 * Last Modified by: fasion
 * Last Modified time: 2023-11-23 14:55:05
 */

package baseutils

import (
	"fmt"

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

type Errors = stl.Errors
