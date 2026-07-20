package stl

import (
	"errors"
	"strconv"
)

type Option[Data any] func(Data)

func NewOption[Data any](opt func(Data)) Option[Data] {
	return opt
}

func (opt Option[Data]) Apply(data Data) Data {
	if opt != nil {
		opt(data)
	}
	return data
}

func (opt Option[Data]) IsNil() bool {
	return opt == nil
}

type Options[Data any] []Option[Data]

func NewOptions[Data any](opts ...Option[Data]) Options[Data] {
	return opts
}

func (opts Options[Data]) Append(more... Option[Data]) Options[Data] {
	return append(opts, more...)
}

func (opts Options[Data]) Apply(data Data) Data {
	ForEachUnaryByMapper(opts, Option[Data].Apply, data)
	return data
}

func (opts Options[Data]) ParseAndAppendIntOptions(metas IntOptionMetas[Data]) (Options[Data], error) {
	intOpts, err := metas.NewOptions()
	if err != nil {
		return nil, err
	}

	return opts.Append(intOpts.PurgeNil()...), nil
}

func (opts Options[Data]) PurgeNil() Options[Data] {
	return Purge(opts, Option[Data].IsNil)
}

type IntOptionMeta[Data any] struct {
	Opt func(int) Option[Data]
	Value string
	Required bool
	Default *int
}

func NewIntOptionMeta[Data any](opt func(int) Option[Data], value string) IntOptionMeta[Data] {
	return IntOptionMeta[Data]{
		Opt: opt,
		Value: value,
	}
}

func (meta IntOptionMeta[Data]) NewOption() (Option[Data], error) {
	value := meta.Value
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		return meta.Opt(intValue), nil
	}

	if meta.Required {
		return nil, errors.New("value is required")
	}

	if defaultValue := meta.Default; defaultValue != nil {
		return meta.Opt(*defaultValue), nil
	}

	return nil, nil
}

type IntOptionMetas[Data any] []IntOptionMeta[Data]

func (metas IntOptionMetas[Data]) NewOptions() (Options[Data], error) {
	opts, errs := MapWithError(metas, true, IntOptionMeta[Data].NewOption)
	if err := errs.Simplify(); err != nil {
		return nil, err
	}

	return opts, nil
}