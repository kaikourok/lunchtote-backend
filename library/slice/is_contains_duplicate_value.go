package slice

import "golang.org/x/exp/constraints"

type Comparable interface {
	constraints.Ordered | rune | string
}

func IsContainsDuplicateValue[T Comparable](arr *[]T) bool {
	visiteds := make(map[T]bool, len(*arr))

	for i := range *arr {
		if visiteds[(*arr)[i]] {
			return true
		} else {
			visiteds[(*arr)[i]] = true
		}
	}

	return false
}
