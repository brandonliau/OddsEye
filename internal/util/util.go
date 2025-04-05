package util

import (
	"maps"

	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](data []T) T {
	var maxv T
	for _, i := range data {
		if i > maxv {
			maxv = i
		}
	}
	return maxv
}

func Keys[T comparable, K any](data map[T]K) []T {
	keys := make([]T, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func MapCopy[T comparable, K any](data map[T]K) map[T]K {
	copy := make(map[T]K)
	maps.Copy(copy, data)
	return copy
}
