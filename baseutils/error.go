/*
 * Author: fasion
 * Created time: 2022-11-02 21:47:37
 * Last Modified by: fasion
 * Last Modified time: 2022-12-25 17:57:21
 */

package baseutils

import "fmt"

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
