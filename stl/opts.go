package stl

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

type Options[Data any] []func(Data)

func NewOptions[Data any](opts ...func(Data)) Options[Data] {
	return opts
}

func (options Options[Data]) Apply(data Data) Data {
	for _, option := range options {
		option(data)
	}
	return data
}
