package contest

import (
	"fmt"
)

type NumQuestion int

func (n NumQuestion) Int() int32 {
	return int32(n)
}

func newNumQuestion(n int32) (NumQuestion, error) {
	if n <= 0 {
		return 0, fmt.Errorf("num questions cannot be negative")
	}
	return NumQuestion(n), nil
}
