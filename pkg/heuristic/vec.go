package heuristic

import (
	"math"
	"sort"
)

type Vec []float64

func (vec Vec) New() Vec {
	return make(Vec, 0, vec.Len())
}

func (vec Vec) Copy() Vec {
	return append(vec.New(), vec...)
}

func (vec Vec) Sorted() Vec {
	var sorted = vec.Copy()
	sort.Float64s(sorted)
	return sorted
}

func (vec Vec) Median() float64 {
	var sorted = vec.Sorted()
	var l = sorted.Len()
	if l == 0 {
		return 0
	}
	if l%2 == 0 {
		return sorted[l/2]/2 + sorted[l/2-1]/2
	}
	return sorted[l/2]
}

func (vec Vec) Sum() float64 {
	var s float64
	for _, p := range vec {
		s += p
	}
	return s
}

func (vec Vec) Avg() float64 {
	var l = float64(len(vec))
	return vec.Sum() / l
}

func (vec Vec) Len() int {
	return len(vec)
}

func (vec Vec) Max() float64 {
	if vec.Len() == 0 {
		return 0
	}
	if vec.Len() == 1 {
		return vec[0]
	}
	var max = vec[0]
	for _, v := range vec[1:] {
		max = math.Max(max, v)
	}
	return max
}

func (vec Vec) Map(op func(x float64) float64) Vec {
	var mapped = vec.New()
	for _, x := range vec {
		mapped = append(mapped, op(x))
	}
	return mapped
}
