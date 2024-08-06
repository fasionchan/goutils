/*
 * Author: fasion
 * Created time: 2022-11-12 21:45:25
 * Last Modified by: fasion
 * Last Modified time: 2024-08-06 13:32:15
 */

package queryutils

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

type SetinTester[Datas ~[]DataPtr, DataPtr ~*Data, Data any] func(ctx context.Context, datas Datas, setin string) (matched bool, err error)
type SetinHandler[Datas ~[]DataPtr, DataPtr ~*Data, Data any] func(ctx context.Context, datas Datas) error

type NamedSetinHandler[Datas ~[]DataPtr, DataPtr ~*Data, Data any] struct {
	name    string
	handler SetinHandler[Datas, DataPtr, Data]
}

func NewNamedSetinHandler[Datas ~[]DataPtr, DataPtr ~*Data, Data any](name string, handler SetinHandler[Datas, DataPtr, Data]) *NamedSetinHandler[Datas, DataPtr, Data] {
	return &NamedSetinHandler[Datas, DataPtr, Data]{
		name:    name,
		handler: handler,
	}
}

func (handler *NamedSetinHandler[Datas, DataPtr, Data]) GetName() string {
	if handler == nil {
		return ""
	}
	return handler.name
}

func (handler *NamedSetinHandler[Datas, DataPtr, Data]) GetHandler() SetinHandler[Datas, DataPtr, Data] {
	if handler == nil {
		return nil
	}
	return handler.handler
}

func NewSetinerHandlerFromComputedHandler[Datas ~[]DataPtr, DataPtr ~*Data, Data any](handler SetinComputedHandler[Datas, DataPtr, Data]) SetinHandler[Datas, DataPtr, Data] {
	return func(ctx context.Context, datas Datas) error {
		return handler(datas)
	}
}

type SetinComputedHandler[Datas ~[]DataPtr, DataPtr ~*Data, Data any] func(datas Datas) error

func NewSetinComputedHandlerFromSingleton[Datas ~[]DataPtr, DataPtr ~*Data, Data any](handler func(DataPtr) error) SetinComputedHandler[Datas, DataPtr, Data] {
	return func(datas Datas) error {
		return stl.Errors(stl.Map(datas, handler)).Simplify()
	}
}

func PartialUnarySetinHandler[
	Data any,
	Arg any,
](handler func(Data, Arg) Data, arg Arg) func(Data) Data {
	return func(data Data) Data {
		return handler(data, arg)
	}
}

func BatchCallUnarySetinHandler[
	Datas ~[]DataPtr,
	Arg any,
	DataPtr ~*Data,
	Data any,
](datas Datas, handler func(DataPtr, Arg) DataPtr, arg Arg) {
	stl.ForEachByMapper(datas, PartialUnarySetinHandler(handler, arg))
}

// Create A Setin Handler
func NewSetinHandlerPro[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDataMapping ~map[Key]SubData,
	DataPtr ~*Data,
	Data any,
	SubData any,
	Key comparable,
](
	dataSubKeys func(DataPtr) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKey func(SubData) Key,
	setinHandler func(data DataPtr, subDataMapping SubDataMapping) DataPtr,
) SetinHandler[Datas, DataPtr, Data] {
	return func(ctx context.Context, datas Datas) error {
		return SetinPro(ctx, datas, dataSubKeys, subDataFetcher, subDataKey, setinHandler)
	}
}

func NewSetinHandler[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDataMapping ~map[Key]SubData,
	DataPtr ~*Data,
	Data any,
	SubData interface {
		GetId() Key
	},
	Key comparable,
](
	dataSubKeys func(DataPtr) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	setinHandler func(data DataPtr, subDataMapping SubDataMapping) DataPtr,
) SetinHandler[Datas, DataPtr, Data] {
	return NewSetinHandlerPro[Datas](dataSubKeys, subDataFetcher, SubData.GetId, setinHandler)
}

func NewReversedSetinHandlerPro[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDatasMapping ~map[Key]SubDatas,
	DataPtr ~*Data,
	Data any,
	SubData any,
	Key comparable,
](
	dataKey func(DataPtr) Key,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) DataPtr,
) SetinHandler[Datas, DataPtr, Data] {
	return func(ctx context.Context, datas Datas) error {
		return ReversedSetinPro(ctx, datas, dataKey, subDataFetcher, subDataKeys, setinHandler)
	}
}

func NewReversedSetinHandler[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDatasMapping ~map[Key]SubDatas,
	DataPtr interface {
		~*Data
		GetId() Key
	},
	Data any,
	SubData any,
	Key comparable,
](
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) DataPtr,
) SetinHandler[Datas, DataPtr, Data] {
	return NewReversedSetinHandlerPro[Datas](DataPtr.GetId, subDataFetcher, subDataKeys, setinHandler)
}

func GenericSetinPro[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDatasMapping ~map[Key]SubDatas,
	DataPtr ~*Data,
	Data any,
	SubData any,
	Key comparable,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(DataPtr) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) DataPtr,
) error {
	keys := stl.MapAndConcat(datas, dataSubKeys)
	keys = Keys(stl.NewSet(keys...).Slice()) // 集合去重

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

func SetinPro[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDataMapping ~map[Key]SubData,
	DataPtr ~*Data,
	Data any,
	SubData any,
	Key comparable,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(DataPtr) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKey func(SubData) Key,
	setinHandler func(data DataPtr, subDataMapping SubDataMapping) DataPtr,
) error {
	keys := stl.MapAndConcat(datas, dataSubKeys)
	keys = Keys(stl.NewSet(keys...).Slice()) // 集合去重

	subDatas, err := subDataFetcher(ctx, keys)
	if err != nil {
		return err
	}

	subDataMapping := stl.MappingByKey(subDatas, subDataKey)
	BatchCallUnarySetinHandler(datas, setinHandler, subDataMapping)

	return nil
}

func Setin[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDataMapping ~map[Key]SubData,
	DataPtr ~*Data,
	Data any,
	SubData interface {
		GetId() Key
	},
	Key comparable,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(DataPtr) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	setinHandler func(data DataPtr, subDataMapping SubDataMapping) DataPtr,
) error {
	return SetinPro(ctx, datas, dataSubKeys, subDataFetcher, SubData.GetId, setinHandler)
}

// 调整类型参数顺序，方便指定
func SetinX[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDataMapping ~map[Key]SubData,
	DataPtr ~*Data,
	Data any,
	SubData interface {
		GetId() Key
	},
	Key comparable,
](
	ctx context.Context,
	datas Datas,
	dataSubKeys func(DataPtr) Keys,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	setinHandler func(data DataPtr, subDataMapping SubDataMapping) DataPtr,
) error {
	return Setin(ctx, datas, dataSubKeys, subDataFetcher, setinHandler)
}

func ReversedSetinPro[
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDatasMapping ~map[Key]SubDatas,
	DataPtr ~*Data,
	Data any,
	SubData any,
	Key comparable,
](
	ctx context.Context,
	datas Datas,
	dataKey func(DataPtr) Key,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) DataPtr,
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
	Datas ~[]DataPtr,
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDatasMapping ~map[Key]SubDatas,
	DataPtr interface {
		~*Data
		GetId() Key
	},
	Data any,
	SubData any,
	Key comparable,
](
	ctx context.Context,
	datas Datas,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) DataPtr,
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
	SubDatas ~[]SubData,
	Keys ~[]Key,
	SubDatasMapping ~map[Key]SubDatas,
	DataPtr interface {
		~*Data
		GetId() Key
	},
	Data any,
	Key comparable,
	SubData any,
](
	ctx context.Context,
	datas Datas,
	subDataFetcher func(ctx context.Context, keys Keys) (SubDatas, error),
	subDataKeys func(SubData) Keys,
	setinHandler func(data DataPtr, subDatasMapping SubDatasMapping) DataPtr,
) error {
	return ReversedSetin(ctx, datas, subDataFetcher, subDataKeys, setinHandler)
}

type SetinHandlerMapping[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any] map[string]SetinHandler[Datas, DataPtr, Data]

func NewSetinHandlerMapping[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any]() SetinHandlerMapping[Datas, DataInstances, DataPtr, Data] {
	return SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]{}
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) Self() SetinHandlerMapping[Datas, DataInstances, DataPtr, Data] {
	return mapping
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) WithHandler(name string, handler SetinHandler[Datas, DataPtr, Data]) SetinHandlerMapping[Datas, DataInstances, DataPtr, Data] {
	mapping[name] = handler
	return mapping
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) WithComputedHandler(name string, handler SetinComputedHandler[Datas, DataPtr, Data]) SetinHandlerMapping[Datas, DataInstances, DataPtr, Data] {
	return mapping.WithHandler(name, NewSetinerHandlerFromComputedHandler(handler))
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) WithNamedHandler(handler *NamedSetinHandler[Datas, DataPtr, Data]) SetinHandlerMapping[Datas, DataInstances, DataPtr, Data] {
	if handler != nil {
		mapping.WithHandler(handler.name, handler.handler)
	}
	return mapping
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) SetinOne(ctx context.Context, datas Datas, setin string) (bool, error) {
	handler, ok := mapping[setin]
	if !ok {
		return false, nil
	}

	if err := handler(ctx, datas); err != nil {
		return false, err
	}

	return true, nil
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) Setin(ctx context.Context, datas Datas, setins []string) error {
	for _, setin := range setins {
		if ok, err := mapping.SetinOne(ctx, datas, setin); err != nil {
			return err
		} else if !ok {
			return NewUnknownSetinError(setin)
		}
	}

	return nil
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) NewTesters() SetinTesters[Datas, DataInstances, DataPtr, Data] {
	return NewSetinTesters[Datas, DataInstances](mapping.SetinOne)
}

func (mapping SetinHandlerMapping[Datas, DataInstances, DataPtr, Data]) NewSetiner() *Setiner[Datas, DataInstances, DataPtr, Data] {
	return mapping.NewTesters().NewSetiner()
}

type UnknownSetinError struct {
	setin   string
	setiner string
}

func NewUnknownSetinError(setin string) UnknownSetinError {
	return UnknownSetinError{
		setin: setin,
	}
}

func (e UnknownSetinError) Error() string {
	if e.setiner == "" {
		return fmt.Sprintf("unknown setin: %s", e.setin)
	}
	return fmt.Sprintf("unknown setin: %s, from setiner: %s", e.setin, e.setiner)
}

func (e UnknownSetinError) WithSetiner(setiner string) UnknownSetinError {
	e.setiner = setiner
	return e
}

type SetinTesters[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any] []SetinTester[Datas, DataPtr, Data]

func NewSetinTesters[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any](testers ...SetinTester[Datas, DataPtr, Data]) SetinTesters[Datas, DataInstances, DataPtr, Data] {
	return testers
}

func (testers SetinTesters[Datas, DataInstances, DataPtr, Data]) Append(more ...SetinTester[Datas, DataPtr, Data]) SetinTesters[Datas, DataInstances, DataPtr, Data] {
	return append(testers, more...)
}

func (testers SetinTesters[Datas, DataInstances, DataPtr, Data]) SetinOne(ctx context.Context, datas Datas, setin string) error {
	matched, err := testers.TestSetinOne(ctx, datas, setin)
	if err != nil {
		return err
	}

	if matched {
		return nil
	}

	return NewUnknownSetinError(setin)
}

func (testers SetinTesters[Datas, DataInstances, DataPtr, Data]) TestSetinOne(ctx context.Context, datas Datas, setin string) (matched bool, err error) {
	for _, tester := range testers {
		matched, err = tester(ctx, datas, setin)
		if err != nil {
			break
		}

		if matched {
			break
		}
	}

	return
}

func (testers SetinTesters[Datas, DataInstances, DataPtr, Data]) Setin(ctx context.Context, datas Datas, setins []string) error {
	for _, setin := range setins {
		if err := testers.SetinOne(ctx, datas, setin); err != nil {
			return err
		}
	}

	return nil
}

func (testers SetinTesters[Datas, DataInstances, DataPtr, Data]) NewSetiner() *Setiner[Datas, DataInstances, DataPtr, Data] {
	return NewSetiner[Datas, DataInstances, DataPtr, Data](testers)
}

type Setiner[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any] struct {
	SetinTesters[Datas, DataInstances, DataPtr, Data]
	subDataSetiner SubDataSetiner
}

func NewSetiner[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any](testers SetinTesters[Datas, DataInstances, DataPtr, Data]) *Setiner[Datas, DataInstances, DataPtr, Data] {
	return &Setiner[Datas, DataInstances, DataPtr, Data]{
		SetinTesters: testers,
	}
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) BindingDefaultSubDataSetiner() *Setiner[Datas, DataInstances, DataPtr, Data] {
	return setiner.BindingSubDataSetiner(GetDefaultSubDataSetiner())
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) BindingSubDataSetiner(subDataSetiner SubDataSetiner) *Setiner[Datas, DataInstances, DataPtr, Data] {
	return setiner.ProvidingSubDataSetiner(subDataSetiner).UsingSubDataSetiner(subDataSetiner)
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) ProvidingSubDataSetiner(subDataSetiner SubDataSetiner) *Setiner[Datas, DataInstances, DataPtr, Data] {
	var data Data
	subDataSetiner.WithHandlerByData(setiner.SetinFor, data)
	return setiner
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) Setin(ctx context.Context, datas Datas, setins []string) error {
	for _, setin := range setins {
		if err := setiner.SetinOne(ctx, datas, setin); err != nil {
			return err
		}
	}
	return nil
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) SetinFor(ctx context.Context, dataX any, setins []string) error {
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

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) SetinForData(ctx context.Context, data *Data, setins []string) error {
	return setiner.Setin(ctx, Datas{data}, setins)
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) SetinForDatas(ctx context.Context, datas Datas, setins []string) error {
	return setiner.Setin(ctx, datas, setins)
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) SetinForDataInstances(ctx context.Context, instances DataInstances, setins []string) error {
	return setiner.Setin(ctx, DatasFromInstances[Datas](instances), setins)
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) SetinOne(ctx context.Context, datas Datas, setin string) error {
	matched, err := setiner.SetinTesters.TestSetinOne(ctx, datas, setin)
	if err != nil {
		return err
	}

	if matched {
		return nil
	}

	subDataSetiner := setiner.subDataSetiner
	if subDataSetiner != nil {
		if strings.ContainsAny(setin, ".-(") {
			return subDataSetiner.setinOne(ctx, datas, setin)
		}
	}

	return NewUnknownSetinError(setin).WithSetiner(DataTypeIdent(setiner))
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) UsingSubDataSetiner(subDataSetiner SubDataSetiner) *Setiner[Datas, DataInstances, DataPtr, Data] {
	setiner.subDataSetiner = subDataSetiner
	return setiner
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) Action() *SetinAction[Datas, DataInstances, DataPtr, Data] {
	return setiner.ActionFor(nil)
}

func (setiner *Setiner[Datas, DataInstances, DataPtr, Data]) ActionFor(datas Datas) *SetinAction[Datas, DataInstances, DataPtr, Data] {
	return NewSetinAction(setiner, datas)
}

// A SetinAction is A Setiner with datas, it can Dup, Clone, WithData and finally Setin.
type SetinAction[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any] struct {
	*Setiner[Datas, DataInstances, DataPtr, Data]
	datas  Datas
	setins []string
}

func NewSetinAction[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any](setiner *Setiner[Datas, DataInstances, DataPtr, Data], datas Datas) *SetinAction[Datas, DataInstances, DataPtr, Data] {
	return &SetinAction[Datas, DataInstances, DataPtr, Data]{
		Setiner: setiner,
		datas:   datas,
	}
}

func (setiner *SetinAction[Datas, DataInstances, DataPtr, Data]) Dup() *SetinAction[Datas, DataInstances, DataPtr, Data] {
	dup := *setiner
	return &dup
}

func (setiner *SetinAction[Datas, DataInstances, DataPtr, Data]) Clone() ClonableSetinerInterface {
	return setiner
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithSetinsX(setins ...string) *SetinAction[Datas, DataInstances, DataPtr, Data] {
	action.setins = setins
	return action
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithSetins(setins ...string) CommonSetinerInterface {
	action.setins = setins
	return action
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithSetineds(setineds any) (CommonSetinerInterface, error) {
	return action.WithDataX(setineds)
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithDataX(dataX any) (*SetinAction[Datas, DataInstances, DataPtr, Data], error) {
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

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithData(data *Data) *SetinAction[Datas, DataInstances, DataPtr, Data] {
	if data == nil {
		action.datas = nil
	} else {
		action.datas = Datas{data}
	}
	return action
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithDatas(datas Datas) *SetinAction[Datas, DataInstances, DataPtr, Data] {
	action.datas = stl.PurgeZero(datas)
	return action
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) WithDataInstances(datas DataInstances) *SetinAction[Datas, DataInstances, DataPtr, Data] {
	action.datas = DatasFromInstances[Datas](datas)
	return action
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) SetinX(ctx context.Context, setins ...string) error {
	return action.Setin(ctx, setins)
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) Setin(ctx context.Context, setins []string) error {
	if len(setins) == 0 {
		setins = action.setins
	}
	return action.Setiner.SetinForDatas(ctx, action.datas, setins)
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) SetinForData(ctx context.Context, data *Data) error {
	return action.Setiner.SetinForData(ctx, data, action.setins)
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) SetinForDatas(ctx context.Context, datas Datas) error {
	return action.Setiner.SetinForDatas(ctx, datas, action.setins)
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) SetinForDataInstances(ctx context.Context, instances DataInstances) error {
	return action.Setiner.SetinForDataInstances(ctx, instances, action.setins)
}

func (action *SetinAction[Datas, DataInstances, DataPtr, Data]) SetinFor(ctx context.Context, dataX any) error {
	return action.Setiner.SetinFor(ctx, dataX, action.setins)
}

func DatasFromInstances[Datas ~[]DataPtr, DataInstances ~[]Data, DataPtr ~*Data, Data any](instances DataInstances) Datas {
	datas := make(Datas, 0, len(instances))
	for i, _ := range instances {
		datas = append(datas, &instances[i])
	}
	return datas
}
