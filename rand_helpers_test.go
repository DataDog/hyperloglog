package hyperloglog

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"
)

// RandFill does what `math/rand.Read` did -- fill a slice of bytes with random
// values. The `math/rand/v2` package does not include a `Read` method
func RandFill(dst []byte) {
	remnant := len(dst) % 8
	for i := 0; i < len(dst)-remnant; i += 8 {
		binary.LittleEndian.PutUint64(dst[i:], rand.Uint64())
	}
	if remnant == 0 {
		return
	}
	v := rand.Uint64()
	for i := len(dst) - remnant; i < len(dst); i++ {
		dst[i] = byte(v % 256)
		v >>= 8
	}
}

func TestRandFill(t *testing.T) {
	// This test is going to make slices of different sizes
	// Each slice will be filled 5 times and its data OR'ed onto an aggregate
	// The aggregate should have no zero-valued bytes at the end
	// The probability of a zero-valued byte should be (1/256)**5, which is ~9.1E-13
	// 3 is probably sufficient, but this is fast enough
	for test := range 256 {
		data := make([]byte, test)
		agg := make([]byte, test)

		for range 5 {
			RandFill(data)
			for i := range test {
				agg[i] |= data[i]
			}
		}

		for i := range test {
			if agg[i] == 0 {
				t.Errorf("failed to fill byte at offset %d in slice of length %d", i, test)
				break
			}
		}
	}
}
