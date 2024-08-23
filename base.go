package gotk

import (
	// "fmt"
	"math/rand/v2"
	"time"
)

var (
	Rand *rand.Rand
)

func init() {
	source := rand.NewPCG(42, uint64(time.Now().UnixNano()))
	Rand = rand.New(source)
}
