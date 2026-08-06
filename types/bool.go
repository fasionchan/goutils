package types

const (
	BoolOptionNone  = iota // 选项未指定
	BoolOptionTrue         // 选项开启
	BoolOptionFalse        // 选项关闭

	BoolOptionOn  = BoolOptionTrue
	BoolOptionOff = BoolOptionFalse

	BoolOptionEnabled  = BoolOptionTrue
	BoolOptionDisabled = BoolOptionFalse
)

type BoolOption bool

func NewBoolOption(option bool) *BoolOption {
	result := BoolOption(option)
	return &result
}

func (option *BoolOption) Option() int {
	if option == nil {
		return BoolOptionNone
	} else if *option {
		return BoolOptionOn
	} else {
		return BoolOptionOff
	}
}

func (option *BoolOption) IsNone() bool {
	return option == nil
}

func (option *BoolOption) IsTrue() bool {
	if option == nil {
		return false
	}
	return *option == true
}

func (option *BoolOption) IsFalse() bool {
	if option == nil {
		return false
	}
	return *option == false
}

func (option *BoolOption) IsOn() bool {
	return option.IsTrue()
}

func (option *BoolOption) IsOff() bool {
	return option.IsFalse()
}

func (option *BoolOption) IsEnabled() bool {
	return option.IsTrue()
}

func (option *BoolOption) IsDisabled() bool {
	return option.IsFalse()
}