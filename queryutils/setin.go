/*
 * Author: fasion
 * Created time: 2022-11-12 21:45:25
 * Last Modified by: fasion
 * Last Modified time: 2024-01-30 11:33:42
 */

package queryutils

import (
	"context"
	"errors"
	"fmt"

	"github.com/fasionchan/goutils/stl"
)

var BadDataTypeError = errors.New("bad data type")

type CommonSetinerInterface interface {
	WithSetineds(interface{}) (CommonSetinerInterface, error)
	Setin(context.Context, []string) error
}

type ClonableSetinerInterface interface {
	CommonSetinerInterface
	Clone() ClonableSetinerInterface
}

type SetinTester[Data any, DataPtr ~*Data, Datas ~[]DataPtr] func(ctx context.Context, datas Datas, setin string) (bool, error)
type SetinHandler[Data any, DataPtr ~*Data, Datas ~[]DataPtr] func(ctx context.Context, datas Datas) error

// Create A Setin Handler
func NewSetinHandlerPro[
	Data any,
	Datas ~[]*Data,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDataMapping ~map[Key]SubData,
](
	dataSubKeys func(*Data) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKey func(SubData) Key,
	setinHandler func(data *Data, subDataMapping SubDataMapping) *Data,
) SetinHandler[Data, *Data, Datas] {
	return func(ctx context.Context, datas Datas) error {
		return SetinPro(ctx, datas, dataSubKeys, subDataFetcher, subDataKey, setinHandler)
	}
}

func NewSetinHandler[
	Data any,
	Datas ~[]*Data,
	Key comparable,
	Keys ~[]Key,
	SubData interface {
		GetId() Key
	},
	SubDatas ~[]SubData,
	SubDataMapping ~map[Key]SubData,
](
	dataSubKeys func(*Data) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	setinHandler func(data *Data, subDataMapping SubDataMapping) *Data,
) SetinHandler[Data, *Data, Datas] {
	return NewSetinHandlerPro[Data, Datas](dataSubKeys, subDataFetcher, SubData.GetId, setinHandler)
}

func NewReversedSetinHandlerPro[
	Data any,
	DataPtr ~*Data,
	Datas ~[]DataPtr,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDatasMapping ~map[Key]SubDatas,
](
	dataKey func(DataPtr) Key,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) *Data,
) SetinHandler[Data, DataPtr, Datas] {
	return func(ctx context.Context, datas Datas) error {
		return ReversedSetinPro(ctx, datas, dataKey, subDataFetcher, subDataKeys, setinHandler)
	}
}

func NewReversedSetinHandler[
	Data any,
	DataPtr interface {
		~*Data
		GetId() Key
	},
	Datas ~[]DataPtr,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDatasMapping ~map[Key]SubDatas,
](
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) *Data,
) SetinHandler[Data, DataPtr, Datas] {
	return NewReversedSetinHandlerPro[Data, DataPtr, Datas](DataPtr.GetId, subDataFetcher, subDataKeys, setinHandler)
}

func GenericSetinPro[
	Data any,
	Datas ~[]*Data,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDatasMapping ~map[Key]SubDatas,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(*Data) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data *Data, subDatasMapping SubDatasMapping) *Data,
) error {
	keys := stl.MapAndConcat(datas, dataSubKeys)
	keys = Keys(stl.NewSet(keys...).Slice()) // 集合去重

	subDatas, err := subDataFetcher(ctx, keys)
	if err != nil {
		return err
	}

	subDatasMapping := stl.SliceMappingByKeys(subDatas, subDataKeys)
	stl.ForEach(datas, func(data *Data) {
		setinHandler(data, subDatasMapping)
	})

	return nil
}

func SetinPro[
	Data any,
	Datas ~[]*Data,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDataMapping ~map[Key]SubData,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(*Data) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKey func(SubData) Key,
	setinHandler func(data *Data, subDataMapping SubDataMapping) *Data,
) error {
	keys := stl.MapAndConcat(datas, dataSubKeys)
	keys = Keys(stl.NewSet(keys...).Slice()) // 集合去重

	subDatas, err := subDataFetcher(ctx, keys)
	if err != nil {
		return err
	}

	subDataMapping := stl.MappingByKey(subDatas, subDataKey)
	stl.ForEach(datas, func(data *Data) {
		setinHandler(data, subDataMapping)
	})

	return nil
}

func Setin[
	Data any,
	Datas ~[]*Data,
	Key comparable,
	Keys ~[]Key,
	SubData interface {
		GetId() Key
	},
	SubDatas ~[]SubData,
	SubDataMapping ~map[Key]SubData,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(*Data) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	setinHandler func(data *Data, subDataMapping SubDataMapping) *Data,
) error {
	return SetinPro(ctx, datas, dataSubKeys, subDataFetcher, SubData.GetId, setinHandler)
}

// 调整类型参数顺序，方便指定
func SetinX[
	Datas ~[]*Data,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	Data any,
	Key comparable,
	SubData interface {
		GetId() Key
	},
	SubDataMapping ~map[Key]SubData,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(*Data) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	setinHandler func(data *Data, subDataMapping SubDataMapping) *Data,
) error {
	return Setin(ctx, datas, dataSubKeys, subDataFetcher, setinHandler)
}

func ReversedSetinPro[
	Data any,
	DataPtr ~*Data,
	Datas ~[]DataPtr,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDatasMapping ~map[Key]SubDatas,
](
	ctx context.Context,
	datas Datas,
	dataKey func(DataPtr) Key,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) *Data,
) error {
	keys := stl.Map(datas, dataKey)
	keys = stl.NewSet(keys...).Slice() // 集合去重

	subDatas, err := subDataFetcher(ctx, keys)
	if err != nil {
		return err
	}

	subDatasMapping := stl.SliceMappingByKeys(subDatas, subDataKeys)
	stl.ForEach(datas, func(data DataPtr) {
		setinHandler(data, subDatasMapping)
	})

	return nil
}

func ReversedSetin[
	Data interface {
	},
	DataPtr interface {
		~*Data
		GetId() Key
	},
	Datas ~[]DataPtr,
	Key comparable,
	Keys ~[]Key,
	SubData any,
	SubDatas ~[]SubData,
	SubDatasMapping ~map[Key]SubDatas,
](
	ctx context.Context,
	datas Datas,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) *Data,
) error {
	keys := stl.Map(datas, DataPtr.GetId)
	keys = stl.NewSet(keys...).Slice() // 集合去重

	subDatas, err := subDataFetcher(ctx, keys)
	if err != nil {
		return err
	}

	subDatasMapping := stl.SliceMappingByKeys(subDatas, subDataKeys)
	stl.ForEach(datas, func(data DataPtr) {
		setinHandler(data, subDatasMapping)
	})

	return nil
}

// 调整类型参数顺序，方便指定
func ReversedSetinX[
	Datas ~[]DataPtr,
	Keys ~[]Key,
	SubDatas ~[]SubData,
	Data interface {
	},
	DataPtr interface {
		~*Data
		GetId() Key
	},
	Key comparable,
	SubData any,
	SubDatasMapping ~map[Key]SubDatas,
](
	ctx context.Context,
	datas Datas,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) *Data,
) error {
	return ReversedSetin(ctx, datas, subDataFetcher, subDataKeys, setinHandler)
}

type SetinHandlerMapping[Data any, Datas ~[]*Data] map[string]SetinHandler[Data, *Data, Datas]

func NewSetinHandlerMapping[Data any, Datas ~[]*Data]() SetinHandlerMapping[Data, Datas] {
	return SetinHandlerMapping[Data, Datas]{}
}

func (mapping SetinHandlerMapping[Data, Datas]) Self() SetinHandlerMapping[Data, Datas] {
	return mapping
}

func (mapping SetinHandlerMapping[Data, Datas]) WithHandler(name string, handler SetinHandler[Data, *Data, Datas]) SetinHandlerMapping[Data, Datas] {
	mapping[name] = handler
	return mapping
}

func (mapping SetinHandlerMapping[Data, Datas]) SetinOne(ctx context.Context, datas Datas, setin string) (bool, error) {
	handler, ok := mapping[setin]
	if !ok {
		return false, nil
	}

	if err := handler(ctx, datas); err != nil {
		return false, err
	}

	return true, nil
}

func (mapping SetinHandlerMapping[Data, Datas]) Setin(ctx context.Context, datas Datas, setins []string) error {
	for _, setin := range setins {
		if ok, err := mapping.SetinOne(ctx, datas, setin); err != nil {
			return err
		} else if !ok {
			return NewUnknownSetinError(setin)
		}
	}

	return nil
}

func (mapping SetinHandlerMapping[Data, Datas]) NewTesters() SetinTesters[Data, Datas] {
	return NewSetinTesters(mapping.SetinOne)
}

func (mapping SetinHandlerMapping[Data, Datas]) NewSetiner() *Setiner[Data, Datas, []Data] {
	return mapping.NewTesters().NewSetiner()
}

type UnknownSetinError struct {
	setin string
}

func NewUnknownSetinError(setin string) UnknownSetinError {
	return UnknownSetinError{
		setin: setin,
	}
}

func (e UnknownSetinError) Error() string {
	return fmt.Sprintf("unknown setin: %s", e.setin)
}

type SetinTesters[Data any, Datas ~[]*Data] []SetinTester[Data, *Data, Datas]

func NewSetinTesters[Data any, Datas ~[]*Data](testers ...SetinTester[Data, *Data, Datas]) SetinTesters[Data, Datas] {
	return testers
}

func (testers SetinTesters[Data, Datas]) Append(more ...SetinTester[Data, *Data, Datas]) SetinTesters[Data, Datas] {
	return append(testers, more...)
}

func (testers SetinTesters[Data, Datas]) SetinOne(ctx context.Context, datas Datas, setin string) error {
	for _, tester := range testers {
		if ok, err := tester(ctx, datas, setin); err != nil {
			return err
		} else if ok {
			return nil
		}
	}

	return NewUnknownSetinError(setin)
}

func (testers SetinTesters[Data, Datas]) Setin(ctx context.Context, datas Datas, setins []string) error {
	for _, setin := range setins {
		if err := testers.SetinOne(ctx, datas, setin); err != nil {
			return err
		}
	}

	return nil
}

func (testers SetinTesters[Data, Datas]) NewSetiner() *Setiner[Data, Datas, []Data] {
	return NewSetiner[Data, Datas, []Data](testers)
}

type Setiner[Data any, Datas ~[]*Data, DataInstances ~[]Data] struct {
	SetinTesters[Data, Datas]
}

func NewSetiner[Data any, Datas ~[]*Data, DataInstances ~[]Data](testers SetinTesters[Data, Datas]) *Setiner[Data, Datas, DataInstances] {
	return &Setiner[Data, Datas, DataInstances]{
		SetinTesters: testers,
	}
}

func (setiner *Setiner[Data, Datas, DataInstances]) SetinFor(ctx context.Context, dataX any, setins []string) error {
	if dataX == nil {
		return nil
	}

	if data, ok := dataX.(*Data); ok {
		return setiner.SetinForData(ctx, data, setins)
	} else if datas, ok := dataX.(Datas); ok {
		return setiner.SetinForDatas(ctx, datas, setins)
	} else if datasPtr, ok := dataX.(*Datas); ok {
		return setiner.SetinForDatas(ctx, *datasPtr, setins)
	} else if instances, ok := dataX.(DataInstances); ok {
		return setiner.SetinForDataInstances(ctx, instances, setins)
	}

	return BadDataTypeError
}

func (setiner *Setiner[Data, Datas, DataInstances]) SetinForData(ctx context.Context, data *Data, setins []string) error {
	return setiner.Setin(ctx, Datas{data}, setins)
}

func (setiner *Setiner[Data, Datas, DataInstances]) SetinForDatas(ctx context.Context, datas Datas, setins []string) error {
	return setiner.Setin(ctx, datas, setins)
}

func (setiner *Setiner[Data, Datas, DataInstances]) SetinForDataInstances(ctx context.Context, instances DataInstances, setins []string) error {
	return setiner.Setin(ctx, stl.GetSliceElemPointers(instances), setins)
}

func (setiner *Setiner[Data, Datas, DataInstances]) Action() *SetinAction[Data, Datas, DataInstances] {
	return setiner.ActionFor(nil)
}

func (setiner *Setiner[Data, Datas, DataInstances]) ActionFor(datas Datas) *SetinAction[Data, Datas, DataInstances] {
	return NewSetinAction(setiner, datas)
}

// A SetinAction is A Setiner with datas, it can Dup, Clone, WithData and finally Setin.
type SetinAction[Data any, Datas ~[]*Data, DataInstances ~[]Data] struct {
	*Setiner[Data, Datas, DataInstances]
	datas  Datas
	setins []string
}

func NewSetinAction[Data any, Datas ~[]*Data, DataInstances ~[]Data](setiner *Setiner[Data, Datas, DataInstances], datas Datas) *SetinAction[Data, Datas, DataInstances] {
	return &SetinAction[Data, Datas, DataInstances]{
		Setiner: setiner,
		datas:   datas,
	}
}

func (setiner *SetinAction[Data, Datas, DataInstances]) Dup() *SetinAction[Data, Datas, DataInstances] {
	dup := *setiner
	return &dup
}

func (setiner *SetinAction[Data, Datas, DataInstances]) Clone() ClonableSetinerInterface {
	return setiner
}

func (action *SetinAction[Data, Datas, DataInstances]) WithSetinsX(setins ...string) *SetinAction[Data, Datas, DataInstances] {
	action.setins = setins
	return action
}

func (action *SetinAction[Data, Datas, DataInstances]) WithSetins(setins ...string) CommonSetinerInterface {
	action.setins = setins
	return action
}

func (action *SetinAction[Data, Datas, DataInstances]) WithSetineds(setineds any) (CommonSetinerInterface, error) {
	return action.WithDataX(setineds)
}

func (action *SetinAction[Data, Datas, DataInstances]) WithDataX(dataX any) (*SetinAction[Data, Datas, DataInstances], error) {
	if data, ok := dataX.(*Data); ok {
		return action.WithData(data), nil
	} else if datas, ok := dataX.(Datas); ok {
		return action.WithDatas(datas), nil
	} else if datasPtr, ok := dataX.(*Datas); ok {
		return action.WithDatas(*datasPtr), nil
	} else if instances, ok := dataX.(DataInstances); ok {
		return action.WithDataInstances(instances), nil
	}

	return nil, BadDataTypeError
}

func (action *SetinAction[Data, Datas, DataInstances]) WithData(data *Data) *SetinAction[Data, Datas, DataInstances] {
	if data == nil {
		action.datas = nil
	} else {
		action.datas = Datas{data}
	}
	return action
}

func (action *SetinAction[Data, Datas, DataInstances]) WithDatas(datas Datas) *SetinAction[Data, Datas, DataInstances] {
	action.datas = stl.PurgeZero(datas)
	return action
}

func (action *SetinAction[Data, Datas, DataInstances]) WithDataInstances(datas DataInstances) *SetinAction[Data, Datas, DataInstances] {
	action.datas = stl.GetSliceElemPointers(datas)
	return action
}

func (action *SetinAction[Data, Datas, DataInstances]) SetinX(ctx context.Context, setins ...string) error {
	return action.Setin(ctx, setins)
}

func (action *SetinAction[Data, Datas, DataInstances]) Setin(ctx context.Context, setins []string) error {
	if len(setins) == 0 {
		setins = action.setins
	}
	return action.Setiner.SetinForDatas(ctx, action.datas, setins)
}

func (action *SetinAction[Data, Datas, DataInstances]) SetinForData(ctx context.Context, data *Data) error {
	return action.Setiner.SetinForData(ctx, data, action.setins)
}

func (action *SetinAction[Data, Datas, DataInstances]) SetinForDatas(ctx context.Context, datas Datas) error {
	return action.Setiner.SetinForDatas(ctx, datas, action.setins)
}

func (action *SetinAction[Data, Datas, DataInstances]) SetinForDataInstances(ctx context.Context, instances DataInstances) error {
	return action.Setiner.SetinForDataInstances(ctx, instances, action.setins)
}

func (action *SetinAction[Data, Datas, DataInstances]) SetinFor(ctx context.Context, dataX any) error {
	return action.Setiner.SetinFor(ctx, dataX, action.setins)
}
