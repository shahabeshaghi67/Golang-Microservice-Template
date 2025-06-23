package utils

import (
	"reflect"
	"strings"
)

// Index takes a slice and an indexing function f and creates a map by indexing each element
// of the slice using result of f as map's key.
func Index[S ~[]T, T any, K comparable](s S, f func(t T) K) map[K]T {
	res := make(map[K]T)
	for _, val := range s {
		res[f(val)] = val
	}
	return res
}

func JsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", -1)[0]
	if name == "-" {
		return ""
	}
	return name
}
