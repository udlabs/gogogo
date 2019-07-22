package util

import "strings"

func JoinBy(separator string, args ...string) string{
	return strings.Join(args, separator)
}
