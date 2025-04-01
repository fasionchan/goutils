/*
 * Author: fasion
 * Created time: 2025-04-01 21:28:04
 * Last Modified by: fasion
 * Last Modified time: 2025-04-01 22:43:33
 */

package types

import "github.com/fasionchan/goutils/stl"

func Float32ToFloat64(value float32) float64 {
	return float64(value)
}

func Float64ToFloat32(value float64) float32 {
	return float32(value)
}

type FloatVector = Float32Vector

type Float32Vector []float32

func (vector Float32Vector) IsNil() bool {
	return vector == nil
}

func (vector Float32Vector) Empty() bool {
	return len(vector) == 0
}

func (vector Float32Vector) Len() int {
	return len(vector)
}

func (vector Float32Vector) Dim() int {
	return len(vector)
}

func (vector Float32Vector) Float32() Float32Vector {
	return vector
}

func (vector Float32Vector) Float64() Float64Vector {
	return stl.Map(vector, Float32ToFloat64)
}

func (vector Float32Vector) Native() []float32 {
	return vector
}

type Float64Vector []float64

func (vector Float64Vector) IsNil() bool {
	return vector == nil
}

func (vector Float64Vector) Empty() bool {
	return len(vector) == 0
}

func (vector Float64Vector) Len() int {
	return len(vector)
}

func (vector Float64Vector) Dim() int {
	return len(vector)
}

func (vector Float64Vector) Float32() Float32Vector {
	return stl.Map(vector, Float64ToFloat32)
}

func (vector Float64Vector) Float64() Float64Vector {
	return vector
}

func (vector Float64Vector) Native() []float64 {
	return vector
}

type FloatVectors = Float64Vectors

type Float32Vectors []Float32Vector

func (vectors Float32Vectors) Len() int {
	return len(vectors)
}

func (vectors Float32Vectors) FirstOneOrNil() Float32Vector {
	return stl.FirstOneOrZero(vectors)
}

func (vectors Float32Vectors) Float32() Float32Vectors {
	return vectors
}

func (vectors Float32Vectors) Float64() Float64Vectors {
	return stl.Map(vectors, Float32Vector.Float64)
}

func (vectors Float32Vectors) PurgeEmpty() Float32Vectors {
	return stl.Purge(vectors, Float32Vector.Empty)
}

func (vectors Float32Vectors) PurgeNil() Float32Vectors {
	return stl.Purge(vectors, Float32Vector.IsNil)
}

type Float64Vectors []Float64Vector

func (vectors Float64Vectors) Len() int {
	return len(vectors)
}

func (vectors Float64Vectors) FirstOneOrNil() Float64Vector {
	return stl.FirstOneOrZero(vectors)
}

func (vectors Float64Vectors) Float32() Float32Vectors {
	return stl.Map(vectors, Float64Vector.Float32)
}

func (vectors Float64Vectors) Float64() Float64Vectors {
	return vectors
}

func (vectors Float64Vectors) PurgeEmpty() Float64Vectors {
	return stl.Purge(vectors, Float64Vector.Empty)
}

func (vectors Float64Vectors) PurgeNil() Float64Vectors {
	return stl.Purge(vectors, Float64Vector.IsNil)
}
