package gotk

import (
	"golang.org/x/exp/constraints"
)

func VectorIndex[T constraints.Ordered](list []T, v T) int {
	for i := range list {
		if list[i] == v {
			return i
		}
	}

	return -1
}

func EqualVector[T constraints.Ordered](arr1, arr2 []T) (ok bool) {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func UniqVector[T constraints.Ordered](arr []T) (list []T) {
	n := len(arr)
	list = make([]T, 0, n)

	if len(arr) == 0 {
		return list
	}

	mp := make(map[T]bool, n)
	for _, v := range arr {
		if !mp[v] {
			list = append(list, v)
			mp[v] = true
		}
	}

	return list
}

func First[T any](v []T) *T {
	if len(v) == 0 {
		return nil
	}
	return &v[0]
}

func Last[T any](v []T) *T {
	var n = len(v)

	if n == 0 {
		return nil
	}
	return &v[n-1]
}

func SliceGet[T any](slice []T, index int) (val T, exists bool) {
	if index > len(slice)-1 {
		return val, false
	}

	return slice[index], true
}

func Slice2Map[K comparable, T, V any](items []T, getKey func(*T) K, getValue func(*T) V) (
	mp map[K]V) {
	mp = make(map[K]V, len(items))

	for i := range items {
		mp[getKey(&items[i])] = getValue(&items[i])
	}

	return mp
}

func JoinSlices[K comparable, T, V any](items []T, exts map[K]V,
	getKey func(*T) K, setValue func(*T, V)) (n int) {
	var (
		v  V
		ok bool
	)

	if len(exts) == 0 {
		return 0
	}

	for i := range items {
		if v, ok = exts[getKey(&items[i])]; ok {
			n++
			setValue(&items[i], v)
		}
	}

	return n
}

func PickOne[T any](items []T) (item T) {
	if len(items) == 0 {
		return item
	}

	return items[Rand.IntN(len(items))]
}

func PickSome[T any](items []T, num int) (ans []T) {
	if len(items) <= num {
		return items
	}

	ans = make([]T, num)
	if num == 0 {
		return ans
	}

	slice := Rand.Perm(len(items))
	for i := 0; i < num; i++ {
		ans[i] = items[slice[i]]
	}

	return ans
}

func PickSomeIndex[T any](items []T, num int) (ans []int) {
	if len(items) <= num {
		ans = make([]int, len(items))
		for i := range items {
			ans[i] = i
		}

		return ans
	}

	ans = make([]int, num)
	if num == 0 {
		return ans
	}

	slice := Rand.Perm(len(items))
	for i := 0; i < num; i++ {
		ans[i] = slice[i]
	}

	return ans
}

func NewSliceWith[T any](item T, caps ...int) (slice []T) {
	if len(caps) > 0 { // caps[0] > 0
		slice = make([]T, 0, caps[0])
	} else {
		slice = make([]T, 0, 1)
	}

	slice = append(slice, item)
	return slice
}
