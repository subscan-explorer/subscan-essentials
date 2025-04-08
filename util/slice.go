package util

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func MapStringToSlice(m map[string]bool) []string {
	var l []string
	for v := range m {
		l = append(l, v)
	}
	return l
}

func ContinuousNums(start, count int, order string) (r []int) {
	if count <= 0 {
		return
	}
	for i := 0; i < count; i++ {
		if order == "desc" {
			if start-i < 0 {
				break
			}
			r = append(r, start-i)
		} else {
			r = append(r, start+i)
		}
	}
	return
}

func SortUintSlice(s []uint) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

func StringInSliceFold(a string, list []string) bool {
	return SliceIndex(a, list, true) != -1
}

func SliceIndex(a string, list []string, fold bool) int {
	for index, b := range list {
		if fold {
			if strings.EqualFold(b, a) {
				return index
			}
		} else {
			if b == a {
				return index
			}
		}
	}
	return -1
}

// reverse
func Reverse(a interface{}) interface{} {
	switch reflect.TypeOf(a).Kind() {
	case reflect.Slice:
		size := reflect.ValueOf(a).Len()
		swap := reflect.Swapper(a)
		for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
			swap(i, j)
		}
		return a
	default:
		panic(fmt.Errorf("invalid Enumerable"))
	}
}
