package util

type any []interface{}

func Extend(src any, inc int) any {
	increment := make(any, inc)
	return append(src, increment...)
}
