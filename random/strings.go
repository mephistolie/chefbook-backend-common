package random

import (
	"math/rand"
	"time"
)

var digits = []rune("0123456789")

func DigitString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	b := make([]rune, length)
	for i := range b {
		b[i] = digits[r.Intn(len(digits))]
	}
	return string(b)
}
