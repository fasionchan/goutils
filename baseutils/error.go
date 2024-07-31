/*
 * Author: fasion
 * Created time: 2022-11-02 21:47:37
 * Last Modified by: fasion
 * Last Modified time: 2024-07-22 13:17:44
 */

package baseutils

import (
	"fmt"
	"reflect"
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

type BadTypeError struct {
	expected string
	given    string
}

func NewBlankBadTypeError() BadTypeError {
	return BadTypeError{}
}

func NewBadTypeError(expected, given string) BadTypeError {
	return BadTypeError{
		expected: expected,
		given:    given,
	}
}

func (err BadTypeError) WithExpected(expected string) BadTypeError {
	err.expected = expected
	return err
}

func (err BadTypeError) WithExpectedReflectType(expected reflect.Type) BadTypeError {
	err.expected = expected.String()
	return err
}

func (err BadTypeError) WithGiven(given string) BadTypeError {
	err.given = given
	return err
}

func (err BadTypeError) WithGivenReflectType(given reflect.Type) BadTypeError {
	err.given = given.String()
	return err
}

func NewBadTypeErrorWithGivenReflectType(expected, given reflect.Type) BadTypeError {
	return NewBadTypeError(expected.String(), given.String())
}

func (e BadTypeError) Error() string {
	return fmt.Sprintf("bad type error: [%s] expected, but given [%s]", e.expected, e.given)
}

type NilError struct {
	name string
}

func NewNilError(name string) NilError {
	return NilError{name: name}
}

func NewNilErrorFromNames(names ...string) NilError {
	return NewNilError(strings.Join(names, "."))
}

func (err NilError) Error() string {
	return fmt.Sprintf("%s is nil", err.name)
}

type Errors = stl.Errors
