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
	"hash"
	"math"
	"math/bits"
)

// const (
// 	exp32 = 1 << 32 // 2^32
// )

// This is a 64-bit version of the HyperLogLog algorithm that use 64-bit hash
// and is thus more accurate than the 32-bit version and can be used to count
// values much larger than 2^32. It uses exactly the same amount of memory as
// the 32-bit version with the same configs.
//
// A HyperLogLog64 is not compatible with a HyperLogLog and cannot be merged
// with a HyperLogLog.
type HyperLogLog64 struct {
	M         uint    // Number of registers
	B         uint32  // Number of bits used to determine register index
	Alpha     float64 // Bias correction constant
	Registers []uint8
}

// // Compute bias correction alpha_m.
// func getAlpha(m uint) (result float64) {
// 	switch m {
// 	case 16:
// 		result = 0.673
// 	case 32:
// 		result = 0.697
// 	case 64:
// 		result = 0.709
// 	default:
// 		result = 0.7213 / (1.0 + 1.079/float64(m))
// 	}
// 	return result
// }

// New creates a HyperLogLog64 with the given number of registers. More
// registers leads to lower error in your estimated count, at the expense of
// memory.
//
// Choose a power of two number of registers, depending on the amount
// of memory you're willing to use and the error you're willing to
// tolerate. Each register uses one byte of memory.
//
// Standard error will be: σ ≈ 1.04 / sqrt(registers)
// The estimates provided by hyperloglog are expected to be within σ, 2σ, 3σ
// of the exact count in respectively 65%, 95%, 99% of all the cases.
func New64(registers uint) (*HyperLogLog64, error) {
	if registers == 0 {
		panic("cannot have zero registers")
	}
	if (registers & (registers - 1)) != 0 {
		return nil, fmt.Errorf("number of registers %d not a power of two", registers)
	}
	h := &HyperLogLog64{}
	h.M = registers
	h.B = uint32(math.Log2(float64(registers)))
	h.Alpha = getAlpha(registers)
	h.Registers = make([]uint8, h.M)
	return h, nil
}

// Reset all internal variables and set the count to zero.
func (h *HyperLogLog64) Reset() {
	for i := range h.Registers {
		h.Registers[i] = 0
	}
}

// Add a hash directly to the count. val should be a 64 bit unsigned integer
// from a good hash function.
//
// To be able to merge with another HyperLogLog64, hash value must be
// deterministic in the sense that the same input will always produce the same
// hash.
func (h *HyperLogLog64) AddHash(val uint64) {
	k := 64 - h.B
	slice := (val << h.B) | (1 << (h.B - 1))
	r := uint8(bits.LeadingZeros64(slice) + 1)
	j := val >> uint(k)
	if r > h.Registers[j] {
		h.Registers[j] = r
	}
}

// Add a Hashable object to the count.
//
// To be able to merge with another HyperLogLog64, the object's hash function
// must be deterministic in the sense that the same input will always produce
// the same hash.
func (h *HyperLogLog64) Add(obj hash.Hash64) {
	h.AddHash(obj.Sum64())
}

// Count returns the estimated cardinality.
func (h *HyperLogLog64) Count() uint64 {
	return h.count(true)
}

// CountWithoutLargeRangeCorrection returns the estimated cardinality, without applying
// the large range correction proposed by Flajolet et al. as it can lead to significant
// overcounting.
//
// See https://github.com/DataDog/hyperloglog/pull/15
func (h *HyperLogLog64) CountWithoutLargeRangeCorrection() uint64 {
	return h.count(false)
}

func (h *HyperLogLog64) count(withLargeRangeCorrection bool) uint64 {
	sum := 0.0
	m := float64(h.M)
	for _, val := range h.Registers {
		sum += 1.0 / float64(int(1)<<val)
	}
	estimate := h.Alpha * m * m / sum
	if estimate <= 5.0/2.0*m {
		// Small range correction
		v := 0
		for _, r := range h.Registers {
			if r == 0 {
				v++
			}
		}
		if v > 0 {
			estimate = m * math.Log(m/float64(v))
		}
	} else if estimate > 1.0/30.0*exp32 && withLargeRangeCorrection {
		// Large range correction
		estimate = -exp32 * math.Log(1-estimate/exp32)
	}
	return uint64(estimate)
}

// Merge another HyperLogLog64 into this one. The number of registers in each
// must be the same. `other` must be using the same deterministic hash function
// as `h`.
func (h *HyperLogLog64) Merge(other *HyperLogLog64) error {
	if h.M != other.M {
		return fmt.Errorf("number of registers doesn't match: %d != %d",
			h.M, other.M)
	}
	for j, r := range other.Registers {
		if r > h.Registers[j] {
			h.Registers[j] = r
		}
	}
	return nil
}
