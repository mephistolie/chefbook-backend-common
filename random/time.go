package random

import (
	"math/rand"
	"time"
)

func DurationSeconds(limit int) time.Duration {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	offsetTime := int64(r.Intn(limit+1)) * time.Second.Nanoseconds()
	return time.Duration(offsetTime)
}
