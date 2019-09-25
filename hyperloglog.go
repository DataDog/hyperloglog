// Package hyperloglog implements the HyperLogLog algorithm for
// cardinality estimation. In English: it counts things. It counts
// things using very small amounts of memory compared to the number of
// objects it is counting.
//
// For a full description of the algorithm, see the paper HyperLogLog:
// the analysis of a near-optimal cardinality estimation algorithm by
// Flajolet, et. al.
package hyperloglog

import (
	"fmt"
	"math"
)

var (
	exp32 = float64(4294967296) // 2**32
)

// A HyperLogLog is a deterministic cardinality estimator.  This version
// exports its fields so that it is suitable for saving eg. to a database.
type HyperLogLog struct {
	registers []uint8
	alpha     float64 // Bias correction constant
	b         uint8   // Number of bits used to determine register index
	m         int     // Number of registers
}

// New creates a HyperLogLog with the given number of registers. More
// registers leads to lower error in your estimated count, at the
// expense of memory.
//
// Choose a power of two number of registers, depending on the amount
// of memory you're willing to use and the error you're willing to
// tolerate. Each register uses one byte of memory.
//
// Approximate error will be:
//     1.04 / sqrt(registers)
//
func New(m int) *HyperLogLog {
	if (m & (m - 1)) != 0 {
		panic(fmt.Errorf("number of registers %d not a power of two", m))
	}

	return &HyperLogLog{
		registers: make([]uint8, m),
		alpha:     getAlpha(m),
		b:         getLog(m),
		m:         m,
	}
}

// Reset all internal variables and set the count to zero.
func (h *HyperLogLog) Reset() {
	for i := 0; i < h.m; i++ {
		h.registers[i] = 0
	}
}

// Add to the count. val should be a 64 bit unsigned integer from a
// good hash function.
func (h *HyperLogLog) Add(val uint32) {
	k := 32 - h.b
	r := rho(val<<h.b, k)
	j := val >> k

	if r > h.registers[j] {
		h.registers[j] = r
	}
}

// Count returns the estimated cardinality.
func (h *HyperLogLog) Count() uint64 {
	sum := 0.0
	m := float64(h.m)
	for _, val := range h.registers {
		sum += 1.0 / float64(uint64(1)<<val)
	}
	estimate := h.alpha * m * m / sum

	if estimate <= 2.5*m {
		// Small range correction
		v := 0
		for _, r := range h.registers {
			if r == 0 {
				v++
			}
		}
		if v > 0 {
			estimate = m * math.Log(m/float64(v))
		}
	} else if estimate > 0.03*exp32 {
		// Large range correction
		estimate = -exp32 * math.Log(1-estimate/exp32)
	}
	return uint64(estimate)
}

// Merge another HyperLogLog into this one. The number of registers in
// each must be the same.
func (h *HyperLogLog) Merge(other *HyperLogLog) {
	if h.m != other.m {
		panic(fmt.Errorf("number of registers doesn't match: %d != %d", h.m, other.m))
	}

	for i := 0; i < h.m; i++ {
		if other.registers[i] > h.registers[i] {
			h.registers[i] = other.registers[i]
		}
	}
}

// Calculate the position of the leftmost 1-bit.
func rho(val uint32, max uint8) uint8 {
	r := uint8(1)
	for val&0x80000000 == 0 && r <= max {
		r++
		val <<= 1
	}
	return r
}

// Compute bias correction alpha_m.
func getAlpha(m int) (result float64) {
	switch m {
	case 16:
		result = 0.673102023867666
	case 32:
		result = 0.6971226338010241
	case 64:
		result = 0.7092084528700233
	case 128:
		result = 0.7152711899613394
	case 256:
		result = 0.7183076381918139
	case 512:
		result = 0.7198271478204001
	case 1024:
		result = 0.7205872259764527
	case 2048:
		result = 0.720967346136219
	case 4096:
		result = 0.7211574265173785
	default:
		result = 0.721347607152952 / (1.0 + 1.08018007534368/float64(m))
	}
	return result
}

// Calculate the number of bits necessary to reprsent integers
// from 0 to m-1 (logarithm in base 2 of m)
func getLog(m int) uint8 {
	r := uint8(0)
	for m&1 == 0 {
		r++
		m >>= 1
	}
	return r
}
