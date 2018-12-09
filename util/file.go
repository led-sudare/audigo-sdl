package util

import "strings"

func ext(file string) string {
	pos := strings.LastIndex(file, ".")
	return file[pos:]
}
