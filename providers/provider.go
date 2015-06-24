package providers

import (
	"sort"
)

func median(v []float64) float64 {
	sort.Float64s(v)
	// [0 1 (2) 3 4]
	if len(v)%2 == 1 {
		return v[len(v)/2]
	}
	// interpolate the 2 middle
	return (v[len(v)/2] + v[len(v)/2+1]) / 2
}
