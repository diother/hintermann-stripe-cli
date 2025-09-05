package helper

import (
	"fmt"
	"strconv"
)

func MustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("bad int in CSV: %q", s))
	}
	return n
}
