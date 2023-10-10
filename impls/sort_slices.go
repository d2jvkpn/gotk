package impls

import (
	"sort"

	"golang.org/x/exp/constraints"
)

/*
demo:

	[[1,2], [3]] => [[3], [1,2]]
	[[2,1], [1,2]] => [[1,2], [2, 1]]
*/
func SortSliceOfSlice[T constraints.Ordered](slice [][]T) {
	sort.SliceStable(slice, func(i int, j int) bool {
		if len(slice[i]) != len(slice[j]) {
			return len(slice[i]) < len(slice[j])
		}

		for k := range slice[i] {
			if slice[i][k] != slice[j][k] {
				return slice[i][k] < slice[j][k]
			}
		}

		return false
	})
}
