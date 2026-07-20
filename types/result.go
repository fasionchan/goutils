package types

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/stl"
)

type (
	ResponseResult = TypedResponseResult[any]
)

var (
	NewBlankResponseResult = NewBlankTypedResponseResult[any]

	NewResponseResult              = NewTypedResponseResult[any]
	NewResponseResultFromData      = NewTypedResponseResultFromData[any]
	NewResponseResultFromError     = NewTypedResponseResultFromError[any]
	NewResponseResultFromPagedData = NewTypedResponseResultFromPagedData[any]
)

type TypedResponseResult[T any] struct {
	Success    bool        `json:"Success"`
	Code       int         `json:"Code"`
	Error      string      `json:"Error"`
	Message    string      `json:"Message"`
	Pagination *Pagination `json:"Pagination,omitempty"`
	Data       T           `json:"Data"`
	Version    string      `json:"Version,omitempty"`
}

func NewTypedResponseResult[T any](
	success bool,
	code int,
	err string,
	message string,
	data T,
	pagination *Pagination,
	version string,
) *TypedResponseResult[T] {
	return &TypedResponseResult[T]{
		Success:    success,
		Code:       code,
		Error:      err,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Version:    version,
	}
}

func NewBlankTypedResponseResult[T any]() *TypedResponseResult[T] {
	return &TypedResponseResult[T]{}
}

func NewTypedResponseResultFromError[T any](code int, err error, msg string) *TypedResponseResult[T] {
	errstr := err.Error()
	if msg == "" {
		msg = errstr
	}

	var zero T
	return NewTypedResponseResult[T](false, code, errstr, msg, zero, nil, "")
}

func NewTypedResponseResultFromData[T any](data T) *TypedResponseResult[T] {
	return NewTypedResponseResultFromPagedData(data, nil)
}

func NewTypedResponseResultFromPagedData[T any](datas T, pagination *Pagination) *TypedResponseResult[T] {
	return NewTypedResponseResult(true, 0, "", "", datas, pagination, "")
}

func NewTypedResponseResultFromPagedDatas[Datas ~[]Data, Data any](datas Datas, pagination *Pagination) *TypedResponseResult[Datas] {
	return NewTypedResponseResult(true, 0, "", "", datas, pagination, "")
}

func (result *TypedResponseResult[T]) Dup() *TypedResponseResult[T] {
	return stl.Dup(result)
}

func (result *TypedResponseResult[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result.Dup().WriteHttpResponseWithStatus(w, http.StatusOK)
}

func (result *TypedResponseResult[T]) WriteHttpResponse(w http.ResponseWriter) error {
	return result.WriteHttpResponseWithStatus(w, http.StatusOK)
}

func (result *TypedResponseResult[T]) WriteHttpResponseWithStatus(w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(result)
}

func (result *TypedResponseResult[T]) GetData() (t T) {
	if result == nil {
		return
	}

	return result.Data
}

func (result *TypedResponseResult[T]) GetError() error {
	if result == nil {
		return baseutils.NewNilError("TypedResponseResult.Validate/result")
	}

	if result.Success {
		return nil
	}

	return fmt.Errorf("code: %d |messsage: %s |error: %s", result.Code, result.Message, result.Error)
}

func (result *TypedResponseResult[T]) SetCode(code int) {
	if result != nil {
		result.Code = code
	}
}

func (result *TypedResponseResult[T]) SetData(data T) {
	if result != nil {
		result.Data = data
	}
}

func (result *TypedResponseResult[T]) SetError(err string) {
	if result != nil {
		result.Error = err
	}
}

func (result *TypedResponseResult[T]) SetMessage(msg string) {
	if result != nil {
		result.Message = msg
	}
}

func (result *TypedResponseResult[T]) SetSuccess(success bool) {
	if result != nil {
		result.Success = success
	}
}

func (result TypedResponseResult[T]) SwaggerDoc() map[string]string {
	return map[string]string{}
}

func (result *TypedResponseResult[T]) WithCode(code int) *TypedResponseResult[T] {
	result.SetCode(code)
	return result
}

func (result *TypedResponseResult[T]) WithData(data T) *TypedResponseResult[T] {
	result.SetData(data)
	return result
}

func (result *TypedResponseResult[T]) WithError(err string) *TypedResponseResult[T] {
	result.SetError(err)
	return result
}

func (result *TypedResponseResult[T]) WithMessage(msg string) *TypedResponseResult[T] {
	result.SetMessage(msg)
	return result
}

func (result *TypedResponseResult[T]) WithSuccess(success bool) *TypedResponseResult[T] {
	result.SetSuccess(success)
	return result
}

func (result *TypedResponseResult[T]) TemplateFuncs() template.FuncMap {
	return stl.ConcatMaps(result.TemplateWithSetFuncs(), result.TemplateWithFuncs())
}

func (result *TypedResponseResult[T]) TemplateWithSetFuncs() template.FuncMap {
	return template.FuncMap{
		"setCode":    result.SetCode,
		"setData":    result.SetData,
		"setError":   result.SetError,
		"setMessage": result.SetMessage,
		"setSuccess": result.SetSuccess,
	}
}

func (result *TypedResponseResult[T]) TemplateWithFuncs() template.FuncMap {
	return template.FuncMap{
		"withCode":    result.WithCode,
		"withData":    result.WithData,
		"withError":   result.WithError,
		"withMessage": result.WithMessage,
		"withSuccess": result.WithSuccess,
	}
}
