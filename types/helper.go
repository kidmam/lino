package types

import (
	"strconv"
)

// MustParseInt64 - panic if error
func MustParseInt64(s string, base int, bitSize int) int64 {
	rst, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		panic(err)
	}
	return rst
}
